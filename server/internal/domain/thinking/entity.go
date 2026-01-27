package thinking

import "time"

type ThinkingMetrics struct {
	ThinkingID int64
	Views      int64
	Likes      int
	Comments   int
	UpdatedAt  time.Time
}

type Thinking struct {
	ID        int64
	CommentID int64
	Content   string
	AuthorID  int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Metrics   ThinkingMetrics
}
