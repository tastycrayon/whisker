package ws

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"golang.org/x/time/rate"
)

type Hub struct {
	Broadcast  chan *MessageFeed
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
			fmt.Println("Register participant")
			if room, doesRoomExist := h.Rooms[participant.RoomSlug]; doesRoomExist {
				if _, doesParticipantExist := room.Participants[participant.Id]; !doesParticipantExist {
					room.ParticipantMu.Lock()
					room.Participants[participant.Id] = participant
					room.ParticipantMu.Unlock()
					go func(r *Room) {
						r.ParticipantMu.Lock()
						defer r.ParticipantMu.Unlock()
						// fetch exisiting messages and post
						mList := GenerateCustomMessage(
							participant,
							GenerateId(h),
							History,
							room.LocalMessageQueue.GetInternalSlice(),
						)
						participant.WriteCustomMessage(mList)
						// fetching participant list
						pList := GenerateCustomMessage(
							participant,
							GenerateId(h),
							Ping,
							GetAllParticipants(r),
						)
						participant.WriteCustomMessage(pList)
					}(room)
					h.Broadcast <- GenerateUserJoinedMessage(participant, GenerateId(h))
				}
			}
		case participant := <-h.Unregister:
			if room, doesRoomExist := h.Rooms[participant.RoomSlug]; doesRoomExist {
				if _, doesParticipant := room.Participants[participant.Id]; doesParticipant {
					fmt.Println("delete connection")
					post := GenerateUserLeftMessage(participant, GenerateId(h))
					h.Broadcast <- post
					// delete participant
					room.ParticipantMu.Lock()
					delete(h.Rooms[participant.RoomSlug].Participants, participant.Id)
					room.ParticipantMu.Unlock()
				}
				// delete room if no one left
				if room.Type != PublicRoom { // must not be a public room
					time.AfterFunc(1*time.Hour, func() {
						if room, doesRoomExist := h.Rooms[participant.RoomSlug]; doesRoomExist {
							room.ParticipantMu.Lock()
							if len(room.Participants) == 0 {
								delete(h.Rooms, participant.RoomSlug)
							}
							room.ParticipantMu.Unlock()
						}
					})
				}
			}

		case message := <-h.Broadcast:
			if room, exist := h.Rooms[message.RoomSlug]; exist { // find the room
				room.ParticipantMu.Lock()
				for _, participant := range room.Participants { // find all the paerticipants
					participant.Message <- message // send them messages
				}
				room.ParticipantMu.Unlock()
				if room.LocalMessageQueue.IsFull() {
					room.LocalMessageQueue.Dequeue()
				}
				room.LocalMessageQueue.Enqueue(message)
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
		Broadcast:               make(chan *MessageFeed, 16),
		Register:                make(chan *Participant, 16),
		Unregister:              make(chan *Participant, 16),
		Rooms:                   make(map[string]*Room),
		Node:                    node,
		SubscriberMessageBuffer: 16,
		PublishLimiter:          rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
	}
}
