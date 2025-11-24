package events

import (
	"goalify/internal/responses"
	"log/slog"
	"net/http"
	"slices"

	"github.com/gorilla/websocket"
)

const EventBufferSize = 256

var allowedOrigins []string = []string{
	"http://localhost:5173",
}

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// by default the origin for the vue frontend was not allowed
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return slices.Contains(allowedOrigins, origin)
	},
}

type webSocketConn struct {
	eventQueue chan Event
	conn       *websocket.Conn
	userID     string
}

func newWebSocketConn(conn *websocket.Conn, userID string) *webSocketConn {
	return &webSocketConn{
		eventQueue: make(chan Event, EventBufferSize),
		conn:       conn,
		userID:     userID,
	}
}

func (w *webSocketConn) HandleEvent(event Event) {
	w.eventQueue <- event
}

func (em *EventManager) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("WebSocketHandler: upgrader.Upgrade:", "err", err)
		responses.SendAPIError(
			w,
			r,
			http.StatusBadRequest,
			"error upgrading socket connection",
			nil)
		return
	}

	userID := r.Header.Get("user_id")
	if userID == "" {
		responses.SendAPIError(
			w,
			r,
			http.StatusBadRequest,
			"user not authenticated",
			nil)
		return
	}

	wsConn := newWebSocketConn(conn, userID)
	em.SubscribeToUserEvents(userID, wsConn)
	defer func() {
		em.UnsubscribeFromUserEvents(userID, wsConn)
		err := wsConn.conn.Close()
		if err != nil {
			slog.Error("WebSocketHandler: wsConn.conn.Close:", "err", err)
		}
	}()

	for {
		select {
		case event := <-wsConn.eventQueue:
			userID := event.UserID
			if userID.ValueOrZero() == wsConn.userID {
				err := conn.WriteJSON(event)
				if err != nil {
					slog.Error("WebSocketHandler: conn.WriteJSON:", "err", err)
					return
				}
			}
		case <-r.Context().Done():
			return
		}
	}
}
