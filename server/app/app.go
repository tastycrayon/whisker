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
