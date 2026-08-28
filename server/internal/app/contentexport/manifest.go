package contentexport

import "time"

const ArchiveFormatVersion = 1

// FailedDownload 记录一次失败的外链图片下载。
type FailedDownload struct {
	URL          string `json:"url"`
	ReferencedBy string `json:"referencedBy,omitempty"`
	Error        string `json:"error"`
}

// Manifest 描述一次导出包的内容与完整性信息。
type Manifest struct {
	FormatVersion int            `json:"formatVersion"`
	ExportID      string         `json:"exportId"`
	Mode          string         `json:"mode"`
	CreatedAt     time.Time      `json:"createdAt"`
	AppVersion    string         `json:"appVersion"`
	SiteName      string         `json:"siteName"`
	SiteURL       string         `json:"siteUrl"`
	Counts        ManifestCounts `json:"counts"`
	Images        ManifestImages `json:"images"`
	// ContainsSensitive 恒为 true：admin 全集导出包含未发布内容。
	ContainsSensitive bool              `json:"containsSensitive"`
	Checksums         map[string]string `json:"checksums"`
}

type ManifestCounts struct {
	Articles   int64 `json:"articles"`
	Moments    int64 `json:"moments"`
	Pages      int64 `json:"pages"`
	Thinkings  int64 `json:"thinkings"`
	Categories int64 `json:"categories"`
	Columns    int64 `json:"columns"`
	Tags       int64 `json:"tags"`
	Total      int64 `json:"total"`
}

type ManifestImages struct {
	SelfHosted      int64            `json:"selfHosted"`
	External        int64            `json:"external"`
	FailedCount     int64            `json:"failedCount"`
	FailedDownloads []FailedDownload `json:"failedDownloads"`
}
