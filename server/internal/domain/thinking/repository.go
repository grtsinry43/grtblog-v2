package thinking

import "context"

type ThinkingRepository interface {
	FindByID(ctx context.Context, id int64) (*Thinking, error)
	List(ctx context.Context, limit, offset int) ([]*Thinking, int64, error)
	Create(ctx context.Context, thinking *Thinking) error
	Update(ctx context.Context, thinking *Thinking) error
	Delete(ctx context.Context, id int64) error
	IncView(ctx context.Context, id int64) error
	IncLike(ctx context.Context, id int64) error
	DecLike(ctx context.Context, id int64) error
	IncComment(ctx context.Context, id int64) error
	DecComment(ctx context.Context, id int64) error
}
