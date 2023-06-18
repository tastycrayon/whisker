package ws

import (
	"sync"

	"github.com/tastycrayon/go-chat/app/queue"
)

type RoomType string

const (
	PublicRoom   RoomType = "public"
	PersonalRoom RoomType = "personal"
)

type Room struct {
	// URL for Room, Eg: /rooms/my-awesome-room
	RoomSlug string `json:"roomSlug"`
	// Room name, Eg: My awesome room
	RoomName string `json:"roomName"`
	// Public room and personal room
	// Public - created by default
	// Personal - created by participants
	Type RoomType `json:"roomType"`
	// Map of participants
	Participants map[string]*Participant `json:"-"`
	// The admin of the room
	CreatedBy string `json:"createdBy"`
	// A circular queue to persist messages, Default 128 messages
	LocalMessageQueue *queue.CQueue[MessageFeed] `json:"-"`

	// Syncing
	ParticipantMu sync.Mutex `json:"-"`
}

func NewRoom(roomSlug, roomName, createdBy string, roomType RoomType) *Room {
	return &Room{
		RoomSlug:          roomSlug,
		RoomName:          roomName,
		CreatedBy:         createdBy,
		Type:              roomType,
		Participants:      make(map[string]*Participant),
		LocalMessageQueue: queue.NewCQueue[MessageFeed](128),
	}
}
