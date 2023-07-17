package ws

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/tastycrayon/pb-svelte-chatapp/app/queue"
	"golang.org/x/time/rate"
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

	// publishLimiter controls the rate limit applied to the publish endpoint.
	//
	// Defaults to one publish every 100ms with a burst of 8.
	PublishLimiter *rate.Limiter `json:"-"`

	// inactive means it is scheduled for death
	Inactive bool
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
		PublishLimiter:    rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
		Inactive:          false,
	}
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
	// for _, participant := range r.Participants { // find all the paerticipants
	// 	participant.Message <- message // send them messages
	// }
	r.PublishLimiter.Wait(context.Background())

	for _, p := range r.Participants { // find all the paerticipants
		select {
		case p.Message <- message: // send them messages
		default:
			// It never blocks and so messages to slow subscribers are dropped.
			go p.CloseSlow(p.Conn)
		}
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

func (r *Room) CheckHealth(h *Hub) {
	r.ParticipantMu.Lock()
	defer r.ParticipantMu.Unlock()

	// must be a personal room
	if r.Type != PersonalRoom {
		return
	}
	// anytime a new user joins, room will be active
	if len(r.Participants) != 0 {
		r.Inactive = false
		return
	}

	r.Inactive = true // mark as inactive
	h.RoomReaper.Enqueue(InactiveRoom{room: r, timestamp: time.Now()})

}

func (r *Room) DeleteRoomFromDB(h *Hub, pb *pocketbase.PocketBase) error {
	// if record found, proceed with delete
	fmt.Printf("deleted room from db: %v | ", r.RoomName)
	return pb.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		record, err := txDao.FindRecordById("rooms", r.RoomId)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if err := txDao.DeleteRecord(record); err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
}

func (r *Room) DeleteRoomFromMemory(h *Hub, pb *pocketbase.PocketBase) {
	r.ParticipantMu.Lock()
	defer r.ParticipantMu.Unlock()
	// disconnects everyone
	for _, p := range r.Participants {
		//TODO: sent message informing the deletion
		h.Unregister <- p
	}
	// removes from memory
	delete(h.Rooms, r.RoomSlug)
	fmt.Println("deleted room from memory: ", r.RoomName)
}
