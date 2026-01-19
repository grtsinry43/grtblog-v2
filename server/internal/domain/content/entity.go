package content

import "time"

type ArticleCategory struct {
	ID        int64
	Name      string
	ShortURL  *string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type TOCNode struct {
	Name     string    `json:"name"`
	Anchor   string    `json:"anchor"`
	Children []TOCNode `json:"children"`
}

type MomentColumn struct {
	ID        int64
	Name      string
	ShortURL  *string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Tag struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type ArticleTag struct {
	ID        int64
	ArticleID int64
	TagID     int64
}

type MomentTopic struct {
	ID       int64
	MomentID int64
	TagID    int64
}

type Article struct {
	ID          int64
	Title       string
	Summary     string
	AISummary   *string
	LeadIn      *string
	TOC         []TOCNode
	Content     string
	ContentHash string
	AuthorID    int64
	Cover       *string
	CategoryID  *int64
	CommentID   *int64
	ShortURL    string
	IsPublished bool
	IsTop       bool
	IsHot       bool
	IsOriginal  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type ArticleMetrics struct {
	ArticleID int64
	Views     int64
	Likes     int
	Comments  int
	UpdatedAt time.Time
}

type Moment struct {
	ID          int64
	Title       string
	Summary     string
	AISummary   *string
	Content     string
	ContentHash string
	AuthorID    int64
	TOC         []TOCNode
	Image       *string
	ColumnID    *int64
	CommentID   *int64
	ShortURL    string
	IsPublished bool
	IsTop       bool
	IsHot       bool
	IsOriginal  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type MomentMetrics struct {
	MomentID  int64
	Views     int64
	Likes     int
	Comments  int
	UpdatedAt time.Time
}

type Page struct {
	ID          int64
	Title       string
	Description *string
	AISummary   *string
	ShortURL    string
	IsEnabled   bool
	IsBuiltin   bool
	TOC         []TOCNode
	Content     string
	ContentHash string
	CommentID   *int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type PageMetrics struct {
	PageID    int64
	Views     int64
	Likes     int
	Comments  int
	UpdatedAt time.Time
}

type Thinking struct {
	ID        int64
	Content   string
	Author    string
	CreatedAt time.Time
}
