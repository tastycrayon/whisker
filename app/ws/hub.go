package ws

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"golang.org/x/time/rate"
)

type Hub struct {
	Broadcast  chan *MessageResponse
	Register   chan *Participant
	Unregister chan *Participant
	Rooms      map[string]*Room
	// snowflake generator node is used for making message id
	Node *snowflake.Node
	// SubscriberMessageBuffer controls the max number
	// of messages that can be queued for a subscriber
	// before it is kicked.
	//
	// Defaults to 16.
	SubscriberMessageBuffer int
	// publishLimiter controls the rate limit applied to the publish endpoint.
	//
	// Defaults to one publish every 100ms with a burst of 8.
	PublishLimiter *rate.Limiter
}

func (h *Hub) Run() {
	for {
		select {
		case participant := <-h.Register:
			if room, ok := h.Rooms[participant.RoomSlug]; ok {
				if _, ok := room.Participants[participant.Id]; !ok {
					fmt.Println("Register participant")
					room.AddParticipant(participant)
					go participant.HandleInit(room, h)
				}
			}
		case participant := <-h.Unregister:
			if room, ok := h.Rooms[participant.RoomSlug]; ok {
				if _, ok := room.Participants[participant.Id]; ok {
					fmt.Println("delete connection")
					// delete participant
					room.RemoveParticipant(participant)
					participant.HandleBailout(room, h)
				}
				// // delete room if no one left
				// if room.Type == PersonalRoom && len(room.Participants) == 0 { // must not be a public room
				// 	// time.AfterFunc(1*time.Hour, func() {
				// 	time.AfterFunc(10*time.Second, func() {
				// 		fmt.Println("hit")
				// 		if room, doesRoomExist := h.Rooms[participant.RoomSlug]; doesRoomExist {
				// 			room.ParticipantMu.Lock()
				// 			if len(room.Participants) == 0 {
				// 				fmt.Println("hi2")
				// 				delete(h.Rooms, room.RoomSlug)
				// 			}
				// 			room.ParticipantMu.Unlock()
				// 		}
				// 	})
				// }
			}

		case message := <-h.Broadcast:
			if room, ok := h.Rooms[message.Slug]; ok { // find the room
				room.BroadcastToPublic(message)
				room.Cache(message)
			}
		}
	}
}

func NewHub() *Hub {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return &Hub{
		Broadcast:               make(chan *MessageResponse, 16),
		Register:                make(chan *Participant, 16),
		Unregister:              make(chan *Participant, 16),
		Rooms:                   make(map[string]*Room),
		Node:                    node,
		SubscriberMessageBuffer: 16,
		PublishLimiter:          rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
	}
}
