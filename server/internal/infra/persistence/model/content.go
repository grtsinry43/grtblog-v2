package model

import (
	"time"

	"gorm.io/gorm"
)

type ArticleCategory struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	Name      string         `gorm:"column:name;size:45;not null"`
	ShortURL  string         `gorm:"column:short_url;size:255"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (ArticleCategory) TableName() string { return "article_category" }

type MomentColumn struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	Name      string         `gorm:"column:name;size:45;not null"`
	ShortURL  string         `gorm:"column:short_url;size:255"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (MomentColumn) TableName() string { return "moment_column" }

type Tag struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	Name      string         `gorm:"column:name;size:45;not null"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Tag) TableName() string { return "tag" }

type ArticleTag struct {
	ID        int64 `gorm:"column:id;primaryKey"`
	ArticleID int64 `gorm:"column:article_id;not null"`
	TagID     int64 `gorm:"column:tag_id;not null"`
}

func (ArticleTag) TableName() string { return "article_tag" }

type MomentTopic struct {
	ID       int64 `gorm:"column:id;primaryKey"`
	MomentID int64 `gorm:"column:moment_id;not null"`
	TagID    int64 `gorm:"column:tag_id;not null"`
}

func (MomentTopic) TableName() string { return "moment_topic" }

type Article struct {
	ID          int64          `gorm:"column:id;primaryKey"`
	Title       string         `gorm:"column:title;size:255;not null"`
	Summary     string         `gorm:"column:summary;type:text;not null"`
	AISummary   *string        `gorm:"column:ai_summary;type:text"`
	LeadIn      *string        `gorm:"column:lead_in;type:text"`
	TOC         []byte         `gorm:"column:toc;type:jsonb;not null"`
	Content     string         `gorm:"column:content;type:text;not null"`
	AuthorID    int64          `gorm:"column:author_id;not null"`
	Cover       *string        `gorm:"column:cover;size:255"`
	CategoryID  *int64         `gorm:"column:category_id"`
	CommentID   *int64         `gorm:"column:comment_id"`
	ShortURL    string         `gorm:"column:short_url;size:255;not null"`
	IsPublished bool           `gorm:"column:is_published"`
	IsTop       bool           `gorm:"column:is_top"`
	IsHot       bool           `gorm:"column:is_hot"`
	IsOriginal  bool           `gorm:"column:is_original"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Article) TableName() string { return "article" }

type ArticleMetrics struct {
	ArticleID int64     `gorm:"column:article_id;primaryKey"`
	Views     int64     `gorm:"column:views"`
	Likes     int       `gorm:"column:likes"`
	Comments  int       `gorm:"column:comments"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ArticleMetrics) TableName() string { return "article_metrics" }

type Moment struct {
	ID          int64          `gorm:"column:id;primaryKey"`
	Title       string         `gorm:"column:title;size:255;not null"`
	Summary     string         `gorm:"column:summary;type:text;not null"`
	AISummary   string         `gorm:"column:ai_summary;type:text"`
	Content     string         `gorm:"column:content;type:text;not null"`
	AuthorID    int64          `gorm:"column:author_id;not null"`
	TOC         []byte         `gorm:"column:toc;type:jsonb;not null"`
	Image       string         `gorm:"column:img"`
	ColumnID    *int64         `gorm:"column:column_id"`
	CommentID   *int64         `gorm:"column:comment_id"`
	ShortURL    string         `gorm:"column:short_url;size:255;not null"`
	IsPublished bool           `gorm:"column:is_published"`
	IsTop       bool           `gorm:"column:is_top"`
	IsHot       bool           `gorm:"column:is_hot"`
	IsOriginal  bool           `gorm:"column:is_original"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Moment) TableName() string { return "moment" }

type MomentMetrics struct {
	MomentID  int64     `gorm:"column:moment_id;primaryKey"`
	Views     int64     `gorm:"column:views"`
	Likes     int       `gorm:"column:likes"`
	Comments  int       `gorm:"column:comments"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (MomentMetrics) TableName() string { return "moment_metrics" }

type Page struct {
	ID          int64          `gorm:"column:id;primaryKey"`
	Title       string         `gorm:"column:title;size:255;not null"`
	Description string         `gorm:"column:description;size:255"`
	AISummary   string         `gorm:"column:ai_summary;type:text"`
	ShortURL    string         `gorm:"column:short_url;size:255;not null"`
	IsEnabled   bool           `gorm:"column:is_enabled"`
	IsBuiltin   bool           `gorm:"column:is_builtin"`
	TOC         []byte         `gorm:"column:toc;type:jsonb;not null"`
	Content     string         `gorm:"column:content;type:text;not null"`
	CommentID   *int64         `gorm:"column:comment_id"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Page) TableName() string { return "page" }

type PageMetrics struct {
	PageID    int64     `gorm:"column:page_id;primaryKey"`
	Views     int64     `gorm:"column:views"`
	Likes     int       `gorm:"column:likes"`
	Comments  int       `gorm:"column:comments"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (PageMetrics) TableName() string { return "page_metrics" }

type Thinking struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Content   string    `gorm:"column:content;type:text;not null"`
	Author    string    `gorm:"column:author;size:45;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (Thinking) TableName() string { return "thinking" }
