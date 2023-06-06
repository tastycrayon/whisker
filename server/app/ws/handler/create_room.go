package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/tastycrayon/go-chat/app/ws"
)

// import (
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/nekonako/moechat/app/user/events/ws"
// )

type CreateRoomReq struct {
	RoomName string `json:"roomname"`
	RoomId   string `json:"roomId"`
}

func CreateRoom(h *ws.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get user
		authRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

		room := new(CreateRoomReq)

		if err := c.Bind(room); err != nil {
			panic(err)
		}

		h.Rooms[room.RoomId] = &ws.Room{
			RoomId:    room.RoomId,
			RoomName:  room.RoomName,
			Clients:   make(map[string]*ws.Participant),
			CreatedBy: authRecord.Id,
		}

		return c.JSON(http.StatusOK, room)
	}
}

// func CreateRoom(c *fiber.Ctx, h *ws.Hub) error {

// 	room := new(MyRoom)

// 	if err := c.BodyParser(room); err != nil {
// 		panic(err)
// 	}

// 	h.Rooms[room.RoomId] = &ws.Room{
// 		RoomId:   room.RoomId,
// 		RoomName: room.RoomName,
// 		Clients:  make(map[string]*ws.Client),
// 	}

// 	return c.JSON(room)

// }
