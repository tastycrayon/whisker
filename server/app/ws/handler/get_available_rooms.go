package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/tastycrayon/go-chat/app/ws"
)

func GetAvailableRooms(h *ws.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		rooms := make([]ws.Room, 0)

		for _, room := range h.Rooms {
			rooms = append(rooms, *room)
		}
		return c.JSON(http.StatusOK, rooms)
	}
}
