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

	gormDB, err := gorm.Open(db, &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to open database:", err)
		return
	}
	SetCustomDB(gormDB)
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

func Create(model interface{}) error {
	return GetDB().Create(model).Error
}

func GetAll(model interface{}) error {
	return GetDB().Find(model).Error
}
func GetQueueBy(field string, value interface{}) ([]Queue, error) {
	var queues []Queue
	if err := GetDB().Where(field+"= ?", value).Find(&queues).Error; err != nil {
		return nil, err
	}
	return queues, nil
}

func GetDownloadBy(field string, value interface{}) ([]Download, error) {
	var downloads []Download
	if err := GetDB().Preload("Queue").Where(field+"= ?", value).Find(&downloads).Error; err != nil {
		return nil, err
	}
	return downloads, nil
}

func Save(model interface{}) error {
	return GetDB().Save(model).Error
}
