package friendlink

import (
	"context"
	"errors"
	"strings"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/social"
)

type Service struct {
	repo social.FriendLinkApplicationRepository
}

func NewService(repo social.FriendLinkApplicationRepository) *Service {
	return &Service{repo: repo}
}

type SubmitCmd struct {
	Name        string
	URL         string
	Logo        string
	Description string
	Message     string
	RSSURL      string
	UserID      *int64
}

type SubmitResult struct {
	Application social.FriendLinkApplication
	Created     bool
}

func (s *Service) Submit(ctx context.Context, cmd SubmitCmd) (*SubmitResult, error) {
	url := strings.TrimSpace(cmd.URL)
	existing, err := s.repo.FindByURL(ctx, url)
	if err != nil && !errors.Is(err, social.ErrFriendLinkApplicationNotFound) {
		return nil, err
	}

	if existing == nil {
		app := &social.FriendLinkApplication{
			Name:              toOptionalString(cmd.Name),
			URL:               url,
			Logo:              toOptionalString(cmd.Logo),
			Description:       toOptionalString(cmd.Description),
			ApplyChannel:      "user",
			RequestedSyncMode: requestedSyncMode(cmd.RSSURL),
			RSSURL:            toOptionalString(cmd.RSSURL),
			UserID:            cmd.UserID,
			Message:           toOptionalString(cmd.Message),
			Status:            "pending",
		}
		if err := s.repo.Create(ctx, app); err != nil {
			return nil, err
		}
		return &SubmitResult{Application: *app, Created: true}, nil
	}

	existing.Name = toOptionalString(cmd.Name)
	existing.Logo = toOptionalString(cmd.Logo)
	existing.Description = toOptionalString(cmd.Description)
	existing.ApplyChannel = "user"
	existing.RequestedSyncMode = requestedSyncMode(cmd.RSSURL)
	existing.RSSURL = toOptionalString(cmd.RSSURL)
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

func requestedSyncMode(rssURL string) string {
	if strings.TrimSpace(rssURL) == "" {
		return "none"
	}
	return "rss"
}
