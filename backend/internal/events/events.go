// Package events contains utilities for publishing and subscribing to events
package events

import (
	"encoding/json"
	"fmt"
	"goalify/pkg/lists"
	"goalify/pkg/options"
	"reflect"
	"sync"
)

type Event struct {
	Data      any                    `json:"data"`
	EventType string                 `json:"event_type"`
	UserID    options.Option[string] `json:"user_id"`
}

const (
	QueueMaxSize        int    = 1000
	UserCreated         string = "user_created"
	GoalCreated         string = "goal_created"
	GoalUpdated         string = "goal_updated"
	UserUpdated         string = "user_updated"
	GoalCategoryCreated string = "goal_category_created"
	DefaultGoalCreated  string = "default_goal_created"
	SSEConnected        string = "sse_connected"
	XPUpdated           string = "xp_updated"
)

func ParseEventData[T any](event Event) (T, error) {
	var val T
	var ok bool
	val, ok = event.Data.(T)
	valType := reflect.TypeOf(val).String()
	if !ok {
		return val, fmt.Errorf("%T: type assertion failed", valType)
	}
	return val, nil
}

func (e *Event) EncodeEvent() ([]byte, error) {
	return json.Marshal(e.Data)
}

type Subscriber interface {
	HandleEvent(event Event)
}

type EventPublisher interface {
	Subscribe(eventType string, subscriber Subscriber)
	Publish(event Event)
	Unsubscribe(eventType string, subscriber Subscriber)
	SubscribeToUserEvents(userID string, subscriber Subscriber)
	UnsubscribeFromUserEvents(userID string, subscriber Subscriber)
}

type EventManager struct {
	eventQueue  chan Event
	subscribers map[string]*lists.TypedList[Subscriber]
	userSubs    map[string]*lists.TypedList[Subscriber]
	mu          sync.Mutex
}

func NewEvent(eventType string, data any) Event {
	return Event{
		Data:      data,
		EventType: eventType,
		UserID:    options.None[string](),
	}
}

func NewEventWithUserID(eventType string, data any, userID string) Event {
	return Event{
		Data:      data,
		EventType: eventType,
		UserID:    options.Some(userID),
	}
}

func NewEventManager() *EventManager {
	em := &EventManager{
		eventQueue:  make(chan Event, QueueMaxSize),
		subscribers: make(map[string]*lists.TypedList[Subscriber]),
		userSubs:    make(map[string]*lists.TypedList[Subscriber]),
		mu:          sync.Mutex{},
	}
	go em.processEvents()
	return em
}

func (em *EventManager) IsSubscribed(eventType string, subscriber Subscriber) bool {
	em.mu.Lock()
	defer em.mu.Unlock()

	if _, ok := em.subscribers[eventType]; !ok {
		return false
	}
	subList := em.subscribers[eventType].GetList()
	for e := subList.Front(); e != nil; e = e.Next() {
		if e.Value == subscriber {
			return true
		}
	}
	return false
}

func (em *EventManager) Subscribe(eventType string, subscriber Subscriber) {
	em.mu.Lock()
	defer em.mu.Unlock()
	if _, ok := em.subscribers[eventType]; !ok {
		em.subscribers[eventType] = lists.New[Subscriber]()
	}
	em.subscribers[eventType].PushBack(subscriber)
}

func (em *EventManager) Publish(event Event) {
	em.eventQueue <- event
}

// snapshotSubscribers copies, under the lock, the subscribers that should
// receive an event: those registered for its type (internal services) plus
// those registered for its user (external SSE/WebSocket clients). Delivery then
// happens off the lock so a slow or blocking subscriber can never freeze
// publishing, subscribing, or disconnecting for everyone else.
func (em *EventManager) snapshotSubscribers(event Event) []Subscriber {
	em.mu.Lock()
	defer em.mu.Unlock()

	var subs []Subscriber
	collect := func(list *lists.TypedList[Subscriber]) {
		for e := list.GetList().Front(); e != nil; e = e.Next() {
			if sub, ok := e.Value.(Subscriber); ok {
				subs = append(subs, sub)
			}
		}
	}

	if list, ok := em.subscribers[event.EventType]; ok {
		collect(list)
	}
	if event.UserID.IsPresent() {
		if list, ok := em.userSubs[event.UserID.ValueOrZero()]; ok {
			collect(list)
		}
	}
	return subs
}

func (em *EventManager) processEvents() {
	for event := range em.eventQueue {
		for _, sub := range em.snapshotSubscribers(event) {
			sub.HandleEvent(event)
		}
	}
}

func (em *EventManager) Unsubscribe(eventType string, subscriber Subscriber) {
	em.mu.Lock()
	defer em.mu.Unlock()
	subList := em.subscribers[eventType].GetList()
	for e := subList.Front(); e != nil; e = e.Next() {
		if e.Value == subscriber {
			subList.Remove(e)
			break
		}
	}
}

func (em *EventManager) SubscribeToUserEvents(userID string, subscriber Subscriber) {
	em.mu.Lock()
	defer em.mu.Unlock()
	if _, ok := em.userSubs[userID]; !ok {
		em.userSubs[userID] = lists.New[Subscriber]()
	}
	em.userSubs[userID].PushBack(subscriber)
}

func (em *EventManager) UnsubscribeFromUserEvents(userID string, subscriber Subscriber) {
	em.mu.Lock()
	defer em.mu.Unlock()
	subList := em.userSubs[userID].GetList()
	for e := subList.Front(); e != nil; e = e.Next() {
		if e.Value == subscriber {
			subList.Remove(e)
			break
		}
	}
}
