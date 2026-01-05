package websiteinfo

import (
	"context"
	"strings"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/config"
)

// Service 编排 WebsiteInfo 相关用例。
type Service struct {
	repo config.WebsiteInfoRepository
}

func NewService(repo config.WebsiteInfoRepository) *Service {
	return &Service{repo: repo}
}

type CreateCmd struct {
	Key   string
	Value string
}

type UpdateCmd struct {
	Key   string
	Value string
}

func (s *Service) List(ctx context.Context) ([]config.WebsiteInfo, error) {
	return s.repo.List(ctx)
}

func (s *Service) Create(ctx context.Context, cmd CreateCmd) (*config.WebsiteInfo, error) {
	info := &config.WebsiteInfo{
		Key:   strings.TrimSpace(cmd.Key),
		Value: cmd.Value,
	}
	if err := s.repo.Create(ctx, info); err != nil {
		return nil, err
	}
	return info, nil
}

func (s *Service) Update(ctx context.Context, cmd UpdateCmd) (*config.WebsiteInfo, error) {
	info := &config.WebsiteInfo{
		Key:   strings.TrimSpace(cmd.Key),
		Value: cmd.Value,
	}
	if err := s.repo.Update(ctx, info); err != nil {
		return nil, err
	}
	return info, nil
}

func (s *Service) Delete(ctx context.Context, key string) error {
	return s.repo.Delete(ctx, strings.TrimSpace(key))
}

func (s *Service) Get(ctx context.Context, key string) (*config.WebsiteInfo, error) {
	return s.repo.GetByKey(ctx, strings.TrimSpace(key))
}
