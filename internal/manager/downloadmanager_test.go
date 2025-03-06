package manager

import (
	"os"
	"path"
	"testing"
	"time"
)

func TestDownloadManager(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("can't get download folder, %s", err)
		return
	}
	queue1 := Queue{
		SaveDir:                   path.Join(homeDir, "Downloads/gdm-test"), // download folder
		MaxConcurrentDownloads:    5,
		StartAtOneWorkerAvailable: true,
		MaxBandwidth:              -1, // for no limit
		ActiveStartTime:           "00:00",
		ActiveEndTime:             "23:59",
		MaxRetries:                1,
	}
	queue2 := Queue{
		SaveDir:                   path.Join(homeDir, "Downloads/gdm-test"),
		MaxConcurrentDownloads:    10,
		StartAtOneWorkerAvailable: false,
		MaxBandwidth:              500, // for 500 Kb/s limit
		ActiveStartTime:           "00:00",
		ActiveEndTime:             "23:59",
		MaxRetries:                4,
	}
	download1 := Download{
		QueueID:    queue1.ID,
		Queue:      &queue1,
		Status:     "pending",
		OutputFile: "example.html",
		URL:        "https://example.com",
		Retries:    0,
	}
	download2 := Download{
		QueueID:    queue2.ID,
		Queue:      &queue2,
		Status:     "pending",
		OutputFile: "google.html",
		URL:        "https://google.com",
		Retries:    0,
	}
	download3 := Download{
		QueueID:    queue1.ID,
		Queue:      &queue2,
		Status:     "pending",
		OutputFile: "10mb.zip",
		URL:        "http://212.183.159.230/10MB.zip",
		Retries:    0,
	}

	MaxParts := 10 // Maximum number of parts for one download
	PartSize := 3  // create new part downloader per each PartSize mb
	downloadManager := NewManager(MaxParts, PartSize)

	downloadManager.AddQueue(&queue1)
	downloadManager.AddQueue(&queue2)

	// downloadManager.AddDownload(&download1)
	downloadManager.AddDownload(&download3)
	// downloadManager.AddDownload(&download2)
	_, _ = download1, download2

	time.Sleep(time.Second * 2)
	downloadManager.PauseDownload(&download1)
	time.Sleep(time.Second * 10)
	downloadManager.ResumeDownload(&download1)
	queue1.SetBandwith(10)
	time.Sleep(time.Second * 10)

	for {
		ended := true
		for _, queue := range downloadManager.Queues {
			for _, download := range queue.Downloads {
				if download.Status == "initializing" ||
					download.Status == "pending" ||
					download.Status == "downloading" {
					ended = false
					time.Sleep(5 * time.Second)
					break
				}
			}
		}
		if ended {
			break
		}
	}
}
