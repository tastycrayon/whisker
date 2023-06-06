package ws

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// 	"github.com/gofiber/websocket/v2"

type Participant struct {
	Conn     *websocket.Conn
	ClientId string `json:"clientId"`
	Username string `json:"username"`
	RoomId   string `json:"roomId"`
	Avatar   string `json:"avatar"`
	Message  chan *Post
}

// const (
// 	// Time allowed to write a message to the peer.
// 	writeWait = 10 * time.Second

// 	// Time allowed to read the next pong message from the peer.
// 	pongWait = 60 * time.Second

// 	// Send pings to peer with this period. Must be less than pongWait.
// 	pingPeriod = (pongWait * 9) / 10

// 	// Maximum message size allowed from peer.
// 	maxMessageSize = 512
// )

// from webscoket Connections to Hub
func (c *Participant) ReadMessage(h *Hub, r *http.Request) {
	defer func() {
		// send client to unregister and disconnect him
		h.Unregister <- c
		c.Conn.Close(websocket.StatusNormalClosure, "")

	}()

	for {
		// _, m, err := c.Conn.ReadMessage() // from user get new msg
		var m MessageReq
		// ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		// defer cancel() // read message for 10 second
		if err := wsjson.Read(context.Background(), c.Conn, &m); err != nil {
			fmt.Println("failed at reader", err)
			if strings.Contains(err.Error(), "websocket: close") {
				fmt.Println("close Connection")
			}
			break
		}

		fmt.Println("message from client (first message)", m)
		// mesaage := Message{
		// 	Message:  string(v.Message),
		// 	ClientId: c.ClientId,
		// 	RoomId:   c.RoomId,
		// 	Username: c.Username,
		// }
		post := &Post{
			RoomId: c.RoomId,
			Message: Message{
				Id:      GenerateId(),
				Content: "An user has left the room.",
				Type:    Bailout,
			},
			Sender: MessageSender{
				ClientId: c.ClientId,
				Username: c.Username,
				RoomId:   c.RoomId,
				Avatar:   c.Avatar,
			},
			Created: time.Now(),
		}
		h.Broadcast <- post
	}
}

// from Hub to websocket Connection
func (c *Participant) WriteMessage() {
	defer func() {
		fmt.Println("Connection was closed")
	}()
	for {
		message, ok := <-c.Message
		if !ok { // close(c.Message) will happen when client unsubscribes
			return
		}
		// c.Conn.WriteJSON(message)
		// ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		// defer cancel()
		wsjson.Write(context.Background(), c.Conn, message)
	}
}
