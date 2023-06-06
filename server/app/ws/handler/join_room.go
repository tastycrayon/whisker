package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/tastycrayon/go-chat/app/ws"
	"nhooyr.io/websocket"
)

// import (
// 	"fmt"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/websocket/v2"
// 	"github.com/nekonako/moechat/app/user/events/ws"
// )

// func JoinRoom(c echo.Context) error {
// 	r, w := c.Request(), c.Response().Writer
// 	wbsock, err := websocket.Accept(w, r, nil)
// 	if err != nil {
// 		return apis.NewApiError(http.StatusInternalServerError, "could not accept websocket connection", nil)
// 	}
// 	defer wbsock.Close(websocket.StatusInternalError, "the sky is falling")

// 	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
// 	defer cancel()

// 	var v ws.Message

// 	err = wsjson.Read(ctx, wbsock, &v)
// 	if err != nil {
// 		return apis.NewApiError(http.StatusBadRequest, "failed receiving body", nil)
// 	}

//		log.Printf("received: %v", v)
//		wbsock.Close(websocket.StatusNormalClosure, "")
//		return nil
//	}

func JoinRoom(h *ws.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		opts := &websocket.AcceptOptions{
			OriginPatterns: []string{"localhost:5173"},
		}
		// get user
		authRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		fmt.Println("knocked", authRecord.Id)
		fmt.Println("knocked", authRecord.Get("avatar"))
		r, w := c.Request(), c.Response()
		wbsock, err := websocket.Accept(w, r, opts)
		if err != nil {
			fmt.Println(err)
			return apis.NewApiError(http.StatusInternalServerError, "could not accept websocket connection", nil)
		}

		roomId := c.PathParam("roomId")
		// clientId := c.QueryParam("userId")
		// username := c.QueryParam("username")
		avatar := authRecord.Get("avatar").(string)
		clientId := authRecord.Id
		username := authRecord.Username()

		client := &ws.Participant{
			ClientId: clientId,
			Username: username,
			Avatar:   avatar,
			RoomId:   roomId,
			Conn:     wbsock,
			Message:  make(chan *ws.Post, 10),
		}
		fmt.Println("client", client)
		// m := ws.Message{
		// 	Message:  "A new user has joined the room.",
		// 	ClientId: client.ClientId,
		// 	RoomId:   client.RoomId,
		// 	Username: username,
		// }
		m := ws.Post{
			RoomId: client.RoomId,
			Message: ws.Message{
				Id:      "not_avaliable_yet",
				Content: "A new user has joined the room.",
				Type:    ws.Welcome,
			},
			Sender: ws.MessageSender{
				ClientId: clientId,
				Username: username,
				RoomId:   client.RoomId,
				Avatar:   avatar,
			},
			Created: time.Now(),
		}
		fmt.Println("mess", m)
		h.Register <- client
		h.Broadcast <- &m

		go client.WriteMessage()
		client.ReadMessage(h, r)

		return nil
	}
}
