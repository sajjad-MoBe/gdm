package db

import (
	"fmt"
	"log"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateAndGetObject(t *testing.T) {
	// Initialize the database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&Queue{}, &Download{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}
	SetCustomDB(db)

	// Create a new queue entry
	queue := Queue{
		SaveDir:                "/downloads",
		MaxConcurrentDownloads: 5,
		MaxBandwidth:           -1,
		ActiveStartTime:        "00:00",
		ActiveEndTime:          "23:59",
		MaxRetries:             1,
	}

	// Test creating the queue
	if err := Create(&queue); err != nil {
		t.Fatalf("failed to create queue: %v", err)
	}

	// Test retrieving queues
	var queues []Queue

	if err := GetAll(&queues); err != nil {
		t.Fatalf("failed to get queues: %v", err)
	}

	// Check if the queue was created
	if len(queues) != 1 {
		t.Fatalf("expected 1 queue, got %d", len(queues))
	}

	// Validate the created queue
	if queues[0].SaveDir != queue.SaveDir {
		t.Errorf("expected SaveDir %s, got %s", queue.SaveDir, queues[0].SaveDir)
	}

	queues, err = GetQueueBy("save_dir", "/downloads")
	if err != nil {
		log.Fatalf("Error retrieving queues with %s='%s': %v", "save_dir", "/downloads", err)
	}
	for _, p := range queues {
		fmt.Printf("Queue: %+v\n", p)
	}

	download := Download{
		QueueID: queue.ID, // Use the created Queue ID
		// Queue:   queue,
		Status:  "pending",
		URL:     "https://example.com/download",
		Retries: 0,
	}

	// Save the Download to the database
	if err := Create(&download); err != nil {
		t.Fatalf("Failed to create download: %v", err)
	}

	// Verify the Download was saved correctly
	downloads, err := GetDownloadBy("id", download.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve download: %v", err)
	}
	for _, p := range downloads {
		fmt.Printf("Download: %+v\n", p)
	}
	savedDownload := downloads[0]

	// Assertions
	if savedDownload.QueueID != download.QueueID {
		t.Errorf("Expected QueueID %v, got %v", download.QueueID, savedDownload.QueueID)
	}
	if savedDownload.Status != download.Status {
		t.Errorf("Expected Status %s, got %s", download.Status, savedDownload.Status)
	}
	if savedDownload.URL != download.URL {
		t.Errorf("Expected URL %s, got %s", download.URL, savedDownload.URL)
	}
	if savedDownload.Retries != download.Retries {
		t.Errorf("Expected Retries %d, got %d", download.Retries, savedDownload.Retries)
	}

	// test for save object
	savedDownload.Retries = 1
	if err := Save(&savedDownload); err != nil {
		log.Fatalf("Failed to save updated download: %v", err)
	}
	// Verify the Download was changed correctly
	downloads, err = GetDownloadBy("id", download.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve download: %v", err)
	}
	if downloads[0].Retries != savedDownload.Retries {
		t.Errorf("Expected Retries %d, got %d", downloads[0].Retries, savedDownload.Retries)
	}
}
