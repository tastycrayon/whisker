package app

import (
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

	rooms := []struct {
		name string
		slug string
	}{
		{"Circular Chatters", "circular-chatters"},
		{"Square Earth Society", "square-earth society"},
		{"Triangular Conversations", "triangular-conversations"},
		{"Polygon Planet", "polygon-planet"},
		{"Elliptical Discussions", "elliptical-discussions"},
		{"Rhombus Roundtable", "rhombus-roundtable"},
		{"Polyhedron Talk", "polyhedron-talk"},
		{"Conical Connections", "conical-connections"},
	}
	for _, room := range rooms {
		hub.Rooms[room.slug] = ws.NewRoom(room.slug, room.name, "Admin", ws.PublicRoom)
	}
	go hub.Run()

	pb.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		InitRoutes(pb, e, hub)
		return nil
	})

	pb.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		if e.Record.Collection().Name == "rooms" {
			authRecord, _ := e.HttpContext.Get(apis.ContextAuthRecordKey).(*models.Record)
			if authRecord == nil {
				return apis.NewForbiddenError("Only authenticated records can perform this endpoint", nil)
			}

			slug := e.Record.GetString("slug")
			name := e.Record.GetString("name")

			if _, roomFound := hub.Rooms[slug]; roomFound {
				return apis.NewApiError(http.StatusBadRequest, "room already exists", nil)
			}

			hub.Rooms[slug] = ws.NewRoom(slug, name, authRecord.Id, ws.PersonalRoom)

		}
		return nil
	})

	if err := pb.Start(); err != nil {
		log.Fatal(err)
	}
}
