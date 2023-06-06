package ws

import (
	"fmt"
	"time"
)

type Hub struct {
	Broadcast  chan *Post
	Register   chan *Participant
	Unregister chan *Participant
	Rooms      map[string]*Room
}

func (h *Hub) Run() {
	for {
		select {
		case participant := <-h.Register:
			fmt.Println("Register client")
			if _, isRoomExist := h.Rooms[participant.RoomId]; !isRoomExist {
				h.Rooms[participant.RoomId] = &Room{
					RoomId:  participant.RoomId,
					Clients: make(map[string]*Participant),
				}
			}
			room := h.Rooms[participant.RoomId]
			if _, doesParticipantExist := room.Clients[participant.ClientId]; !doesParticipantExist {
				room.Clients[participant.ClientId] = participant
			}

		case participant := <-h.Unregister:
			if room, doesRoomExist := h.Rooms[participant.RoomId]; doesRoomExist {
				if _, isCLientExist := room.Clients[participant.ClientId]; isCLientExist {
					fmt.Println("delete connection")
					if len(room.Clients) != 0 {
						post := &Post{
							RoomId: participant.RoomId,
							Message: Message{
								Id:      GenerateId(),
								Content: "An user has left the room.",
								Type:    Bailout,
							},
							Sender: MessageSender{
								ClientId: participant.ClientId,
								Username: participant.Username,
								RoomId:   participant.RoomId,
								Avatar:   participant.Avatar,
							},
							Created: time.Now(),
						}
						h.Broadcast <- post
						// h.Broadcast <- &Message{
						// 	Message:  "disconnect_user",
						// 	ClientId: client.ClientId,
						// 	RoomId:   client.RoomId,
						// 	Username: client.Username,
						// }
					}
					delete(h.Rooms[participant.RoomId].Clients, participant.ClientId)
					close(participant.Message) // stop receiving messages
				}

				// remove room if no one clinet
				clients := h.Rooms[participant.RoomId].Clients
				if len(clients) == 0 {
					delete(h.Rooms, participant.RoomId)
				}
			}

		case message := <-h.Broadcast:
			if room, exist := h.Rooms[message.RoomId]; exist { // find the room
				for _, client := range room.Clients { // find all the paerticipants
					if client.RoomId == message.RoomId {
						client.Message <- message // send them messages
					}
				}
				if room.PostsQueue.IsFull() {
					room.PostsQueue.Dequeue()
				}
				room.PostsQueue.Enqueue(message)
			}
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan *Post, 5),
		Register:   make(chan *Participant),
		Unregister: make(chan *Participant),
		Rooms:      make(map[string]*Room),
	}
}
