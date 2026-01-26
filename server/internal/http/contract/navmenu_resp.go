package contract

import "time"

type NavMenuResp struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	URL       string        `json:"url"`
	Icon      *string       `json:"icon,omitempty"`
	Sort      int           `json:"sort"`
	ParentID  *int64        `json:"parentId,omitempty"`
	Children  []NavMenuResp `json:"children,omitempty"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}
