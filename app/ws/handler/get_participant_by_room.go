package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/tastycrayon/whisker/app/ws"
)

func GetParticipantsByRoom(h *ws.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomSlug := c.QueryParam("roomSlug")
		if roomSlug == "" {
			return apis.NewBadRequestError("room not found", nil)
		}
		participants := make([]*ws.Participant, 0)
		if room, roomExist := h.Rooms[roomSlug]; roomExist && len(room.Participants) != 0 {
			for _, v := range room.Participants {
				participants = append(participants, v)
			}
			return c.JSON(http.StatusOK, participants)
		}
		return c.JSON(http.StatusOK, participants)
	}
}
