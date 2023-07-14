package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/tastycrayon/pb-svelte-chatapp/app/ws"
)

type CreateRoomReq struct {
	RoomName        string `json:"name"`
	RoomSlug        string `json:"slug"`
	RoomDescription string `json:"description"`
}

func CreateRoom(h *ws.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		// authRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

		roomReq := new(CreateRoomReq)

		if err := c.Bind(roomReq); err != nil {
			return apis.NewApiError(http.StatusBadRequest, "could not open room with given parameters", nil)
		}

		if _, roomFound := h.Rooms[roomReq.RoomSlug]; roomFound {
			return apis.NewApiError(http.StatusBadRequest, "room already exists", nil)
		}

		// h.Rooms[roomReq.RoomSlug] = ws.NewRoom("-", roomReq.RoomSlug, roomReq.RoomName, roomReq.RoomDescription, authRecord.Id, ws.PersonalRoom)

		return c.JSON(http.StatusOK, roomReq)
	}
}
