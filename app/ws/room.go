package ws

import (
	"sync"

	"github.com/tastycrayon/pb-svelte-chatapp/app/queue"
	"nhooyr.io/websocket"
)

const messageCacheSize = 16

type RoomType string

const (
	PublicRoom   RoomType = "public"
	PersonalRoom RoomType = "personal"
)

type Room struct {
	RoomId string `json:"id"`
	// URL for Room, Eg: /rooms/my-awesome-room
	RoomSlug string `json:"slug"`
	// Room name, Eg: My awesome room
	RoomName string `json:"name"`
	// cover image of the room
	RoomCover string `json:"cover"`
	// A short description about the room
	RoomDescription string `json:"description"`
	// Public room and personal room
	// Public - created by default
	// Personal - created by participants
	Type RoomType `json:"type"`
	// Map of participants
	Participants ParticipantList `json:"-"`
	// The id of the admin of the room
	CreatedBy string `json:"createdBy"`
	// Time the room was made
	Created string `json:"created"`
	// A circular queue to persist messages, Default 128 messages
	LocalMessageQueue *queue.CQueue[MessageResponse] `json:"-"`

	// Syncing
	ParticipantMu *sync.Mutex `json:"-"`
}

func NewRoom(roomId, roomSlug, roomName, roomCover, description, createdBy string, created string, roomType RoomType) *Room {
	return &Room{
		RoomId:            roomId,
		RoomSlug:          roomSlug,
		RoomName:          roomName,
		RoomCover:         roomCover,
		RoomDescription:   description,
		CreatedBy:         createdBy,
		Type:              roomType,
		Created:           created,
		Participants:      make(map[string]*Participant),
		LocalMessageQueue: queue.NewCQueue[MessageResponse](messageCacheSize),
		ParticipantMu:     &sync.Mutex{},
	}
}

func (r *Room) DeleteRoom(h *Hub) {
	r.ParticipantMu.Lock()
	defer r.ParticipantMu.Unlock()
	for _, p := range r.Participants {
		//TODO: sent message informing the deletion
		h.Unregister <- p
	}
	delete(h.Rooms, r.RoomSlug)
}

// adds a new participant to the room
func (r *Room) AddParticipant(p *Participant) {
	r.ParticipantMu.Lock()
	defer r.ParticipantMu.Unlock()

	p.RoomSlug = r.RoomSlug  // add room to participant
	r.Participants[p.Id] = p // add participant to room
}

// removes a new participant from the room
func (r *Room) RemoveParticipant(p *Participant) {
	r.ParticipantMu.Lock()
	defer r.ParticipantMu.Unlock()

	if _, ok := r.Participants[p.Id]; ok {
		delete(r.Participants, p.Id)
		p.Conn.Close(websocket.StatusNormalClosure, "deleting participant "+p.Id)
	}
}

func (r *Room) BroadcastToPublic(message *MessageResponse) {
	r.ParticipantMu.Lock()
	defer r.ParticipantMu.Unlock()
	for _, participant := range r.Participants { // find all the paerticipants
		participant.Message <- message // send them messages
	}
}

func (r *Room) Cache(message *MessageResponse) {
	r.ParticipantMu.Lock()
	defer r.ParticipantMu.Unlock()

	if r.LocalMessageQueue.IsFull() {
		r.LocalMessageQueue.Dequeue()
	}
	r.LocalMessageQueue.Enqueue(message)
}
