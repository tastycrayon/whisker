package app

import (
	"fmt"
	"log"
	"net/http"

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
	DefaultRoomList := GetDefaultRooms()
	for _, room := range DefaultRoomList {
		hub.Rooms[room.slug] = ws.NewRoom(
			room.id, room.slug, room.name, room.cover, room.description, "Admin",
			"6/24/2023, 8:51:08 PM",
			ws.PublicRoom)
	}
	// TODO: cover
	go hub.Run()
	// load default rooms end

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
			createdBy := r.GetString("createdBy")
			created := r.GetCreated().String()
			hub.Rooms[slug] = ws.NewRoom(id, slug, name, cover, description, createdBy, created, ws.PersonalRoom)
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
			createdBy := authRecord.GetString("createdBy")
			created := authRecord.GetCreated().String()

			if _, roomFound := hub.Rooms[slug]; roomFound {
				return apis.NewApiError(http.StatusBadRequest, "room already exists", nil)
			}

			hub.Rooms[slug] = ws.NewRoom(id, slug, name, cover, description, createdBy, created, ws.PersonalRoom)
		}
		return nil
	})
	// pb.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
	// 	if e.Record.Collection().Name == "rooms" {
	// 		authRecord, _ := e.HttpContext.Get(apis.ContextAuthRecordKey).(*models.Record)
	// 		if authRecord == nil {
	// 			return apis.NewForbiddenError("Only authenticated records can perform this endpoint", nil)
	// 		}

	// 		// slug := e.Record.GetString("slug")

	// 		// if _, roomFound := hub.Rooms[slug]; !roomFound {
	// 		// 	return apis.NewApiError(http.StatusBadRequest, "room was not found", nil)
	// 		// }
	// 		// delete(hub.Rooms, slug)
	// 		fmt.Println("hit", e.Record.GetString("slug"))
	// 	}
	// 	return nil
	// })

	if err := pb.Start(); err != nil {
		log.Fatal(err)
	}
}
