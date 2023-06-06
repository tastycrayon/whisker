package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/tastycrayon/go-chat/app/ws"
	"github.com/tastycrayon/go-chat/model/api"
)

type PostList struct {
	*api.BaseResponse
	Data []ws.Post `json:"data"`
}

func GetPostByRoom(h *ws.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		var posts []ws.Post
		roomId := c.PathParam("roomId")
		fmt.Println(roomId)

		room, doesExist := h.Rooms[roomId]
		if !doesExist || room.PostsQueue.IsEmpty() {
			res := PostList{
				BaseResponse: &api.BaseResponse{
					Success: true,
					Code:    200,
					Message: "no client",
				},
				Data: make([]ws.Post, 0),
			}
			return c.JSON(http.StatusOK, res)
		}

		for {
			p, ok := room.PostsQueue.Next()
			if !ok {
				break
			}
			posts = append(posts, *p)
		}

		res := PostList{
			BaseResponse: &api.BaseResponse{
				Success: true,
				Code:    200,
				Message: "success get clients this room",
			},
			Data: posts,
		}

		return c.JSON(http.StatusOK, res)
	}
}
