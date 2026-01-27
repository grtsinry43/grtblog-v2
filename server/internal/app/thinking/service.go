package thinking

import (
	"context"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/comment"
	domainthinking "github.com/grtsinry43/grtblog-v2/server/internal/domain/thinking"
)

type Service struct {
	repo        domainthinking.ThinkingRepository
	commentRepo comment.CommentRepository
}

func NewService(repo domainthinking.ThinkingRepository, commentRepo comment.CommentRepository) *Service {
	return &Service{
		repo:        repo,
		commentRepo: commentRepo,
	}
}

func (s *Service) Create(ctx context.Context, cmd CreateThinkingCmd) (*domainthinking.Thinking, error) {
	if cmd.Content == "" {
		return nil, domainthinking.ErrThinkingContentEmpty
	}
	author := cmd.Author
	if author == "" {
		author = "Admin"
	}

	t := &domainthinking.Thinking{
		Content: cmd.Content,
		Author:  author,
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Service) Update(ctx context.Context, cmd UpdateThinkingCmd) (*domainthinking.Thinking, error) {
	t, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	t.Content = cmd.Content
	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]*domainthinking.Thinking, int64, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, t.ID)
}

func (s *Service) FindByID(ctx context.Context, id int64) (*domainthinking.Thinking, error) {
	_ = s.repo.IncView(ctx, id)
	return s.repo.FindByID(ctx, id)
}
