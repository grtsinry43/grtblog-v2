package social

import "time"

type FriendLink struct {
	ID          int64
	Name        string
	URL         string
	Logo        *string
	Description *string
	RSSURL      *string
	UserID      *int64
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type FriendLinkApplication struct {
	ID          int64
	Name        *string
	URL         string
	Logo        *string
	Description *string
	UserID      *int64
	Message     *string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GlobalNotification struct {
	ID         int64
	Content    string
	PublishAt  time.Time
	ExpireAt   time.Time
	AllowClose bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
