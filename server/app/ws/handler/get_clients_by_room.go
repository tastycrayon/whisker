package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/tastycrayon/go-chat/app/ws"
)

func GetParticipantByRoom(h *ws.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomId := c.PathParam("roomSlug")
		participants := make([]*ws.Participant, 0)
		if room, roomExist := h.Rooms[roomId]; roomExist && len(room.Participants) != 0 {
			for _, v := range room.Participants {
				participants = append(participants, v)
			}
			return c.JSON(http.StatusOK, participants)
		}
		return c.JSON(http.StatusOK, participants)
	}
}
