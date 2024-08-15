package events

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	SSE_BUFFER_SIZE = 10
)

type SSEConn struct {
	eventQueue chan Event
	writer     http.ResponseWriter
	userId     string
}

func newSSEConn(writer http.ResponseWriter, userId string) *SSEConn {
	return &SSEConn{
		writer:     writer,
		userId:     userId,
		eventQueue: make(chan Event, SSE_BUFFER_SIZE),
	}
}

func (s *SSEConn) writeEvent(event Event) error {
	eventId := uuid.New().String()
	eventData, err := event.EncodeEvent()
	if err != nil {
		return err
	}
	fmt.Fprintf(s.writer, "id: %s\nevent: %s\ndata: %s\n\n", eventId, event.EventType, eventData)
	return nil
}

func (em *EventManager) addSSEConn(conn *SSEConn) {
	em.mu.Lock()
	defer em.mu.Unlock()
	if _, ok := em.sseConnMap[conn.userId]; !ok {
		em.sseConnMap[conn.userId] = make([]*SSEConn, 0)
	}
	em.sseConnMap[conn.userId] = append(em.sseConnMap[conn.userId], conn)
}

func (em *EventManager) removeSSEConn(conn *SSEConn) {
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
	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	conn := newSSEConn(w, userId)
	em.addSSEConn(conn)
	defer em.removeSSEConn(conn)

	for {
		select {
		case event := <-conn.eventQueue:
			userId := event.UserId
			if userId.ValueOrZero() == conn.userId {
				err := conn.writeEvent(event)
				if err != nil {
					slog.Error("SSEHandler: conn.WriteEvent:", "err", err)
				}
				w.(http.Flusher).Flush()
			}
		case <-time.After(10 * time.Second):
			// Send a keep-alive message
			fmt.Fprintf(w, ": keep-alive\n\n")
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
}
