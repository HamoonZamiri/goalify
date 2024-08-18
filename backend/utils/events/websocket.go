package events

import (
	"goalify/utils/responses"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

const EVENT_BUFFER_SIZE = 256

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
				}
			}
		case <-r.Context().Done():
			return
		}
	}
}
