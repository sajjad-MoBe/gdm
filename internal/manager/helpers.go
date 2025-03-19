package manager

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"
)

func IsWithinActiveHours(start, end string) bool {

	now := time.Now()
	//fmt.Println("Current time:", now)
	startTime, err := time.Parse("15:04", start)
	if err != nil {
		return true
	}
	endTime, err := time.Parse("15:04", end)
	if err != nil {
		return true
	}
	return (now.Hour() > startTime.Hour() ||
		(now.Minute() >= startTime.Minute() && now.Hour() == startTime.Hour())) &&
		(now.Hour() < endTime.Hour() ||
			(now.Hour() == endTime.Hour() && now.Minute() <= endTime.Minute()))
}

func getFileSize(file string) int64 {
	fileInfo, err := os.Stat(file)
	var size int64 = 0
	if err == nil {
		size = fileInfo.Size()
	}
	return size
}

func mergeParts(download *Download) error {
	if err := os.MkdirAll(download.Queue.SaveDir, os.ModePerm); err != nil {
		return err
	}
	fullPath := filepath.Join(download.Queue.SaveDir, download.OutputFile)
	counter := 1
	for {
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			break
		}
		newFilename := fmt.Sprintf("%s(%d)%s",
			download.OutputFile[:len(download.OutputFile)-len(filepath.Ext(download.OutputFile))],
			counter, filepath.Ext(download.OutputFile),
		)
		fullPath = filepath.Join(download.Queue.SaveDir, newFilename)
		counter++
	}
	outFile, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	for _, p := range download.PartDownloaders {
		partFile, err := os.Open(p.TempFile)
		if err != nil {
			return err
		}
		defer os.Remove(p.TempFile)
		defer partFile.Close()

		_, err = io.Copy(outFile, partFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetFileNameFromURL(URL string) (string, error) {
	if len(URL) < 2 {
		return "", errors.New("not enough lenght for URL")
	}
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	return path.Base(parsedURL.Path), nil

}

func IsValidURL(URL string) bool {

	_, err := url.Parse(URL)
	return err == nil

}
