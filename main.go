package main

import (
	"embed"

	"github.com/labstack/echo/v5"
	"github.com/tastycrayon/whisker/app"
)

var (
	//go:embed all:client/build
	buildDir embed.FS

	//go:embed client/build/index.html
	indexDir embed.FS

	buildDirFs  = echo.MustSubFS(buildDir, "client/build")
	indexFileFs = echo.MustSubFS(indexDir, "client/build")
)

func main() {
	app.Run(buildDirFs, indexFileFs)
}
