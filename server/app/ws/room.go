package ws

import (
	"sync"

	"github.com/tastycrayon/go-chat/app/queue"
	"nhooyr.io/websocket"
)

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
	Participants map[string]*Participant `json:"-"`
	// The id of the admin of the room
	CreatedBy string `json:"createdBy"`
	// Time the room was made
	Created string `json:"created"`
	// A circular queue to persist messages, Default 128 messages
	LocalMessageQueue *queue.CQueue[MessageFeed] `json:"-"`

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
		LocalMessageQueue: queue.NewCQueue[MessageFeed](128),
		ParticipantMu:     &sync.Mutex{},
	}
}

func (r *Room) Delete(h *Hub) {
	r.ParticipantMu.Lock()
	for _, p := range r.Participants {
		// h.Unregister <- p
		go p.Conn.Close(websocket.StatusNormalClosure, "room was permanently deleted")
	}
	r.ParticipantMu.Unlock()
	delete(h.Rooms, r.RoomSlug)
}
