package manager

import (
	"os"
	"time"
)

func IsWithinActiveHours(start, end string) bool {
	now := time.Now()
	startTime, err := time.Parse("15:04", start)
	if err != nil {
		return true
	}
	endTime, err := time.Parse("15:04", end)
	if err != nil {
		return true
	}
	return now.Hour() >= startTime.Hour() &&
		now.Minute() >= startTime.Minute() &&
		now.Hour() <= endTime.Hour() &&
		now.Minute() <= endTime.Minute()
}

func getFileSize(file string) int64 {
	fileInfo, err := os.Stat(file)
	var size int64 = 0
	if err == nil {
		size = fileInfo.Size()
	}
	return size
}
