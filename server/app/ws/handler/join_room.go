package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/tastycrayon/go-chat/app/ws"
	"nhooyr.io/websocket"
)

func JoinRoom(h *ws.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomSlug := c.PathParam("roomSlug")
		_, found := h.Rooms[roomSlug]
		if !found {
			return apis.NewApiError(http.StatusInternalServerError, "room not found", nil)
		}
		opts := &websocket.AcceptOptions{
			OriginPatterns: []string{"localhost:5173"},
		}
		// get user
		authRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		r, w := c.Request(), c.Response()
		wbsock, err := websocket.Accept(w, r, opts)
		if err != nil {
			fmt.Println(err)
			return apis.NewApiError(http.StatusInternalServerError, "could not accept websocket connection", nil)
		}
		avatar := authRecord.Get("avatar").(string)
		participantId := authRecord.Id
		username := authRecord.Username()

		participant := &ws.Participant{
			Id:       participantId,
			Username: username,
			Avatar:   avatar,
			RoomSlug: roomSlug,
			Conn:     wbsock,
			Message:  make(chan *ws.MessageFeed, h.SubscriberMessageBuffer),
			CloseSlow: func() {
				wbsock.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
			},
		}
		h.Register <- participant
		go participant.WriteMessage(r.Context())
		participant.ReadMessage(h, r) // subscriber

		return nil
	}
}
