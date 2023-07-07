package ws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Participant struct {
	Conn      *websocket.Conn   `json:"-"`
	Id        string            `json:"id"`
	Username  string            `json:"username"`
	RoomSlug  string            `json:"roomSlug"`
	Avatar    string            `json:"avatar"`
	Message   chan *MessageFeed `json:"-"`
	CloseSlow func()            `json:"-"`
}

// from webscoket Connections to Hub
func (p *Participant) ReadMessage(h *Hub, ctx context.Context) {
	defer func() {
		// send participant to unregister and disconnect him | will run on error(such as disconnect)
		h.Unregister <- p // when participant is not receiving messages
		close(p.Message)  // chan - stop receiving messages
		p.Conn.Close(websocket.StatusNormalClosure, "either participant disconnected or message read failed")

	}()

	for {
		var m MessageReq
		if err := wsjson.Read(ctx, p.Conn, &m); err != nil {
			fmt.Println("failed at message reader", err)
			if strings.Contains(err.Error(), "websocket: close") {
				fmt.Println("closing connection for websocket")
			}
			break
		}
		// intercept messaage
		if m.Type == Swap {
			if _, ok := h.Rooms[m.From]; ok && m.From != p.RoomSlug {
				fmt.Println("swap parse error")
				continue
			}
			h.Unregister <- p // unregister from existing
			p.RoomSlug = m.To // change room id
			h.Register <- p   // re-register to new room
			continue
		}
		// intercept messaage

		message := &MessageFeed{
			SID:      GenerateId(h),
			RoomSlug: p.RoomSlug,
			Content:  m.Content,
			Type:     m.Type,
			Sender:   *p,
			Created:  time.Now(),
		}
		h.Broadcast <- message
	}
}

// from Hub to websocket Connection
func (c *Participant) WriteMessage(ctx context.Context) error {
	defer func() {
		fmt.Println("Connection was closed")
	}()
	for {
		select {
		case message := <-c.Message:
			err := writeTimeout(ctx, time.Second*5, c.Conn, message)
			if err != nil {
				fmt.Println("er", err)
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *Participant) WriteCustomMessage(m interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	wsjson.Write(ctx, c.Conn, m)
}

func GetAllParticipants(room *Room) []*Participant {
	participants := make([]*Participant, 0)
	for _, p := range room.Participants {
		participants = append(participants, p)
	}
	return participants
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return wsjson.Write(ctx, c, msg)
}
