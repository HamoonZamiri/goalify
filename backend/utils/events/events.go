package events

import (
	"encoding/json"
	"fmt"
	"goalify/utils/lists"
	"goalify/utils/options"
	"log/slog"
	"reflect"
	"sync"
)

type Event struct {
	Data      any                    `json:"data"`
	EventType string                 `json:"event_type"`
	UserId    options.Option[string] `json:"user_id"`
}

const (
	QUEUE_MAX_SIZE        int    = 1000
	USER_CREATED          string = "user_created"
	GOAL_CREATED          string = "goal_created"
	GOAL_UPDATED          string = "goal_updated"
	USER_UPDATED          string = "user_updated"
	GOAL_CATEGORY_CREATED string = "goal_category_created"
	DEFAULT_GOAL_CREATED  string = "default_goal_created"
	SSE_CONNECTED         string = "sse_connected"
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
	SubscribeToUserEvents(userId string, subscriber Subscriber)
	UnsubscribeFromUserEvents(userId string, subscriber Subscriber)
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
		UserId:    options.None[string](),
	}
}

func NewEventWithUserId(eventType string, data any, userId string) Event {
	return Event{
		Data:      data,
		EventType: eventType,
		UserId:    options.Some(userId),
	}
}

func NewEventManager() *EventManager {
	em := &EventManager{
		eventQueue:  make(chan Event, QUEUE_MAX_SIZE),
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

func (em *EventManager) broadcastByType(event Event) {
	sublist, ok := em.subscribers[event.EventType]
	if !ok {
		return
	}
	underlyingList := sublist.GetList()
	for e := underlyingList.Front(); e != nil; e = e.Next() {
		sub, ok := e.Value.(Subscriber)
		if !ok {
			slog.Warn("EventManager.Publish: type assertion failed", "subscriber", e.Value)
		}
		sub.HandleEvent(event)
	}
}

func (em *EventManager) broadcastByUser(event Event) {
	if !event.UserId.IsPresent() {
		return
	}

	userId := event.UserId.ValueOrZero()
	userSubList, ok := em.userSubs[userId]
	if !ok {
		return
	}
	underlyingList := userSubList.GetList()
	for e := underlyingList.Front(); e != nil; e = e.Next() {
		sub, ok := e.Value.(Subscriber)
		if !ok {
			slog.Warn("EventManager.Publish: type assertion failed", "subscriber", e.Value)
		}
		sub.HandleEvent(event)
	}
}

func (em *EventManager) processEvents() {
	for event := range em.eventQueue {
		em.mu.Lock()
		var wg sync.WaitGroup
		wg.Add(2)

		// publish events based on their event type (internal to the monolith)
		go func() {
			em.broadcastByType(event)
			wg.Done()
		}()
		// publish events based on userId external clients such as the frontend
		go func() {
			em.broadcastByUser(event)
			wg.Done()
		}()
		wg.Wait()
		em.mu.Unlock()
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

func (em *EventManager) SubscribeToUserEvents(userId string, subscriber Subscriber) {
	em.mu.Lock()
	defer em.mu.Unlock()
	if _, ok := em.userSubs[userId]; !ok {
		em.userSubs[userId] = lists.New[Subscriber]()
	}
	em.userSubs[userId].PushBack(subscriber)
}

func (em *EventManager) UnsubscribeFromUserEvents(userId string, subscriber Subscriber) {
	em.mu.Lock()
	defer em.mu.Unlock()
	subList := em.userSubs[userId].GetList()
	for e := subList.Front(); e != nil; e = e.Next() {
		if e.Value == subscriber {
			subList.Remove(e)
			break
		}
	}
}
