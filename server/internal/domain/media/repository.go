package media

import "context"

// Repository 定义上传文件的持久化操作。
type Repository interface {
	FindByHash(ctx context.Context, hash string) (*UploadFile, error)
	FindByID(ctx context.Context, id int64) (*UploadFile, error)
	Create(ctx context.Context, file *UploadFile) error
	UpdatePath(ctx context.Context, id int64, path string) error
	UpdateName(ctx context.Context, id int64, name string) error
	List(ctx context.Context, offset int, limit int) ([]UploadFile, int64, error)
	DeleteByID(ctx context.Context, id int64) error
}
