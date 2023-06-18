package ws

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
)

type MessageType uint8

const (
	Welcome MessageType = 2
	Ping    MessageType = 4 // slice
	Text    MessageType = 8
	Bailout MessageType = 16
	Swap    MessageType = 32
	History MessageType = 64 // slice
)

type MessageReq struct {
	Content string      `json:"content"`
	Type    MessageType `json:"messageType"`
	From    string      `json:"from"`
	To      string      `json:"to"`
}

type MessageSender struct {
	Id       string `json:"id"`
	Username string `json:"userName"`
	RoomSlug string `json:"roomSlug"`
	Avatar   string `json:"avatar"`
}

type MessageFeed struct {
	SID      snowflake.ID  `json:"sid"`
	Content  interface{}   `json:"content"`
	Type     MessageType   `json:"messageType"`
	Sender   MessageSender `json:"sender"`
	RoomSlug string        `json:"roomSlug"`
	Created  time.Time     `json:"created"`
}

func GenerateId(h *Hub) snowflake.ID {
	return h.Node.Generate()
}

func GenerateUserLeftMessage(p Participant, id snowflake.ID) *MessageFeed {
	return &MessageFeed{
		SID:      id,
		RoomSlug: p.RoomSlug,
		Type:     Bailout,
		Content:  fmt.Sprintf("%v has left the room.", p.Username),
		Sender: MessageSender{
			Id:       p.Id,
			Username: p.Username,
			RoomSlug: p.RoomSlug,
			Avatar:   p.Avatar,
		},
		Created: time.Now(),
	}
}
func GenerateUserJoinedMessage(p Participant, id snowflake.ID) *MessageFeed {
	return &MessageFeed{
		SID:      id,
		RoomSlug: p.RoomSlug,
		Content:  fmt.Sprintf("%v has joined the room.", p.Username),
		Type:     Welcome,
		Sender: MessageSender{
			Id:       p.Id,
			Username: p.Username,
			RoomSlug: p.RoomSlug,
			Avatar:   p.Avatar,
		},
		Created: time.Now(),
	}
}

func GenerateCustomMessage(p Participant, id snowflake.ID, t MessageType, content interface{}) *MessageFeed {
	return &MessageFeed{
		SID:      id,
		RoomSlug: p.RoomSlug,
		Content:  content,
		Type:     t,
		Sender: MessageSender{
			Id:       p.Id,
			Username: p.Username,
			RoomSlug: p.RoomSlug,
			Avatar:   p.Avatar,
		},
		Created: time.Now(),
	}
}
