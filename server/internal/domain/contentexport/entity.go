// Package contentexport 定义内置内容导出（markdown + 打包图片）的领域模型。
package contentexport

import "time"

type Status string

const (
	StatusQueued    Status = "queued"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

// Mode 控制导出包的内容布局。
type Mode string

const (
	ModeStructured Mode = "structured"
	ModeFlatten    Mode = "flatten"
	ModeBoth       Mode = "both"
)

// ValidMode 校验模式取值。
func ValidMode(mode string) (Mode, bool) {
	switch Mode(mode) {
	case ModeStructured, ModeFlatten, ModeBoth:
		return Mode(mode), true
	default:
		return "", false
	}
}

type Record struct {
	ID               string     `json:"id"`
	Filename         string     `json:"filename"`
	Status           Status     `json:"status"`
	Stage            string     `json:"stage"`
	TriggerType      string     `json:"triggerType"`
	Mode             string     `json:"mode"`
	SizeBytes        int64      `json:"sizeBytes"`
	SHA256           string     `json:"sha256,omitempty"`
	AppVersion       string     `json:"appVersion,omitempty"`
	SiteName         string     `json:"siteName,omitempty"`
	SiteURL          string     `json:"siteUrl,omitempty"`
	ArticleCount     int64      `json:"articleCount"`
	MomentsCount     int64      `json:"momentsCount"`
	PagesCount       int64      `json:"pagesCount"`
	ThinkingsCount   int64      `json:"thinkingsCount"`
	ImageCount       int64      `json:"imageCount"`
	FailedImageCount int64      `json:"failedImageCount"`
	ErrorMessage     string     `json:"errorMessage,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
	StartedAt        *time.Time `json:"startedAt,omitempty"`
	CompletedAt      *time.Time `json:"completedAt,omitempty"`
}

type DownloadTicket struct {
	TokenHash string
	ExportID  string
	ExpiresAt time.Time
	CreatedAt time.Time
}
