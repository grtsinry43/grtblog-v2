package taxonomy

import (
	"context"
	"strings"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
)

type CategoryService struct {
	repo content.CategoryRepository
}

func NewCategoryService(repo content.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) List(ctx context.Context) ([]*content.ArticleCategory, error) {
	return s.repo.List(ctx)
}

func (s *CategoryService) Create(ctx context.Context, name string, shortURL *string) (*content.ArticleCategory, error) {
	category := &content.ArticleCategory{
		Name:     strings.TrimSpace(name),
		ShortURL: trimPtr(shortURL),
	}
	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Update(ctx context.Context, id int64, name string, shortURL *string) (*content.ArticleCategory, error) {
	category := &content.ArticleCategory{
		ID:       id,
		Name:     strings.TrimSpace(name),
		ShortURL: trimPtr(shortURL),
	}
	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
