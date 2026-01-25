package comment

import "context"

type CommentRepository interface {
	GetAreaByID(ctx context.Context, id int64) (*CommentArea, error)
	FindByID(ctx context.Context, id int64) (*Comment, error)
	ListByAreaID(ctx context.Context, areaID int64) ([]*Comment, error)
	Create(ctx context.Context, comment *Comment) error
	Update(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, id int64) error
	SetTopStatus(ctx context.Context, id int64, isTop bool) error
}
