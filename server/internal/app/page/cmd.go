package page

import "time"

// CreatePageCmd 创建页面命令。
type CreatePageCmd struct {
	Title       string
	Description *string
	Content     string
	ShortURL    *string
	IsEnabled   bool
	IsBuiltin   bool
	CreatedAt   *time.Time
}

// UpdatePageCmd 更新页面命令。
type UpdatePageCmd struct {
	ID          int64
	Title       string
	Description *string
	Content     string
	ShortURL    string
	IsEnabled   bool
	IsBuiltin   bool
}
