package events

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CountEventHandler struct {
	count int
}

func (h *CountEventHandler) HandleIncrementEvent(event Event) {
	h.count++
}

func (h *CountEventHandler) HandleDecrementEvent(event Event) {
	h.count--
}

func (h *CountEventHandler) HandleEvent(event Event) {
	switch event.EventType {
	case "increment":
		h.HandleIncrementEvent(event)
	case "decrement":
		h.HandleDecrementEvent(event)
	default:
		slog.Error("unhandled event type")
	}
}

func NewCountEventHandler() *CountEventHandler {
	return &CountEventHandler{
		count: 0,
	}
}

func TestSubscribe(t *testing.T) {
	em := NewEventManager()
	h := NewCountEventHandler()

	em.Subscribe("increment", h)
	em.Subscribe("decrement", h)

	assert.True(t, em.IsSubscribed("increment", h))
	assert.True(t, em.IsSubscribed("decrement", h))
}

func TestUnsubscribe(t *testing.T) {
	em := NewEventManager()
	h := NewCountEventHandler()

	em.Subscribe("increment", h)
	assert.True(t, em.IsSubscribed("increment", h))

	em.Unsubscribe("increment", h)
	assert.False(t, em.IsSubscribed("increment", h))
}

func TestPublishEvent(t *testing.T) {
	em := NewEventManager()
	h := NewCountEventHandler()

	em.Subscribe("increment", h)
	em.Subscribe("decrement", h)
	incEvent := NewEvent("increment", 1)
	em.Publish(incEvent)
	assert.Equal(t, 1, h.count)

	decEvent := NewEvent("decrement", 1)
	em.Publish(decEvent)
	assert.Equal(t, 0, h.count)
}

func TestParseEvent(t *testing.T) {
	e := NewEvent("increment", 1)
	eventData, err := ParseEventData[int](e)
	assert.Nil(t, err)
	assert.Equal(t, 1, eventData)
}

func TestMultipleSubscribers(t *testing.T) {
	em := NewEventManager()
	subscribers := make([]*CountEventHandler, 0)
	for i := 0; i < 5; i++ {
		eh := NewCountEventHandler()
		subscribers = append(subscribers, eh)
		em.Subscribe("increment", eh)
	}

	assert.Equal(t, 5, em.subscribers["increment"].Len())

	for i := 0; i < 5; i++ {
		em.Publish(NewEvent("increment", 20))
	}

	for _, sub := range subscribers {
		assert.Equal(t, 5, sub.count)
	}
}
