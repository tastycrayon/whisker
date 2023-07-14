package ws

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
)

// response for text messages
type TextMessageRes struct {
	Id      snowflake.ID `json:"id"`
	Content interface{}  `json:"content"`
	Sender  Participant  `json:"sender"`
	Created time.Time    `json:"created"`
}

// request messages for changing room
type SwapMessageReq struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// full response of messages
type MessageResponse struct {
	Payload interface{} `json:"payload"`
	Type    EventType   `json:"type"`
	Slug    string      `json:"slug"`
}

func GenerateId(h *Hub) snowflake.ID {
	return h.Node.Generate()
}
func GenerateUserLeftMessage(p *Participant, id snowflake.ID) *MessageResponse {
	return &MessageResponse{
		Payload: TextMessageRes{
			Id:      id,
			Content: fmt.Sprintf("%v has left the room.", p.Username),
			Sender:  *p,
			Created: time.Now(),
		},
		Type: Text,
		Slug: p.RoomSlug,
	}
}
func GenerateUserJoinedMessage(p *Participant, id snowflake.ID) *MessageResponse {
	return &MessageResponse{
		Payload: TextMessageRes{
			Id:      id,
			Content: fmt.Sprintf("%v has joined the room.", p.Username),
			Sender:  *p,
			Created: time.Now(),
		},
		Type: Text,
		Slug: p.RoomSlug,
	}
}
func GenerateCustomMessage(p *Participant, id snowflake.ID, t EventType, payload interface{}) *MessageResponse {
	return &MessageResponse{
		Payload: payload,
		Type:    t,
		Slug:    p.RoomSlug,
	}
}
