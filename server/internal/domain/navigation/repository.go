package navigation

import "context"

type Repository interface {
	List(ctx context.Context) ([]*NavMenu, error)
	GetByID(ctx context.Context, id int64) (*NavMenu, error)
	Create(ctx context.Context, menu *NavMenu) error
	Update(ctx context.Context, menu *NavMenu) error
	Delete(ctx context.Context, id int64) error
	NextSort(ctx context.Context, parentID *int64) (int, error)
	UpdateOrder(ctx context.Context, updates []NavMenuOrderUpdate) error
}

type NavMenuOrderUpdate struct {
	ID       int64
	ParentID *int64
	Sort     int
}
