package persistence

import (
	"context"
	"errors"
	"strconv"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/contentutil"
	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/thinking"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type ThinkingRepository struct {
	db   *gorm.DB
	repo *GormRepository[model.Thinking]
}

func NewThinkingRepository(db *gorm.DB) *ThinkingRepository {
	return &ThinkingRepository{
		db:   db,
		repo: NewGormRepository[model.Thinking](db),
	}
}

func (r *ThinkingRepository) FindByID(ctx context.Context, id int64) (*thinking.Thinking, error) {
	var m model.Thinking
	if err := r.db.WithContext(ctx).Preload("Metrics").First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, thinking.ErrThinkingNotFound
		}
		return nil, err
	}
	return mapModelToThinking(&m), nil
}

func (r *ThinkingRepository) List(ctx context.Context, limit, offset int) ([]*thinking.Thinking, int64, error) {
	var models []model.Thinking
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Thinking{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Preload("Metrics").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*thinking.Thinking, len(models))
	for i, m := range models {
		result[i] = mapModelToThinking(&m)
	}
	return result, total, nil
}

func (r *ThinkingRepository) Create(ctx context.Context, t *thinking.Thinking) error {
	m := mapThinkingToModel(t)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&m).Error; err != nil {
			return err
		}

		areaID, err := createCommentArea(tx, contentutil.CommentAreaTypeThinking, "思考", strconv.FormatInt(m.ID, 10), m.ID)
		if err != nil {
			return err
		}
		if err := tx.Model(&model.Thinking{}).Where("id = ?", m.ID).
			Update("comment_id", areaID).Error; err != nil {
			return err
		}
		m.CommentID = areaID

		metrics := model.ThinkingMetrics{
			ThinkingID: m.ID,
		}
		if err := tx.Create(&metrics).Error; err != nil {
			return err
		}

		t.ID = m.ID
		t.CommentID = m.CommentID
		t.CreatedAt = m.CreatedAt
		t.UpdatedAt = m.UpdatedAt
		t.Metrics = thinking.ThinkingMetrics{
			ThinkingID: metrics.ThinkingID,
			Views:      metrics.Views,
			Likes:      metrics.Likes,
			Comments:   metrics.Comments,
			UpdatedAt:  metrics.UpdatedAt,
		}
		return nil
	})
}

func (r *ThinkingRepository) Update(ctx context.Context, t *thinking.Thinking) error {
	return r.db.WithContext(ctx).Model(&model.Thinking{}).
		Where("id = ?", t.ID).
		Updates(map[string]interface{}{
			"content": t.Content,
		}).Error
}

func (r *ThinkingRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var rec model.Thinking
		if err := tx.Select("id", "comment_id").Where("id = ?", id).First(&rec).Error; err != nil {
			return err
		}
		if rec.CommentID != 0 {
			if err := deleteCommentArea(tx, rec.CommentID); err != nil {
				return err
			}
		}
		if err := tx.Delete(&model.ThinkingMetrics{}, "thinking_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.Thinking{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *ThinkingRepository) IncView(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.ThinkingMetrics{}).
		Where("thinking_id = ?", id).
		UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}

func (r *ThinkingRepository) IncLike(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.ThinkingMetrics{}).
		Where("thinking_id = ?", id).
		UpdateColumn("likes", gorm.Expr("likes + ?", 1)).Error
}

func (r *ThinkingRepository) DecLike(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.ThinkingMetrics{}).
		Where("thinking_id = ?", id).
		UpdateColumn("likes", gorm.Expr("likes - ?", 1)).Error
}

func (r *ThinkingRepository) IncComment(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.ThinkingMetrics{}).
		Where("thinking_id = ?", id).
		UpdateColumn("comments", gorm.Expr("comments + ?", 1)).Error
}

func (r *ThinkingRepository) DecComment(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.ThinkingMetrics{}).
		Where("thinking_id = ?", id).
		UpdateColumn("comments", gorm.Expr("comments - ?", 1)).Error
}

func mapModelToThinking(m *model.Thinking) *thinking.Thinking {
	return &thinking.Thinking{
		ID:        m.ID,
		CommentID: m.CommentID,
		Content:   m.Content,
		AuthorID:  m.AuthorID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Metrics: thinking.ThinkingMetrics{
			ThinkingID: m.Metrics.ThinkingID,
			Views:      m.Metrics.Views,
			Likes:      m.Metrics.Likes,
			Comments:   m.Metrics.Comments,
			UpdatedAt:  m.Metrics.UpdatedAt,
		},
	}
}

func mapThinkingToModel(t *thinking.Thinking) model.Thinking {
	return model.Thinking{
		ID:        t.ID,
		CommentID: t.CommentID,
		Content:   t.Content,
		AuthorID:  t.AuthorID,
	}
}
