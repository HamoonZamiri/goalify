package events

import (
	"net/http"

	"github.com/google/uuid"
)

type SSEConn struct {
	writer http.ResponseWriter
	userId uuid.UUID
}

func NewSSEConn(writer http.ResponseWriter, userId uuid.UUID) *SSEConn {
	return &SSEConn{
		writer: writer,
		userId: userId,
	}
}

func (s *SSEConn) WriteEvent(event []byte) error {
	_, err := s.writer.Write(event)
	return err
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
	parsedUserId, _ := uuid.Parse(userId)
	if userId == "" || parsedUserId == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	conn := NewSSEConn(w, parsedUserId)
	em.AddSSEConn(conn)
	defer em.RemoveSSEConn(conn)

	for {
		select {
		case event := <-em.eventQueue:
			userId := event.UserId
			encoded, err := event.EncodeEvent()
			if userId.ValueOrZero() == conn.userId && err == nil {
				conn.WriteEvent(encoded)
			}
		case <-r.Context().Done():
			return
		}
	}
}
