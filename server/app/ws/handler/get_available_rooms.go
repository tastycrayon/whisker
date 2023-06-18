package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/tastycrayon/go-chat/app/ws"
	"github.com/tastycrayon/go-chat/model"
)

type Res struct {
	*model.BaseResponse
	Data *[]RoomList `json:"data"`
}

type RoomList struct {
	RoomSlug  string      `json:"roomSlug"`
	RoomName  string      `json:"roomName"`
	Type      ws.RoomType `json:"roomType"`
	CreatedBy string      `json:"createdBy"`
}

func GetAvailableRooms(h *ws.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		rooms := make([]RoomList, 0)

		for _, room := range h.Rooms {
			rooms = append(rooms, RoomList{
				RoomSlug:  room.RoomSlug,
				RoomName:  room.RoomName,
				Type:      room.Type,
				CreatedBy: room.CreatedBy,
			})
		}
		return c.JSON(http.StatusOK, rooms)
	}
}
