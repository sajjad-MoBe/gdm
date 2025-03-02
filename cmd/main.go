package main

import (
	"fmt"
	"log"

	"github.com/sajjad-mobe/gdm/internal/db"
	"github.com/sajjad-mobe/gdm/internal/tui"
)

func main() {
	db.Initialize()

	database := db.GetDB()
	fmt.Println(database)
	defer db.Close()

	app := tui.NewApp()
	if _, err := app.Run(); err != nil {
		log.Fatalf("Failed to start the app: %v", err)
	}
}
