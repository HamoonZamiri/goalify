package events

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

const (
	SSEBufferSize = 10
)

type SSEConn struct {
	eventQueue chan Event
	writer     http.ResponseWriter
	userID     string
}

func newSSEConn(writer http.ResponseWriter, userID string) *SSEConn {
	return &SSEConn{
		writer:     writer,
		userID:     userID,
		eventQueue: make(chan Event, SSEBufferSize),
	}
}

func (s *SSEConn) HandleEvent(event Event) {
	s.eventQueue <- event
}

func (s *SSEConn) writeEvent(event Event) error {
	eventID := uuid.New().String()
	eventData, err := event.EncodeEvent()
	if err != nil {
		return err
	}
	slog.Info("SSE writing event",
		slog.String("eventId", eventID),
		slog.String("eventType", event.EventType),
		slog.String("userId", s.userID))
	_, err = fmt.Fprintf(
		s.writer,
		"id: %s\nevent: %s\ndata: %s\n\n",
		eventID,
		event.EventType,
		eventData,
	)
	return err
}

func (em *EventManager) SSEHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("user_id")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	conn := newSSEConn(w, userID)
	em.SubscribeToUserEvents(conn.userID, conn)
	defer em.UnsubscribeFromUserEvents(conn.userID, conn)

	// send an initial event for browser connection
	err := conn.writeEvent(NewEventWithUserID(SSEConnected, nil, conn.userID))
	if err != nil {
		slog.Error("SSEHandler: conn.WriteEvent:", "err", err)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		slog.Error("SSEHandler: http.ResponseWriter does not implement http.Flusher")
		return
	}
	flusher.Flush()

	for {
		select {
		case event := <-conn.eventQueue:
			slog.Info("SSE received event from queue",
				slog.String("eventType", event.EventType),
				slog.String("eventUserId", event.UserID.ValueOrZero()),
				slog.String("connUserId", conn.userID))
			userID := event.UserID
			if userID.ValueOrZero() == conn.userID {
				err := conn.writeEvent(event)
				if err != nil {
					slog.Error("SSEHandler: conn.WriteEvent:", "err", err)
				}
				flusher.Flush()
			} else {
				slog.Warn("SSE event userId mismatch, skipping",
					slog.String("eventUserId", userID.ValueOrZero()),
					slog.String("connUserId", conn.userID))
			}
		case <-r.Context().Done():
			return
		}
	}
}
