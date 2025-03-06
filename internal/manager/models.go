package manager

import (
	"sync"
	"time"
)

type Queue struct {
	PartDownloaders           chan *PartDownloader `gorm:"-"`
	Downloads                 []*Download          `gorm:"-"`
	tokenBucket               chan struct{}        `gorm:"-"`
	ticker                    *time.Ticker         `gorm:"-"`
	ID                        uint                 `gorm:"primaryKey" json:"id"`
	IsActive                  bool                 `gorm:"default:false" json:"is_active"`
	SaveDir                   string               `json:"save_dir"`
	MaxConcurrentDownloads    int                  `gorm:"default:5" json:"max_concurrent_downloads"`
	StartAtOneWorkerAvailable bool                 `gorm:"default:false" json:"start_at_one_worker_available"`
	MaxBandwidth              int                  `gorm:"default:-1" json:"max_bandwidth"`
	ActiveStartTime           string               `gorm:"default:'00:00'" json:"active_start_time"`
	ActiveEndTime             string               `gorm:"default:'23:59'" json:"active_end_time"`
	MaxRetries                int                  `gorm:"default:3" json:"max_retries"`
}

type Download struct {
	Temps           *DownloadTemps    `gorm:"-"`
	PartDownloaders []*PartDownloader `gorm:"-"`
	ID              uint              `gorm:"primaryKey" json:"id"`
	QueueID         uint              `gorm:"not null" json:"queue_id"`
	IsActive        bool              `gorm:"default:false" json:"is_active"`
	Status          string            `json:"status"`
	TotalSize       int64             `gorm:"default:0" json:"total_size"`
	IsPartial       bool              `gorm:"default:false" json:"is_partial"`
	OutputFile      string            `json:"output_file"`
	URL             string            `json:"url"`
	Retries         int               `gorm:"default:0" json:"retries"`
	Queue           *Queue            `gorm:"foreignKey:QueueID"`
}

type DownloadTemps struct {
	TotalDownloaded int64
	StartTime       time.Time
	Mutex           *sync.Mutex
}

type PartDownloader struct {
	Index      int
	Start      int64
	End        int64
	Downloaded int64
	Speed      int64
	TempFile   string
	IsFailed   bool
	IsPaused   bool
}

type DownloadManager struct {
	Queues     []*Queue
	MaxParts   int
	PartSize   int
	TempFolder string
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
