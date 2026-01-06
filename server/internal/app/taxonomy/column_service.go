package taxonomy

import (
	"context"
	"strings"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
)

type ColumnService struct {
	repo content.ColumnRepository
}

func NewColumnService(repo content.ColumnRepository) *ColumnService {
	return &ColumnService{repo: repo}
}

func (s *ColumnService) List(ctx context.Context) ([]*content.MomentColumn, error) {
	return s.repo.List(ctx)
}

func (s *ColumnService) Create(ctx context.Context, name string, shortURL *string) (*content.MomentColumn, error) {
	column := &content.MomentColumn{
		Name:     strings.TrimSpace(name),
		ShortURL: trimPtr(shortURL),
	}
	if err := s.repo.Create(ctx, column); err != nil {
		return nil, err
	}
	return column, nil
}

func (s *ColumnService) Update(ctx context.Context, id int64, name string, shortURL *string) (*content.MomentColumn, error) {
	column := &content.MomentColumn{
		ID:       id,
		Name:     strings.TrimSpace(name),
		ShortURL: trimPtr(shortURL),
	}
	if err := s.repo.Update(ctx, column); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ColumnService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
