package db

type Queue struct {
	ID                     uint   `gorm:"primaryKey" json:"id"`
	SaveDir                string `json:"save_dir"`
	MaxConcurrentDownloads int    `gorm:"default:5" json:"max_concurrent_downloads"`
	MaxBandwidth           int    `gorm:"default:-1" json:"max_bandwidth"` // -1 for infinite
	ActiveStartTime        string `gorm:"default:'00:00'" json:"active_start_time"`
	ActiveEndTime          string `gorm:"default:'23:59'" json:"active_end_time"`
	MaxRetries             int    `gorm:"default:0" json:"max_retries"`
}

type Download struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	QueueID uint   `gorm:"not null" json:"queue_id"`
	Status  string `json:"status"`
	URL     string `json:"url"`
	Retries int    `gorm:"default:0" json:"retries"`

	Queue Queue `gorm:"foreignKey:QueueID"`
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
