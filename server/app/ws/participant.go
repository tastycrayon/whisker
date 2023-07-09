package ws

import (
	"context"
	"errors"
	"fmt"
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

func NewParticipant(h *Hub, conn *websocket.Conn, id, username, slug, avatar string) *Participant {
	return &Participant{
		Conn:     conn,
		Id:       id,
		Username: username,
		Avatar:   avatar,
		RoomSlug: slug,
		Message:  make(chan *MessageFeed, h.SubscriberMessageBuffer),
		CloseSlow: func() {
			conn.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		},
	}
}

// from webscoket Connections to Hub
func (p *Participant) ReadMessage(h *Hub, ctx context.Context) {
	defer func() {
		// send participant to unregister and disconnect him | will run on error(such as disconnect)
		h.Unregister <- p // when participant is not receiving messages
		close(p.Message)  // chan - stop receiving messages
		fmt.Println("at reader: close(p.Message)")
		p.Conn.Close(websocket.StatusNormalClosure, "either participant disconnected or message read failed")
		ctx.Done() // TODO
	}()

	for {
		var m MessageReq
		if err := wsjson.Read(ctx, p.Conn, &m); err != nil {
			fmt.Println("failed at message reader", err)
			if errors.Is(err, websocket.CloseError{}) {
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
func (p *Participant) WriteMessage(ctx context.Context) error {
	defer func() {
		fmt.Println("Connection was closed")
		ctx.Done()
		fmt.Println("at writer: ctx.Done()")
	}()
	for {
		select {
		case message := <-p.Message:
			err := writeTimeout(ctx, time.Second*3, p.Conn, message)
			if err != nil {
				fmt.Println("error at writer", err)
				go p.CloseSlow()
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (p *Participant) WriteCustomMessage(m interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	wsjson.Write(ctx, p.Conn, m)
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
