package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/tastycrayon/go-chat/app/ws"
	"github.com/tastycrayon/go-chat/model/api"
)

// import (
// 	"fmt"

// 	"github.com/labstack/echo"
// 	"github.com/tastycrayon/go-chat/model/api"
// 	"golang.org/x/net/websocket"
// )

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/nekonako/moechat/app/user/events/ws"

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

// func GetAvailableRooms(c *fiber.Ctx, h *ws.Hub) error {

// 	rooms := make([]RoomList, 0)

// 	for _, room := range h.Rooms {
// 		rooms = append(rooms, RoomList{
// 			RoomName: room.RoomName,
// 			RoomId:   room.RoomId,
// 		})
// 	}

// 	res := Res{
// 		BaseResponse: &api.BaseResponse{
// 			Success: true,
// 			Code:    200,
// 			Message: "success get rooms",
// 		},
// 		Data: &rooms,
// 	}

// 	return c.JSON(&res)
// }
