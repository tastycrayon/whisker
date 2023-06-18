package main

import (
	"github.com/tastycrayon/go-chat/app"
)

const MaximumMessageQueueSize = 128

func main() {
	app.Run()

	// router
	// room -> join, create, delete

	// app.OnRecordBeforeUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
	// 	if e.Record.Collection().Name == "users" {
	// 		switch e.Record.Get("role").(int) {
	// 		case READER:
	// 			e.Record.Set("role", READER)
	// 		case AUTHOR:
	// 			e.Record.Set("role", AUTHOR)
	// 		}
	// 	}
	// 	return nil
	// })
	// server(app) // init server

}
