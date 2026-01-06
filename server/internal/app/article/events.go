package article

import "time"

type ArticleCreated struct {
	ID        int64
	AuthorID  int64
	Title     string
	ShortURL  string
	Published bool
	At        time.Time
}

func (e ArticleCreated) Name() string { return "article.created" }
func (e ArticleCreated) OccurredAt() time.Time {
	return e.At
}

type ArticleUpdated struct {
	ID        int64
	AuthorID  int64
	Title     string
	ShortURL  string
	Published bool
	At        time.Time
}

func (e ArticleUpdated) Name() string { return "article.updated" }
func (e ArticleUpdated) OccurredAt() time.Time {
	return e.At
}

type ArticlePublished struct {
	ID       int64
	AuthorID int64
	Title    string
	ShortURL string
	At       time.Time
}

func (e ArticlePublished) Name() string { return "article.published" }
func (e ArticlePublished) OccurredAt() time.Time {
	return e.At
}

type ArticleUnpublished struct {
	ID       int64
	AuthorID int64
	Title    string
	ShortURL string
	At       time.Time
}

func (e ArticleUnpublished) Name() string { return "article.unpublished" }
func (e ArticleUnpublished) OccurredAt() time.Time {
	return e.At
}

type ArticleDeleted struct {
	ID       int64
	AuthorID int64
	Title    string
	ShortURL string
	At       time.Time
}

func (e ArticleDeleted) Name() string { return "article.deleted" }
func (e ArticleDeleted) OccurredAt() time.Time {
	return e.At
}
