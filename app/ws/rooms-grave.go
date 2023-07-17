package ws

import (
	"fmt"
	"sync"
	"time"

	"github.com/pocketbase/pocketbase"
)

type InactiveRoom struct {
	room      *Room
	timestamp time.Time
}

type InactiveRooms []InactiveRoom

type RoomHeap struct {
	// rooms
	data InactiveRooms
	// sync for dead rooms
	ReaperMu *sync.Mutex
}

func (rr *RoomHeap) Less(i, j int) bool { return rr.data[i].timestamp.After(rr.data[j].timestamp) }
func (rr *RoomHeap) Swap(i, j int)      { rr.data[i], rr.data[j] = rr.data[j], rr.data[i] }
func (rr *RoomHeap) Push(x any) {
	rr.data = append(rr.data, x.(InactiveRoom))
}
func (rr *RoomHeap) Pop() any {
	old := rr.data
	n := len(old)
	x := old[n-1]
	rr.data = old[0 : n-1]
	return x
}
func (rr *RoomHeap) Len() int {
	rr.ReaperMu.Lock()
	defer rr.ReaperMu.Unlock()
	return len(rr.data)
}
func (rr *RoomHeap) Dequeue() InactiveRoom {
	rr.ReaperMu.Lock()
	defer rr.ReaperMu.Unlock()
	return rr.Pop().(InactiveRoom)
}
func (rr *RoomHeap) Enqueue(x InactiveRoom) {
	rr.ReaperMu.Lock()
	defer rr.ReaperMu.Unlock()
	rr.Push(x)
}

func CleanupDeadRoom(h *Hub, pb *pocketbase.PocketBase) {
	const deleteWaitPeriod = 1000 * time.Millisecond
	ticker := time.NewTicker(deleteWaitPeriod)
	go func() {
		fmt.Println("Cleanup Running")
		for range ticker.C {
			isPoolEmpty := h.RoomReaper.Len() == 0
			if isPoolEmpty {
				continue
			}
			if room := h.RoomReaper.Dequeue().room; room.Inactive && len(room.Participants) == 0 {
				err := room.DeleteRoomFromDB(h, pb)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	// ticker.Stop()
}
