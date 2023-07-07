package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/tastycrayon/go-chat/app/ws"
)

func Run() {
	pb := pocketbase.New()

	// migration
	migratecmd.MustRegister(pb, pb.RootCmd, &migratecmd.Options{
		Automigrate: true, // auto creates migration files when making collection changes
	})

	hub := ws.NewHub()

	// load default rooms
	// DefaultRoomList := GetDefaultRooms()
	// for _, room := range DefaultRoomList {
	// 	hub.Rooms[room.slug] = ws.NewRoom(
	// 		room.id, room.slug, room.name, room.cover, room.description, "Admin",
	// 		"6/24/2023, 8:51:08 PM",
	// 		ws.PublicRoom)
	// }
	go hub.Run()
	// load default rooms end
	pb.OnRecordsListRequest().Add(func(e *core.RecordsListEvent) error {
		if e.Collection.Name == "rooms" {
			time.Sleep(time.Second * 2)
			fmt.Println("rrom")
		}
		return nil
	})
	pb.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		records, err := pb.Dao().FindRecordsByExpr("rooms")
		if err != nil {
			fmt.Println("err", err)
		}

		for _, r := range records {
			id := r.GetString("id")
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
			hub.Rooms[slug] = ws.NewRoom(id, slug, name, cover, description, createdBy, created, roomType)
		}

		InitRoutes(pb, e, hub)
		return nil
	})

	pb.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		if e.Record.Collection().Name == "rooms" {
			authRecord, _ := e.HttpContext.Get(apis.ContextAuthRecordKey).(*models.Record)
			if authRecord == nil {
				return apis.NewForbiddenError("Only authenticated records can perform this endpoint", nil)
			}
			id := e.Record.GetString("id")
			slug := e.Record.GetString("slug")
			name := e.Record.GetString("name")
			cover := e.Record.GetString("cover")
			description := e.Record.GetString("description")
			createdBy := e.Record.GetString("createdBy")
			created := e.Record.GetCreated().String()

			if _, roomFound := hub.Rooms[slug]; roomFound {
				return apis.NewApiError(http.StatusBadRequest, "room already exists", nil)
			}

			hub.Rooms[slug] = ws.NewRoom(id, slug, name, cover, description, createdBy, created, ws.PersonalRoom)
		}
		return nil
	})

	pb.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		if e.Record.Collection().Name == "rooms" {
			authRecord, _ := e.HttpContext.Get(apis.ContextAuthRecordKey).(*models.Record)
			if authRecord == nil {
				return apis.NewForbiddenError("Only authenticated records can perform this endpoint", nil)
			}

			oldSlug := e.Record.OriginalCopy().GetString("slug")
			room, roomFound := hub.Rooms[oldSlug]
			if !roomFound {
				return apis.NewApiError(http.StatusBadRequest, "room was not found", nil)
			}
			room.ParticipantMu.Lock()
			room.RoomName = e.Record.GetString("name")
			// room.RoomSlug = e.Record.GetString("slug") // dont update slug
			room.RoomCover = e.Record.GetString("cover")
			room.RoomDescription = e.Record.GetString("description")
			// room.CreatedBy = e.Record.GetString("createdBy")
			room.Created = e.Record.GetCreated().String()
			room.ParticipantMu.Unlock()
		}
		return nil
	})

	pb.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		if e.Record.Collection().Name == "rooms" {
			authRecord, _ := e.HttpContext.Get(apis.ContextAuthRecordKey).(*models.Record)
			if authRecord == nil {
				return apis.NewForbiddenError("Only authenticated records can perform this endpoint", nil)
			}
			slug := e.Record.GetString("slug")
			if room, roomFound := hub.Rooms[slug]; roomFound {
				room.Delete(hub)
			}
			return apis.NewApiError(http.StatusBadRequest, "room was not found", nil)
		}
		return nil
	})

	if err := pb.Start(); err != nil {
		log.Fatal(err)
	}
}
