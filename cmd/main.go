package main

import (
	"log"

	"github.com/sajjad-mobe/gdm/internal/db"
	"github.com/sajjad-mobe/gdm/internal/tui"
)

func main() {
	db.Initialize()

	database := db.GetDB()
	defer database.Close()

	app := tui.NewApp()
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start the app: %v", err)
	}
}
