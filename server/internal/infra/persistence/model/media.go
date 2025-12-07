package model

import "time"

type UploadFile struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;size:255;not null"`
	Path      string    `gorm:"column:path;size:255;not null"`
	Type      string    `gorm:"column:type;size:45;not null"`
	Size      int64     `gorm:"column:size;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (UploadFile) TableName() string { return "upload_file" }

type Photo struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	URL         string    `gorm:"column:url;size:255;not null"`
	Location    string    `gorm:"column:location;size:255"`
	Device      string    `gorm:"column:device;size:255"`
	Shade       string    `gorm:"column:shade;size:255"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (Photo) TableName() string { return "photo" }
