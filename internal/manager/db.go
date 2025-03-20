package manager

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var mu sync.Mutex // Mutex for thread-safe operations

func LoadData() *DataStore {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Open(getDBPath())
	if err != nil {
		return &DataStore{
			Queues:    make(map[string]*Queue),
			Downloads: make(map[string]*Download),
		}
	}
	defer file.Close()

	// Decode JSON file into a DataStore
	var data DataStore
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return &DataStore{
			Queues:    make(map[string]*Queue),
			Downloads: make(map[string]*Download),
		}
	}

	// Ensure maps are initialized if they are nil
	if data.Queues == nil {
		data.Queues = make(map[string]*Queue)
	}
	if data.Downloads == nil {
		data.Downloads = make(map[string]*Download)
	}

	return &data
}

func getDBPath() string {
	// Check if file exists
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Failed to get config directory:", err)
		return "database.json"
	}

	appConfigDir := filepath.Join(configDir, "gdm")
	if err := os.MkdirAll(appConfigDir, os.ModePerm); err != nil {
		log.Fatal("Failed to create config directory:", err)
		return "database.json"
	}

	return filepath.Join(appConfigDir, "database.json")
}

// SaveData writes the DataStore back to the JSON file
func (data *DataStore) Save() error {
	mu.Lock()
	defer mu.Unlock()

	// Open file for writing
	file, err := os.OpenFile(getDBPath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the DataStore as JSON and write to the file
	return json.NewEncoder(file).Encode(data)
}

// AddQueue adds a new Queue to the DataStore
func (data *DataStore) AddQueue(queue *Queue) {
	data.Queues[strconv.Itoa(queue.ID)] = queue
	data.Save()
}

// RemoveQueue removes a Queue from the DataStore
func (data *DataStore) RemoveQueue(queue *Queue) {
	delete(data.Queues, strconv.Itoa(queue.ID))
	data.Save()
}

// AddDownload adds a new Download to the DataStore
func (data *DataStore) AddDownload(download *Download) {
	data.Downloads[strconv.Itoa(download.ID)] = download
	data.Save()
}

// RemoveDownload removes a Download from the DataStore
func (data *DataStore) RemoveDownload(download *Download) {
	delete(data.Downloads, strconv.Itoa(download.ID))
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
