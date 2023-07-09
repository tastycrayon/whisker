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
		r, w := c.Request(), c.Response()

		roomSlug := c.PathParam("roomSlug")

		if _, roomFound := h.Rooms[roomSlug]; !roomFound {
			return apis.NewApiError(http.StatusNotFound, "room not found", nil)
		}
		// get user
		authRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

		opts := &websocket.AcceptOptions{
			OriginPatterns: []string{"127.0.0.1:5173", "127.0.0.1:8090", "127.0.0.1:8080", "127.0.0.1:80"},
		}

		conn, err := websocket.Accept(w, r, opts)
		if err != nil {
			fmt.Println(err)
			return apis.NewApiError(http.StatusInternalServerError, "could not accept websocket connection", nil)
		}
		avatar := authRecord.Get("avatar").(string)
		participantId := authRecord.Id
		username := authRecord.Username()

		participant := ws.NewParticipant(h, conn, participantId, username, roomSlug, avatar)
		h.Register <- participant

		// be ready to write new messages to this user
		go participant.WriteMessage(r.Context())

		// read messages from this user and forward the message to channel
		participant.ReadMessage(h, r.Context()) // subscriber

		return nil
	}
}
