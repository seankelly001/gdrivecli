package main

import (
	"gdrivecli/pkg/app"
	"gdrivecli/pkg/gdfs"
	"log"
)

func main() {

	gdfs, err := gdfs.NewGDFS()
	if err != nil {
		log.Fatalf("error creating gfds: %s", err.Error())
	}

	app, err := app.NewApp(gdfs)
	if err != nil {
		log.Fatalf("error creating app: %s", err.Error())
	}

	app.Run()
}
