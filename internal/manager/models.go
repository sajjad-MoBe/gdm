package manager

import (
	"strconv"
	"sync"
	"time"
)

// Queue represents a single queue item
type Queue struct {
	PartDownloaders           chan *PartDownloader `json:"-"`
	Downloads                 []*Download          `json:"-"`
	tokenBucket               chan struct{}        `json:"-"`
	ticker                    *time.Ticker         `json:"-"`
	IsRemoved                 bool                 `json:"-"`
	ID                        int                  `json:"id"`
	IsActive                  bool                 `json:"is_active"`
	SaveDir                   string               `json:"save_dir"`
	MaxConcurrentDownloads    int                  `json:"max_concurrent_downloads"` // default 10
	StartAtOneWorkerAvailable bool                 `json:"start_at_one_worker_available"`
	MaxBandwidth              int                  `json:"max_bandwidth"`     // default 0 for unlimited
	ActiveStartTime           string               `json:"active_start_time"` // default 00:00
	ActiveEndTime             string               `json:"active_end_time"`   // default 23:59
	MaxRetries                int                  `json:"max_retries"`       // default 3
}

func (q Queue) FilterValue() string {
	return strconv.Itoa(q.ID)
}

type Download struct {
	Temps           *DownloadTemps    `json:"-"`
	PartDownloaders []*PartDownloader `json:"-"`
	IsRemoved       bool              `json:"-"`
	Queue           *Queue            `json:"_"`
	ID              int               `json:"id"`
	QueueID         int               `json:"queue_id"`
	IsActive        bool              `json:"is_active"`
	Status          string            `json:"status"`
	TotalSize       int64             `json:"total_size"`
	IsPartial       bool              `json:"is_partial"`
	OutputFile      string            `json:"output_file"`
	URL             string            `json:"url"`
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

// DataStore holds the queues and downloads
type DataStore struct {
	Queues    map[string]*Queue    `json:"queues"`    // Map with ID as key and Queue as value
	Downloads map[string]*Download `json:"downloads"` // Map with ID as key and generic download data
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
