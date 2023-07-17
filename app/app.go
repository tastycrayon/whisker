package app

import (
	"io/fs"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/tastycrayon/pb-svelte-chatapp/app/ws"
)

func Run(buildDirFs, indexFileFs fs.FS) {
	pb := pocketbase.New()
	// migration
	migratecmd.MustRegister(pb, pb.RootCmd, &migratecmd.Options{
		Automigrate: true, // auto creates migration files when making collection changes
	})

	hub := ws.NewHub()
	ws.CleanupDeadRoom(hub, pb)

	// hooks
	InitHooks(pb, hub)
	pb.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		InitRoutes(pb, e, hub, buildDirFs, indexFileFs)
		return nil
	})

	go hub.Run()

	if err := pb.Start(); err != nil {
		log.Fatal(err)
	}
}
