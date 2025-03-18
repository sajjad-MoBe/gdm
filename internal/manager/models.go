package manager

import (
	"strconv"
	"sync"
	"time"
)

type Queue struct {
	PartDownloaders           chan *PartDownloader `gorm:"-"`
	Downloads                 []*Download          `gorm:"-"`
	tokenBucket               chan struct{}        `gorm:"-"`
	ticker                    *time.Ticker         `gorm:"-"`
	IsDeleted                 bool                 `gorm:"-"`
	ID                        int                  `gorm:"primaryKey" json:"id"`
	IsActive                  bool                 `gorm:"default:false" json:"is_active"`
	SaveDir                   string               `json:"save_dir"`
	MaxConcurrentDownloads    int                  `gorm:"default:5" json:"max_concurrent_downloads"`
	StartAtOneWorkerAvailable bool                 `gorm:"default:false" json:"start_at_one_worker_available"`
	MaxBandwidth              int                  `gorm:"default:0" json:"max_bandwidth"`
	ActiveStartTime           string               `gorm:"default:'00:00'" json:"active_start_time"`
	ActiveEndTime             string               `gorm:"default:'23:59'" json:"active_end_time"`
	MaxRetries                int                  `gorm:"default:3" json:"max_retries"`
}

func (q Queue) FilterValue() string {
	return strconv.Itoa(q.ID)
}

type Download struct {
	Temps           *DownloadTemps    `gorm:"-"`
	PartDownloaders []*PartDownloader `gorm:"-"`
	IsDeleted       bool              `gorm:"-"`
	ID              int               `gorm:"primaryKey" json:"id"`
	QueueID         int               `gorm:"not null" json:"queue_id"`
	IsActive        bool              `gorm:"default:false" json:"is_active"`
	Status          string            `json:"status"`
	TotalSize       int64             `gorm:"default:0" json:"total_size"`
	IsPartial       bool              `gorm:"default:false" json:"is_partial"`
	OutputFile      string            `json:"output_file"`
	URL             string            `json:"url"`
	Queue           *Queue            `gorm:"foreignKey:QueueID"`
}

type DownloadTemps struct {
	TotalDownloaded int64
	Retries         int
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

func (d *Download) GetStatus() string {
	return d.Status
}
func (d *Download) GetSpeed() int {
	totalKB := 0
	for _, p := range d.PartDownloaders {
		totalKB += int(p.Speed / 1024)
	}
	return totalKB
}
func (d *Download) GetProgress() int {
	if d.TotalSize == 0 || d.Temps == nil {
		return 0
	}
	return int(d.Temps.TotalDownloaded * 100 / d.TotalSize)
}
