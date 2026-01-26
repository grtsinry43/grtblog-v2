package navigation

import (
	"context"
	"strings"

	domain "github.com/grtsinry43/grtblog-v2/server/internal/domain/navigation"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]*domain.NavMenu, error) {
	return s.repo.List(ctx)
}

func (s *Service) Create(ctx context.Context, cmd CreateNavMenuCmd) (*domain.NavMenu, error) {
	if cmd.ParentID != nil {
		if _, err := s.repo.GetByID(ctx, *cmd.ParentID); err != nil {
			return nil, err
		}
	}

	sort, err := s.repo.NextSort(ctx, cmd.ParentID)
	if err != nil {
		return nil, err
	}

	icon := normalizeIcon(cmd.Icon)

	menu := &domain.NavMenu{
		Name:     strings.TrimSpace(cmd.Name),
		URL:      strings.TrimSpace(cmd.URL),
		Icon:     icon,
		Sort:     sort,
		ParentID: cmd.ParentID,
	}

	if err := s.repo.Create(ctx, menu); err != nil {
		return nil, err
	}

	return menu, nil
}

func (s *Service) Update(ctx context.Context, cmd UpdateNavMenuCmd) (*domain.NavMenu, error) {
	menu, err := s.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	menu.Name = strings.TrimSpace(cmd.Name)
	menu.URL = strings.TrimSpace(cmd.URL)

	if cmd.Icon != nil {
		menu.Icon = normalizeIcon(cmd.Icon)
	}

	if cmd.ParentID != nil {
		if menu.ParentID == nil || *menu.ParentID != *cmd.ParentID {
			if _, err := s.repo.GetByID(ctx, *cmd.ParentID); err != nil {
				return nil, err
			}
			menu.ParentID = cmd.ParentID
			if cmd.Sort != nil {
				menu.Sort = *cmd.Sort
			} else {
				nextSort, err := s.repo.NextSort(ctx, cmd.ParentID)
				if err != nil {
					return nil, err
				}
				menu.Sort = nextSort
			}
		}
	} else if cmd.Sort != nil {
		menu.Sort = *cmd.Sort
	}

	if err := s.repo.Update(ctx, menu); err != nil {
		return nil, err
	}

	return menu, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) UpdateOrder(ctx context.Context, items []NavMenuOrderItem) error {
	updates := make([]domain.NavMenuOrderUpdate, 0, len(items))
	for _, item := range items {
		updates = append(updates, domain.NavMenuOrderUpdate{
			ID:       item.ID,
			ParentID: item.ParentID,
			Sort:     item.Sort,
		})
	}
	return s.repo.UpdateOrder(ctx, updates)
}

func normalizeIcon(icon *string) *string {
	if icon == nil {
		return nil
	}
	value := strings.TrimSpace(*icon)
	if value == "" {
		return nil
	}
	return &value
}
