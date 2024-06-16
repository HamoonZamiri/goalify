package events

import (
	"fmt"
	"goalify/utils/lists"
	"reflect"
)

type Event struct {
	Data      any
	EventType string
}

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

type Subscriber interface {
	HandleEvent(event Event)
}

type EventPublisher interface {
	Subscribe(eventType string, subscriber Subscriber)
	Publish(event Event)
	Unsubscribe(eventType string, subscriber Subscriber)
}

type EventManager struct {
	subscribers map[string]*lists.TypedList[Subscriber]
}

func NewEvent(eventType string, data any) Event {
	return Event{
		Data:      data,
		EventType: eventType,
	}
}

func NewEventManager() *EventManager {
	return &EventManager{
		subscribers: make(map[string]*lists.TypedList[Subscriber]),
	}
}

func (em *EventManager) IsSubscribed(eventType string, subscriber Subscriber) bool {
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
	if _, ok := em.subscribers[eventType]; !ok {
		em.subscribers[eventType] = lists.New[Subscriber]()
	}
	em.subscribers[eventType].PushBack(subscriber)
}

func (em *EventManager) Publish(event Event) {
	subList := em.subscribers[event.EventType].GetList()
	for e := subList.Front(); e != nil; e = e.Next() {
		sub := e.Value.(Subscriber)
		sub.HandleEvent(event)
	}
}

func (em *EventManager) Unsubscribe(eventType string, subscriber Subscriber) {
	subList := em.subscribers[eventType].GetList()
	for e := subList.Front(); e != nil; e = e.Next() {
		if e.Value == subscriber {
			subList.Remove(e)
			break
		}
	}
}
