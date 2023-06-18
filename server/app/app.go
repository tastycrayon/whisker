package app

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/tastycrayon/go-chat/app/ws"
)

func Run() {
	pb := pocketbase.New()
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
	migratecmd.MustRegister(pb, pb.RootCmd, &migratecmd.Options{
		Automigrate: true, // auto creates migration files when making collection changes
	})

	if err := pb.Start(); err != nil {
		log.Fatal(err)
	}
}
