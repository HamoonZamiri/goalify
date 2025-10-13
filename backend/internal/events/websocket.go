package events

import (
	"goalify/internal/responses"
	"log/slog"
	"net/http"
	"slices"

	"github.com/gorilla/websocket"
)

const EVENT_BUFFER_SIZE = 256

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
	userId     string
}

func newWebSocketConn(conn *websocket.Conn, userId string) *webSocketConn {
	return &webSocketConn{
		eventQueue: make(chan Event, EVENT_BUFFER_SIZE),
		conn:       conn,
		userId:     userId,
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

	userId := r.Header.Get("user_id")
	if userId == "" {
		responses.SendAPIError(
			w,
			r,
			http.StatusBadRequest,
			"user not authenticated",
			nil)
		return
	}

	wsConn := newWebSocketConn(conn, userId)
	em.SubscribeToUserEvents(userId, wsConn)
	defer func() {
		em.UnsubscribeFromUserEvents(userId, wsConn)
		wsConn.conn.Close()
	}()

	for {
		select {
		case event := <-wsConn.eventQueue:
			userId := event.UserId
			if userId.ValueOrZero() == wsConn.userId {
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
