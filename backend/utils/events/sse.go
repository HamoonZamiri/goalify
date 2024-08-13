package events

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type SSEConn struct {
	eventQueue chan Event
	writer     http.ResponseWriter
	userId     string
}

func NewSSEConn(writer http.ResponseWriter, userId uuid.UUID) *SSEConn {
	return &SSEConn{
		writer:     writer,
		userId:     userId.String(),
		eventQueue: make(chan Event, 10),
	}
}

func (s *SSEConn) WriteEvent(event Event) error {
	eventId := uuid.New().String()
	eventData, err := event.EncodeEvent()
	if err != nil {
		return err
	}
	fmt.Fprintf(s.writer, "id: %s\nevent: %s\ndata: %s\n\n", eventId, event.EventType, eventData)
	return nil
}

func (em *EventManager) AddSSEConn(conn *SSEConn) {
	em.mu.Lock()
	defer em.mu.Unlock()
	if _, ok := em.sseConnMap[conn.userId]; !ok {
		em.sseConnMap[conn.userId] = make([]*SSEConn, 0)
	}
	em.sseConnMap[conn.userId] = append(em.sseConnMap[conn.userId], conn)
}

func (em *EventManager) RemoveSSEConn(conn *SSEConn) {
	em.mu.Lock()
	defer em.mu.Unlock()
	if _, ok := em.sseConnMap[conn.userId]; !ok {
		return
	}
	for i, c := range em.sseConnMap[conn.userId] {
		if c == conn {
			em.sseConnMap[conn.userId] = append(em.sseConnMap[conn.userId][:i], em.sseConnMap[conn.userId][i+1:]...)
			break
		}
	}
}

func (em *EventManager) SSEHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("user_id")
	parsedUserId, err := uuid.Parse(userId)
	if userId == "" || parsedUserId == uuid.Nil || err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	conn := NewSSEConn(w, parsedUserId)
	em.AddSSEConn(conn)
	defer em.RemoveSSEConn(conn)

	for {
		select {
		case event := <-conn.eventQueue:
			log.Println("event:", event)
			log.Println(event.Data)
			userId := event.UserId
			if userId.ValueOrZero() == conn.userId {
				err := conn.WriteEvent(event)
				if err != nil {
					slog.Error("SSEHandler: conn.WriteEvent:", "err", err)
				}
				w.(http.Flusher).Flush()
			}
		case <-r.Context().Done():
			return
		}
	}
}
