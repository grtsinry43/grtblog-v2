package discovery

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("public discovery resource not found")

// Record is a public-content projection; lists never include the article body.
type Record struct {
	ID         int64
	Kind       string
	Slug       string
	Title      string
	Summary    string
	CreatedAt  time.Time
	ModifiedAt time.Time
	AuthorID   int64
	Author     string
	Content    string
}

type Repository interface {
	List(context.Context) ([]Record, error)
	Document(context.Context, string, string) (Record, error)
}
