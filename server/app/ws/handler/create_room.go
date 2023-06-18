package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/tastycrayon/go-chat/app/ws"
)

type CreateRoomReq struct {
	RoomName string `json:"roomname"`
	RoomId   string `json:"roomId"`
}

func CreateRoom(h *ws.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		authRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

		room := new(CreateRoomReq)

		if err := c.Bind(room); err != nil {
			panic(err)
		}
		h.Rooms[room.RoomId] = ws.NewRoom(room.RoomId, room.RoomName, authRecord.Id, ws.PersonalRoom)

		return c.JSON(http.StatusOK, room)
	}
}
