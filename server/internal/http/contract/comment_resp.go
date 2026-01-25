package contract

import "time"

type CreateCommentResp struct {
	ID        int64      `json:"id"`
	AreaID    int64      `json:"areaId"`
	Content   string     `json:"content"`
	NickName  *string    `json:"nickName"`
	Location  *string    `json:"location"`
	Platform  *string    `json:"platform"`
	Browser   *string    `json:"browser"`
	Website   *string    `json:"website"`
	IsOwner   bool       `json:"isOwner"`
	IsFriend  bool       `json:"isFriend"`
	IsAuthor  bool       `json:"isAuthor"`
	IsViewed  bool       `json:"isViewed"`
	IsTop     bool       `json:"isTop"`
	ParentID  *int64     `json:"parentId"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

type CommentNodeResp struct {
	ID        int64             `json:"id"`
	AreaID    int64             `json:"areaId"`
	Content   string            `json:"content"`
	NickName  *string           `json:"nickName"`
	Location  *string           `json:"location"`
	Platform  *string           `json:"platform"`
	Browser   *string           `json:"browser"`
	Website   *string           `json:"website"`
	IsOwner   bool              `json:"isOwner"`
	IsFriend  bool              `json:"isFriend"`
	IsAuthor  bool              `json:"isAuthor"`
	IsViewed  bool              `json:"isViewed"`
	IsTop     bool              `json:"isTop"`
	ParentID  *int64            `json:"parentId"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	DeletedAt *time.Time        `json:"deletedAt,omitempty"`
	Children  []CommentNodeResp `json:"children,omitempty"`
}
