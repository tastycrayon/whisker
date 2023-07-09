package main

import (
	"embed"

	"github.com/tastycrayon/go-chat/app"
)

//go:embed build
var embeddedFiles embed.FS

func main() {
	app.Run(embeddedFiles)
}
