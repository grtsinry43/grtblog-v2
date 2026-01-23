package page

import (
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
)

type PageCreated struct {
	ID       int64
	Title    string
	ShortURL string
	Enabled  bool
	At       time.Time
}

func (e PageCreated) Name() string { return "page.created" }
func (e PageCreated) OccurredAt() time.Time {
	return e.At
}

type PageUpdated struct {
	ID          int64
	Title       string
	ShortURL    string
	Enabled     bool
	ContentHash string
	Description *string
	TOC         []content.TOCNode
	Content     string
	At          time.Time
}

func (e PageUpdated) Name() string { return "page.updated" }
func (e PageUpdated) OccurredAt() time.Time {
	return e.At
}

type PageDeleted struct {
	ID       int64
	Title    string
	ShortURL string
	At       time.Time
}

func (e PageDeleted) Name() string { return "page.deleted" }
func (e PageDeleted) OccurredAt() time.Time {
	return e.At
}
