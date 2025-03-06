package main

import (
	"fmt"
	"log"

	"github.com/sajjad-mobe/gdm/internal/manager"
	"github.com/sajjad-mobe/gdm/internal/tui"
)

func main() {
	manager.InitializeDB()

	database := manager.GetDB()
	fmt.Println(database)
	defer manager.CloseDB()

	app := tui.NewApp()
	if _, err := app.Run(); err != nil {
		log.Fatalf("Failed to start the app: %v", err)
	}
}
