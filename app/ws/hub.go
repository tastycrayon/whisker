package ws

import (
	"container/heap"
	"fmt"
	"sync"

	"github.com/bwmarrin/snowflake"
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

	// Rooms added here will be scheduled for death
	RoomReaper *RoomHeap
}

func (h *Hub) Run() {
	for {
		select {
		case participant := <-h.Register:
			if room, ok := h.Rooms[participant.RoomSlug]; ok {
				if _, ok := room.Participants[participant.Id]; !ok {
					fmt.Println("Register participant", participant.Username)
					room.AddParticipant(participant)
					go participant.HandleInit(room, h)
					room.CheckHealth(h)
				}
			}
		case participant := <-h.Unregister:
			if room, ok := h.Rooms[participant.RoomSlug]; ok {
				if _, ok := room.Participants[participant.Id]; ok {
					fmt.Println("delete connection", participant.Username)
					// delete participant
					room.RemoveParticipant(participant)
					participant.HandleBailout(room, h)
					room.CheckHealth(h)
				}
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

	// var pq InactiveRooms = make(InactiveRooms, 0, 128)
	var pq RoomHeap = RoomHeap{
		data:     make(InactiveRooms, 0, 128),
		ReaperMu: &sync.Mutex{},
	}
	heap.Init(&pq)
	return &Hub{
		Broadcast:               make(chan *MessageResponse, 16),
		Register:                make(chan *Participant, 16),
		Unregister:              make(chan *Participant, 16),
		Rooms:                   make(map[string]*Room),
		Node:                    node,
		SubscriberMessageBuffer: 16,
		RoomReaper:              &pq,
	}
}
