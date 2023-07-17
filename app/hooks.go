package app

import (
	"fmt"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/tastycrayon/pb-svelte-chatapp/app/ws"
)

func InitHooks(pb *pocketbase.PocketBase, h *ws.Hub) {
	// load default rooms
	pb.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		records, err := pb.Dao().FindRecordsByExpr("rooms")
		if err != nil {
			fmt.Println("err", err)
		}
		for _, r := range records {
			id := r.GetId()
			slug := r.GetString("slug")
			name := r.GetString("name")
			description := r.GetString("description")
			cover := r.GetString("cover")

			var roomType ws.RoomType = ws.PersonalRoom
			if rawRoomType := r.GetString("type"); rawRoomType == string(ws.PublicRoom) {
				roomType = ws.PublicRoom
			}
			createdBy := r.GetString("createdBy")
			created := r.GetCreated().String()
			h.Rooms[slug] = ws.NewRoom(id, slug, name, cover, description, createdBy, created, roomType)
		}
		return nil
	})

	// after room creation, load into memory
	pb.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		if e.Record.Collection().Name == "rooms" {
			id := e.Collection.GetId()
			slug := e.Record.GetString("slug")
			name := e.Record.GetString("name")
			cover := e.Record.GetString("cover")
			description := e.Record.GetString("description")
			createdBy := e.Record.GetString("createdBy")
			created := e.Record.GetCreated().String()

			if _, roomFound := h.Rooms[slug]; roomFound {
				return apis.NewApiError(http.StatusBadRequest, "room already exists", nil)
			}
			h.Rooms[slug] = ws.NewRoom(id, slug, name, cover, description, createdBy, created, ws.PersonalRoom)
		}
		return nil
	})

	// after update, mutate local copy
	pb.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		if e.Record.Collection().Name == "rooms" {
			oldSlug := e.Record.OriginalCopy().GetString("slug")
			room, roomFound := h.Rooms[oldSlug]
			if !roomFound {
				return apis.NewApiError(http.StatusBadRequest, "room was not found", nil)
			}
			room.ParticipantMu.Lock()
			defer room.ParticipantMu.Unlock()
			room.RoomName = e.Record.GetString("name")
			// room.RoomSlug = e.Record.GetString("slug") // dont update slug
			room.RoomCover = e.Record.GetString("cover")
			room.RoomDescription = e.Record.GetString("description")
			// room.CreatedBy = e.Record.GetString("createdBy")
			room.Created = e.Record.GetCreated().String()
		}
		return nil
	})

	// after room deletion, remove from local memory
	pb.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		if e.Record.Collection().Name == "rooms" {
			slug := e.Record.GetString("slug")
			if room, ok := h.Rooms[slug]; ok {
				room.DeleteRoomFromMemory(h, pb)
				return nil
			}
			return apis.NewApiError(http.StatusBadRequest, "room was not found", nil)
		}
		return nil
	})

}
