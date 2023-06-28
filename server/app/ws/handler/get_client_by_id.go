package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/tastycrayon/go-chat/app/ws"
)

func GetParticipantById(pb *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		particpantId := c.PathParam("participantId")
		user, err := pb.Dao().FindRecordById("users", particpantId)
		if err != nil {
			return apis.NewApiError(http.StatusBadRequest, "could not find user with given parameters", nil)
		}
		res := &ws.Participant{
			Id:       user.Id,
			Username: user.GetString("username"),
			Avatar:   user.GetString("avatar"),
			RoomSlug: "",
		}
		return c.JSON(http.StatusOK, res)
	}
}
