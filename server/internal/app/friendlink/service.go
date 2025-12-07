package friendlink

import (
	"context"
	"strings"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/social"
)

type Service struct {
	repo social.FriendLinkApplicationRepository
}

func NewService(repo social.FriendLinkApplicationRepository) *Service {
	return &Service{repo: repo}
}

type SubmitCommand struct {
	Name        string
	URL         string
	Logo        string
	Description string
	Message     string
	UserID      *int64
}

type SubmitResult struct {
	Application social.FriendLinkApplication
	Created     bool
}

func (s *Service) Submit(ctx context.Context, cmd SubmitCommand) (*SubmitResult, error) {
	url := strings.TrimSpace(cmd.URL)
	existing, err := s.repo.FindByURL(ctx, url)
	if err != nil && err != social.ErrFriendLinkApplicationNotFound {
		return nil, err
	}

	if existing == nil {
		app := &social.FriendLinkApplication{
			Name:        toOptionalString(cmd.Name),
			URL:         url,
			Logo:        toOptionalString(cmd.Logo),
			Description: toOptionalString(cmd.Description),
			UserID:      cmd.UserID,
			Message:     toOptionalString(cmd.Message),
			Status:      "pending",
		}
		if err := s.repo.Create(ctx, app); err != nil {
			return nil, err
		}
		return &SubmitResult{Application: *app, Created: true}, nil
	}

	existing.Name = toOptionalString(cmd.Name)
	existing.Logo = toOptionalString(cmd.Logo)
	existing.Description = toOptionalString(cmd.Description)
	existing.Message = toOptionalString(cmd.Message)
	existing.UserID = cmd.UserID
	existing.Status = "pending"

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return &SubmitResult{Application: *existing, Created: false}, nil
}

func toOptionalString(val string) *string {
	trimmed := strings.TrimSpace(val)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
