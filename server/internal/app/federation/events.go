package federation

import "time"

type MentionDetected struct {
	ArticleID      int64
	AuthorID       int64
	Title          string
	ShortURL       string
	TargetUser     string
	TargetInstance string
	Context        string
	MentionType    string
	At             time.Time
}

func (e MentionDetected) Name() string { return "federation.mention.detected" }
func (e MentionDetected) OccurredAt() time.Time {
	return e.At
}

type CitationDetected struct {
	ArticleID      int64
	AuthorID       int64
	Title          string
	ShortURL       string
	TargetInstance string
	TargetPostID   string
	Context        string
	CitationType   string
	At             time.Time
}

func (e CitationDetected) Name() string { return "federation.citation.detected" }
func (e CitationDetected) OccurredAt() time.Time {
	return e.At
}
