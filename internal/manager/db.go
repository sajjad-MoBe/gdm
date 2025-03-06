package manager

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var database *gorm.DB

func InitializeDB() {
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
	db := sqlite.Open(dbPath)

	database, err := gorm.Open(db, &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to open database:", err)
		return
	}
	err = database.AutoMigrate(&Queue{}, &Download{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

}

func CloseDB() {
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.Close()
}

func GetDB() *gorm.DB {
	return database
}
func SetCustomDB(db *gorm.DB) {
	database = db
}

func ResetAll(tempDir string) error {
	return filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	})
}
