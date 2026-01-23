package moment

import (
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
)

type MomentCreated struct {
	ID        int64
	AuthorID  int64
	Title     string
	ShortURL  string
	Published bool
	At        time.Time
}

func (e MomentCreated) Name() string { return "moment.created" }
func (e MomentCreated) OccurredAt() time.Time {
	return e.At
}

type MomentUpdated struct {
	ID          int64
	AuthorID    int64
	Title       string
	ShortURL    string
	Published   bool
	ContentHash string
	Summary     string
	TOC         []content.TOCNode
	Content     string
	At          time.Time
}

func (e MomentUpdated) Name() string { return "moment.updated" }
func (e MomentUpdated) OccurredAt() time.Time {
	return e.At
}

type MomentPublished struct {
	ID       int64
	AuthorID int64
	Title    string
	ShortURL string
	At       time.Time
}

func (e MomentPublished) Name() string { return "moment.published" }
func (e MomentPublished) OccurredAt() time.Time {
	return e.At
}

type MomentUnpublished struct {
	ID       int64
	AuthorID int64
	Title    string
	ShortURL string
	At       time.Time
}

func (e MomentUnpublished) Name() string { return "moment.unpublished" }
func (e MomentUnpublished) OccurredAt() time.Time {
	return e.At
}

type MomentDeleted struct {
	ID       int64
	AuthorID int64
	Title    string
	ShortURL string
	At       time.Time
}

func (e MomentDeleted) Name() string { return "moment.deleted" }
func (e MomentDeleted) OccurredAt() time.Time {
	return e.At
}
