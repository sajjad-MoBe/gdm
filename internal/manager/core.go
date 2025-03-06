package manager

import (
	"fmt"
	"io"
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
		configDir = "./"
	}

	TempFolder := filepath.Join(configDir, "gdm/tempparts")
	if err := os.MkdirAll(TempFolder, os.ModePerm); err != nil {
		log.Fatal("Failed to create temp directory:", err)
	}
	return &DownloadManager{Queues: []*Queue{}, MaxParts: maxParts, PartSize: partSize, TempFolder: TempFolder}
}

func (dm *DownloadManager) AddQueue(queue *Queue) {
	queue.IsActive = false
	queue.PartDownloaders = make(chan *PartDownloader, queue.MaxConcurrentDownloads)
	dm.Queues = append(dm.Queues, queue)
	if queue.MaxBandwidth > 0 {
		queue.SetBandwith(queue.MaxBandwidth)
	}
	go func() {
		for {
			if IsWithinActiveHours(queue.ActiveStartTime, queue.ActiveEndTime) {
				queue.IsActive = true
				for _, download := range queue.Downloads {

					for {
						if download.Status != "initializing" {
							break
						}
						time.Sleep(time.Millisecond * 500)
					}

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
					dm.startDownload(download, &StartWG)
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
	download.Temps = &DownloadTemps{0, time.Now(), &sync.Mutex{}}
	download.IsActive = false
	if download.Status != "finished" && download.Status != "failed" && download.Status != "paused" {
		download.Status = "initializing"
	}
	download.Queue.Downloads = append(download.Queue.Downloads, download)
	go func() {
		if download.TotalSize < 1 {
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
			download.TotalSize = resp.ContentLength
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
				download.IsPartial = false
			} else {
				download.IsPartial = true
				// fmt.Println(download.Temps.TotalSize)
			}
		}
		if download.IsPartial {
			numParts := min(dm.MaxParts, max(1, int(download.TotalSize/(int64(dm.PartSize)*1024*1024)))) // each partSize mb add to new part
			partSize := download.TotalSize / int64(numParts)
			for i := 0; i < numParts; i++ {
				start := partSize * int64(i)
				end := start + partSize - 1
				if i == numParts-1 {
					end = download.TotalSize - 1
				}
				tempFile := fmt.Sprintf(download.OutputFile+"-part-%d.tmp", i)
				tempFile = filepath.Join(dm.TempFolder, tempFile)
				download.PartDownloaders = append(
					download.PartDownloaders,
					&PartDownloader{Index: i, Start: start + getFileSize(tempFile), End: end, TempFile: tempFile},
				)
			}
		} else {
			download.PartDownloaders = append(download.PartDownloaders, &PartDownloader{
				Index:    0,
				Start:    0,
				End:      0,
				TempFile: filepath.Join(dm.TempFolder, download.OutputFile+"-part-0.tmp"),
			})
		}
		download.Status = "pending"
	}()
}

func (dm *DownloadManager) PauseDownload(download *Download) {
	download.Status = "paused"
}

func (dm *DownloadManager) ResumeDownload(download *Download) {
	download.Status = "initializing"
}

func (dm *DownloadManager) startDownload(download *Download, StartWG *sync.WaitGroup) {

	var wg sync.WaitGroup
	// fmt.Println("downloading " + download.URL)
	for _, part := range download.PartDownloaders {
		wg.Add(1)
		go func() {
			defer wg.Done()
			download.Queue.PartDownloaders <- part
			StartWG.Done()
			// fmt.Println("part", part.Index, "started")
			err := dm.partDownload(download, part)
			if err != nil {
				// fmt.Println(err)
				part.IsFailed = true
				<-download.Queue.PartDownloaders
				return
			}
			part.IsFailed = false
			// time.Sleep(5 * time.Second)
			// fmt.Println("end of ", download.URL)
			<-download.Queue.PartDownloaders

		}()
	}
	time.Sleep(time.Second * 10)
	// close(download.IsCompletlyStarted)
	download.Temps.StartTime = time.Now()
	download.Status = "downloading"

	go func() {
		wg.Wait()
		IsDone := true
		IsPaused := false
		for _, part := range download.PartDownloaders {
			if part.IsPaused {
				IsPaused = true
				IsDone = false
				break
			}
			if part.IsFailed {
				IsDone = false
				break
			}
		}
		if IsDone {
			download.Status = "finished"
			// fmt.Println(download.URL, "finished")
			mergeParts(download)
		}
		if IsPaused {
			download.Status = "paused"
		}
	}()
}

func (dm *DownloadManager) partDownload(download *Download, partDownloader *PartDownloader) error {

	client := &http.Client{}
	req, err := http.NewRequest("GET", download.URL, nil)
	if err != nil {
		return err
	}
	// fmt.Println(partDownloader.TempFile, partDownloader.Start)
	if partDownloader.Start < partDownloader.End {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", partDownloader.Start, partDownloader.End))
	} else if download.IsPartial {
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	file, err := os.OpenFile(
		partDownloader.TempFile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// file, err := os.Create(partDownloader.TempFile)

	if err != nil {
		return err
	}
	defer file.Close()

	var buf []byte
	bandwidth := download.Queue.MaxBandwidth
	if bandwidth > 0 {
		buf = make([]byte, 1024) // 2^10 or 1 Kb
	} else {
		buf = make([]byte, 1024*1024) // 2^20 or 1 Mb
	}
	startTime := time.Now()
	for {
		if download.Queue.MaxBandwidth != bandwidth {
			time.Sleep(1 * time.Second)
			bandwidth := download.Queue.MaxBandwidth // get new bandwith
			if bandwidth > 0 {
				buf = make([]byte, 1024) // 2^10 or 1 Kb
			} else {
				buf = make([]byte, 1024*1024) // 2^20 or 1 Mb
			}
		}
		if download.Queue.MaxBandwidth > 0 {
			<-download.Queue.tokenBucket
		}
		n, err := resp.Body.Read(buf)
		if n > 0 {
			// limiter.WaitN(context.Background(), n)

			partDownloader.Downloaded += int64(n)
			download.Temps.Mutex.Lock()
			download.Temps.TotalDownloaded += int64(n)
			download.Temps.Mutex.Unlock()
			file.Write(buf[:n])
			elapsed := time.Since(startTime).Seconds()
			partDownloader.Speed = int64(float64(partDownloader.Downloaded) / elapsed)
		}
		if err == io.EOF {
			break
		}
		if !download.Queue.IsActive || download.Status == "paused" {
			partDownloader.IsPaused = true
			break
		}
		if err != nil {
			download.Retries++
			if download.Retries > download.Queue.MaxRetries {
				partDownloader.IsFailed = true
				return err
			}
		}
		if download.Retries > download.Queue.MaxRetries {
			return nil
		}
	}

	return nil
}

func (queue *Queue) SetBandwith(bandwith int) {
	queue.MaxBandwidth = -1
	if queue.ticker != nil {
		queue.ticker.Stop()
	}
	queue.tokenBucket = make(chan struct{}, bandwith)
	go func() {
		tokenInterval := max(1, 1000_000/bandwith)
		queue.ticker = time.NewTicker(time.Duration(tokenInterval) * time.Microsecond)
		queue.MaxBandwidth = bandwith

		defer queue.ticker.Stop()
		for range queue.ticker.C {
			select {
			case queue.tokenBucket <- struct{}{}:
			default:
			}
		}
	}()
}
