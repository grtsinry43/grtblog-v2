package comment

import "time"

type CommentArea struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Comment struct {
	ID        int64
	AreaID    int64
	Content   string
	AuthorID  *int64
	NickName  *string
	IP        *string
	Location  *string
	Platform  *string
	Browser   *string
	Email     *string
	Website   *string
	IsOwner   bool
	IsFriend  bool
	IsAuthor  bool
	IsViewed  bool
	IsTop     bool
	ParentID  *int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
