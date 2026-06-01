package events

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	SSEBufferSize = 10
	// sseWriteTimeout bounds how long a single write/flush to the client may
	// block. The server sets no global WriteTimeout (the stream is intentionally
	// long-lived), so without a per-write deadline a half-dead connection — a
	// suspended browser tab, or a proxy-dropped socket — would hang the writer
	// indefinitely and, via the shared publisher, stall delivery for everyone.
	sseWriteTimeout = 10 * time.Second
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

// HandleEvent queues an event for delivery without blocking the caller. The
// publisher fans out to every subscriber on a single goroutine, so a blocking
// send here would stall delivery for all connections. A full buffer means this
// client is too slow (or its socket is wedged); we drop the event rather than
// freeze the system — the client reconciles its state on the next reconnect.
func (s *SSEConn) HandleEvent(event Event) {
	select {
	case s.eventQueue <- event:
	default:
		slog.Warn("SSE buffer full, dropping event",
			slog.String("eventType", event.EventType),
			slog.String("userId", s.userID))
	}
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
	// no-transform stops Cloudflare and other proxies from buffering or
	// compressing the stream, which would otherwise clump events together
	// instead of flushing them as they happen.
	w.Header().Set("Cache-Control", "no-cache, no-transform")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	conn := newSSEConn(w, userID)
	em.SubscribeToUserEvents(conn.userID, conn)
	defer em.UnsubscribeFromUserEvents(conn.userID, conn)

	flusher, ok := w.(http.Flusher)
	if !ok {
		slog.Error("SSEHandler: http.ResponseWriter does not implement http.Flusher")
		return
	}

	rc := http.NewResponseController(w)
	// refreshDeadline arms the next write with a fresh deadline so a stalled
	// client fails fast instead of hanging. Connections that don't support
	// deadlines (e.g. the test recorder) degrade to the previous behaviour.
	refreshDeadline := func() {
		if err := rc.SetWriteDeadline(time.Now().Add(sseWriteTimeout)); err != nil &&
			!errors.Is(err, errors.ErrUnsupported) {
			slog.Warn("SSEHandler: SetWriteDeadline:", "err", err)
		}
	}

	// tell the browser to reconnect after 3 seconds if the stream drops
	refreshDeadline()
	if _, err := fmt.Fprintf(w, "retry: 3000\n\n"); err != nil {
		slog.Error("SSEHandler: conn.WriteEvent:", "err", err)
		return
	}

	// send an initial event for browser connection
	refreshDeadline()
	if err := conn.writeEvent(NewEventWithUserID(SSEConnected, nil, conn.userID)); err != nil {
		slog.Error("SSEHandler: conn.WriteEvent:", "err", err)
		return
	}
	flusher.Flush()

	heartbeat := time.NewTicker(30 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case event := <-conn.eventQueue:
			slog.Info("SSE received event from queue",
				slog.String("eventType", event.EventType),
				slog.String("eventUserId", event.UserID.ValueOrZero()),
				slog.String("connUserId", conn.userID))
			if event.UserID.ValueOrZero() == conn.userID {
				refreshDeadline()
				if err := conn.writeEvent(event); err != nil {
					slog.Error("SSEHandler: conn.WriteEvent:", "err", err)
					return
				}
				flusher.Flush()
			} else {
				slog.Warn("SSE event userId mismatch, skipping",
					slog.String("eventUserId", event.UserID.ValueOrZero()),
					slog.String("connUserId", conn.userID))
			}
		case <-heartbeat.C:
			refreshDeadline()
			if _, err := fmt.Fprintf(w, ": heartbeat\n\n"); err != nil {
				slog.Error("SSEHandler: heartbeat write failed", "err", err)
				return
			}
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
