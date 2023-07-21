package ws

import (
	"context"
	"fmt"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Participant struct {
	Conn      *websocket.Conn            `json:"-"`
	Id        string                     `json:"id"`
	Username  string                     `json:"username"`
	RoomSlug  string                     `json:"roomSlug"`
	Avatar    string                     `json:"avatar"`
	Message   chan *MessageResponse      `json:"-"`
	CloseSlow func(conn *websocket.Conn) `json:"-"`
}

type ParticipantList map[string]*Participant

func NewParticipant(h *Hub, conn *websocket.Conn, id, username, slug, avatar string) *Participant {
	return &Participant{
		Conn:     conn,
		Id:       id,
		Username: username,
		Avatar:   avatar,
		RoomSlug: slug,
		Message:  make(chan *MessageResponse, h.SubscriberMessageBuffer),
		CloseSlow: func(conn *websocket.Conn) {
			conn.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		},
	}
}

// from webscoket Connections to Hub
func (p *Participant) ReadMessage(ctx context.Context, h *Hub) {
	defer func() {
		fmt.Println("at reader")
		h.Unregister <- p
	}()

	for {
		var event EventReq
		if err := wsjson.Read(ctx, p.Conn, &event); err != nil {
			switch websocket.CloseStatus(err) {
			case websocket.StatusGoingAway:
			case websocket.StatusAbnormalClosure:
			default:
				fmt.Println("failed at message reader: ", err)
			}
			break
		}
		room, ok := h.Rooms[p.RoomSlug]
		if !ok {
			break
		}
		// intercept messaage
		switch event.Type {
		case Text:
			msg := event.handleTextMessage(p, room, h)
			h.Broadcast <- msg
		case Swap:
			err := event.handleSwap(p, room, h)
			if err != nil {
				continue
			}
		default:
			fmt.Println("unknown message/event type")
		}
	}
}

// from Hub to websocket Connection
func (p *Participant) WriteMessage(ctx context.Context, h *Hub) error {
	defer func() {
		fmt.Println("at writer")
		h.Unregister <- p
	}()

	for message := range p.Message {
		err := writeTimeout(ctx, time.Second*2, p.Conn, message)
		if err != nil {
			fmt.Println("error at writer", err)
			return err
		}
	}
	return nil
}

func (p *Participant) WriteCustomMessage(m interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	wsjson.Write(ctx, p.Conn, m)
}

func GetAllParticipants(room *Room) []*Participant {
	participants := make([]*Participant, 0, len(room.Participants))
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

// sends history, particpant-list and a welcome message
func (p *Participant) HandleInit(r *Room, h *Hub) {
	r.ParticipantMu.Lock()
	var pList, mList, welcomeMessage *MessageResponse
	// fetching participant list
	pList = GenerateCustomMessage(p, GenerateId(h), ParticipantHistory, GetAllParticipants(r))
	// fetch exisiting messages and post
	mList = GenerateCustomMessage(p, GenerateId(h), MessageHistory, r.LocalMessageQueue.GetInternalSlice())
	welcomeMessage = GenerateUserJoinedMessage(p, GenerateId(h))

	r.ParticipantMu.Unlock()

	p.WriteCustomMessage(pList)
	p.WriteCustomMessage(mList)

	h.Broadcast <- welcomeMessage
}

func (p *Participant) HandleBailout(r *Room, h *Hub) {
	r.ParticipantMu.Lock()
	defer r.ParticipantMu.Unlock()
	h.Broadcast <- GenerateUserLeftMessage(p, GenerateId(h))
}
