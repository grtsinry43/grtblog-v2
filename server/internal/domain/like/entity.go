package like

import "time"

type TargetType string

const (
	TargetArticle TargetType = "article"
	TargetMoment  TargetType = "moment"
	TargetPage    TargetType = "page"
)

type ContentLike struct {
	ID         int64
	TargetType TargetType
	TargetID   int64
	UserID     *int64
	SessionID  *string
	CreatedAt  time.Time
}
