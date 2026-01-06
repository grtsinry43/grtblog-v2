package taxonomy

import (
	"context"
	"strings"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
)

type TagService struct {
	repo content.TagRepository
}

func NewTagService(repo content.TagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) List(ctx context.Context) ([]*content.Tag, error) {
	return s.repo.List(ctx)
}

func (s *TagService) Create(ctx context.Context, name string) (*content.Tag, error) {
	tag := &content.Tag{
		Name: strings.TrimSpace(name),
	}
	if err := s.repo.Create(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *TagService) Update(ctx context.Context, id int64, name string) (*content.Tag, error) {
	tag := &content.Tag{
		ID:   id,
		Name: strings.TrimSpace(name),
	}
	if err := s.repo.Update(ctx, tag); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *TagService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
