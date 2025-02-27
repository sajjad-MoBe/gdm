package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func Initialize() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Failed to get config directory:", err)
	}

	appConfigDir := filepath.Join(configDir, "gdm")
	if err := os.MkdirAll(appConfigDir, os.ModePerm); err != nil {
		log.Fatal("Failed to create config directory:", err)
	}

	dbPath := filepath.Join(appConfigDir, "database.db")
	fmt.Println("Database path:", dbPath)

	database, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	createTables()

	// handleShutdown()
}

func createTables() {
	query := `
	
	`
	_, err := database.Exec(query)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}
}

func GetDB() *sql.DB {
	return database
}

func Close() {
	if database != nil {
		err := database.Close()
		if err != nil {
			log.Println("Error closing database:", err)
		} else {
			fmt.Println("Database connection closed.")
		}
	}
}

// func handleShutdown() {
// 	sigChan := make(chan os.Signal, 1)
// 	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

// 	go func() {
// 		<-sigChan
// 		fmt.Println("\nclosing app...")
// 		Close()
// 		os.Exit(0)
// 	}()
// }
