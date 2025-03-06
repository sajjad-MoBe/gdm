package manager

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func NewManager(maxParts, partSize int) *DownloadManager {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Failed to get config directory:", err)
	}

	appConfigDir := filepath.Join(configDir, "gdm/tempparts")
	if err := os.MkdirAll(appConfigDir, os.ModePerm); err != nil {
		log.Fatal("Failed to create config directory:", err)
	}
	return &DownloadManager{Queues: []*Queue{}, MaxParts: maxParts, PartSize: partSize, TempFolder: appConfigDir}
}

func (dm *DownloadManager) AddQueue(queue *Queue) {
	queue.IsActive = false
	queue.PartDownloaders = make(chan *PartDownloader, queue.MaxConcurrentDownloads)
	dm.Queues = append(dm.Queues, queue)
	tokenInterval := max(1, 1000_000/queue.MaxBandwidth)

	tokenBucket := make(chan struct{}, queue.MaxBandwidth)
	go func() {
		ticker := time.NewTicker(time.Duration(tokenInterval) * time.Microsecond)
		defer ticker.Stop()
		for range ticker.C {
			select {
			case tokenBucket <- struct{}{}:
			default:
				// Do nothing if the bucket is full
			}
		}
	}()

	go func() {
		for {
			// var progress string
			if IsWithinActiveHours(queue.ActiveStartTime, queue.ActiveEndTime) {
				queue.IsActive = true
				for _, download := range queue.Downloads {

					for {
						if download.Status != "initializing" {
							break
						}
						time.Sleep(time.Millisecond * 500)
					}
					// freeDownloaderss := queue.MaxConcurrentDownloads - len(queue.PartDownloaders)
					// fmt.Println(download.URL, download.Status, freeDownloaderss, len(download.PartDownloaders))
					if download.Status != "pending" {
						continue
					}
					// queueMutex.Lock()
					if !queue.StartAtOneWorkerAvailable {
						// fmt.Println("wait for worker")
						for {
							freeDownloaders := queue.MaxConcurrentDownloads - len(queue.PartDownloaders)
							if freeDownloaders >= len(download.PartDownloaders) {
								break
							}
							time.Sleep(time.Millisecond * 500)
						}
					}
					// queueMutex.Unlock()
					var StartWG sync.WaitGroup
					StartWG.Add(len(download.PartDownloaders))
					// dm.startDownload(download, tokenBucket, &StartWG)
					StartWG.Wait()

				}
			} else {
				// fmt.Println("quqeue", queue.ID, "not working")
				queue.IsActive = false
			}
			// fmt.Printf("\n%s", progress)
			time.Sleep(1 * time.Second)
		}
	}()
}

func (dm *DownloadManager) AddDownload(download *Download) {
	download.Temps = &DownloadTemps{false, 0, 0, time.Now(), &sync.Mutex{}}
	download.IsActive = false
	if download.Status != "finished" && download.Status != "failed" && download.Status != "paused" {
		download.Status = "initializing"
	}
	download.Queue.Downloads = append(download.Queue.Downloads, download)
	go func() {

		resp, err := http.Head(download.URL)
		if err != nil {
			// fmt.Println("Error:", err)
			download.Status = "failed"
			return
		}
		if resp.StatusCode != http.StatusOK {
			download.Status = "failed"
			// fmt.Println("Failed to fetch file details:", resp.Status)
			return
		}
		download.Temps.TotalSize = resp.ContentLength
		resp.Body.Close()

		req, _ := http.NewRequest("GET", download.URL, nil)
		client := &http.Client{}
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", 0, 1))
		resp, err = client.Do(req)
		if err != nil {
			download.Status = "failed"
			// fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusPartialContent {
			download.PartDownloaders = append(download.PartDownloaders, &PartDownloader{
				Index:    0,
				Start:    0,
				End:      0,
				TempFile: download.OutputFile + "-part-0.tmp",
			})

			download.Temps.IsPartial = false
		} else {
			download.Temps.IsPartial = true
			// fmt.Println(download.Temps.TotalSize)
			numParts := min(dm.MaxParts, max(1, int(download.Temps.TotalSize/(int64(dm.PartSize)*1024*1024)))) // each partSize mb add to new part
			partSize := download.Temps.TotalSize / int64(numParts)
			for i := 0; i < numParts; i++ {
				start := partSize * int64(i)
				end := start + partSize - 1
				if i == numParts-1 {
					end = download.Temps.TotalSize - 1
				}
				tempFile := fmt.Sprintf(download.OutputFile+"-part-%d.tmp", i)
				download.PartDownloaders = append(
					download.PartDownloaders,
					&PartDownloader{Index: i, Start: start + getFileSize(tempFile), End: end, TempFile: tempFile},
				)
			}
		}
		download.Status = "pending"
	}()
}
