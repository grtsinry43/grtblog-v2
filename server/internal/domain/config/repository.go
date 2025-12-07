package config

import "context"

// WebsiteInfoRepository 抽象网站信息配置的持久化操作，便于在应用层依赖。
type WebsiteInfoRepository interface {
	List(ctx context.Context) ([]WebsiteInfo, error)
	GetByKey(ctx context.Context, key string) (*WebsiteInfo, error)
	Create(ctx context.Context, info *WebsiteInfo) error
	Update(ctx context.Context, info *WebsiteInfo) error
	Delete(ctx context.Context, key string) error
}
