package events

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

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

func newSSEConn(writer http.ResponseWriter, userId uuid.UUID) *SSEConn {
	return &SSEConn{
		writer:     writer,
		userId:     userId.String(),
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
	parsedUserId, err := uuid.Parse(userId)
	if userId == "" || parsedUserId == uuid.Nil || err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	conn := newSSEConn(w, parsedUserId)
	em.addSSEConn(conn)
	defer em.removeSSEConn(conn)

	for {
		select {
		case event := <-conn.eventQueue:
			log.Println("event:", event)
			log.Println(event.Data)
			userId := event.UserId
			if userId.ValueOrZero() == conn.userId {
				err := conn.writeEvent(event)
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
