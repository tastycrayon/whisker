package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/tastycrayon/go-chat/app/ws"
	"github.com/tastycrayon/go-chat/model/api"
)

type Res struct {
	*api.BaseResponse
	Data *[]RoomList `json:"data"`
}

type RoomList struct {
	RoomName string `json:"roomName"`
	RoomId   string `json:"roomId"`
}

func GetAvailableRooms(h *ws.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		rooms := make([]RoomList, 0)

		for _, room := range h.Rooms {
			rooms = append(rooms, RoomList{
				RoomName: room.RoomName,
				RoomId:   room.RoomId,
			})
		}

		res := Res{
			BaseResponse: &api.BaseResponse{
				Success: true,
				Code:    200,
				Message: "success get rooms",
			},
			Data: &rooms,
		}

		return c.JSON(http.StatusOK, res)
	}
}
