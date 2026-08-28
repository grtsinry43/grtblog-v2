package model

import "time"

type ExportRecord struct {
	ID               string     `gorm:"column:id;primaryKey"`
	Filename         string     `gorm:"column:filename"`
	Status           string     `gorm:"column:status"`
	Stage            string     `gorm:"column:stage"`
	TriggerType      string     `gorm:"column:trigger_type"`
	Mode             string     `gorm:"column:mode"`
	SizeBytes        int64      `gorm:"column:size_bytes"`
	SHA256           string     `gorm:"column:sha256"`
	AppVersion       string     `gorm:"column:app_version"`
	SiteName         string     `gorm:"column:site_name"`
	SiteURL          string     `gorm:"column:site_url"`
	ArticleCount     int64      `gorm:"column:article_count"`
	MomentsCount     int64      `gorm:"column:moments_count"`
	PagesCount       int64      `gorm:"column:pages_count"`
	ThinkingsCount   int64      `gorm:"column:thinkings_count"`
	ImageCount       int64      `gorm:"column:image_count"`
	FailedImageCount int64      `gorm:"column:failed_image_count"`
	ErrorMessage     string     `gorm:"column:error_message"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
	StartedAt        *time.Time `gorm:"column:started_at"`
	CompletedAt      *time.Time `gorm:"column:completed_at"`
}

func (ExportRecord) TableName() string { return "export_ops.export_record" }

type ExportDownloadTicket struct {
	TokenHash string    `gorm:"column:token_hash;primaryKey"`
	ExportID  string    `gorm:"column:export_id"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (ExportDownloadTicket) TableName() string { return "export_ops.download_ticket" }
