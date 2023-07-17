package ws

import (
	"encoding/json"
	"errors"
	"time"
)

type EventType uint8

const (
	Text               EventType = 2
	Swap               EventType = 4
	ParticipantHistory EventType = 6 // slice
	MessageHistory     EventType = 8 // slice
)

type EventReq struct {
	Payload string    `json:"payload"`
	Type    EventType `json:"type"`
}

func (event *EventReq) handleSwap(p *Participant, r *Room, h *Hub) error {
	var payload SwapMessageReq
	err := json.Unmarshal([]byte(event.Payload), &payload)
	if err != nil {
		return err
	}
	if payload.From != p.RoomSlug || payload.From != r.RoomSlug {
		return errors.New("failed to change rooms")
	}
	newRoom, ok := h.Rooms[payload.To]
	if !ok {
		return errors.New("failed to change rooms")
	}

	r.ParticipantMu.Lock()
	delete(r.Participants, p.Id)
	r.ParticipantMu.Unlock()
	r.CheckHealth(h) // check if it is needed to be killed

	p.HandleBailout(r, h)

	newRoom.AddParticipant(p)
	go p.HandleInit(newRoom, h) // handle it with new room
	// p.RoomSlug = payload.To
	// h.Register <- p TODO: use hub later

	return nil
}

func (event *EventReq) handleTextMessage(p *Participant, r *Room, h *Hub) *MessageResponse {
	r.ParticipantMu.Lock()
	defer r.ParticipantMu.Unlock()
	res := &MessageResponse{
		Payload: TextMessageRes{
			Id:      GenerateId(h),
			Content: event.Payload,
			Sender:  *p,
			Created: time.Now(),
		},
		Type: Text,
		Slug: p.RoomSlug,
	}
	return res
}
