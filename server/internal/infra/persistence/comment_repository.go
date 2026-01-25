package persistence

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/comment"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type CommentRepository struct {
	db        *gorm.DB
	commentDB *GormRepository[model.Comment]
	areaDB    *GormRepository[model.CommentArea]
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{
		db:        db,
		commentDB: NewGormRepository[model.Comment](db),
		areaDB:    NewGormRepository[model.CommentArea](db),
	}
}

func (r *CommentRepository) GetAreaByID(ctx context.Context, id int64) (*comment.CommentArea, error) {
	rec, err := r.areaDB.FirstByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, comment.ErrCommentAreaNotFound
		}
		return nil, err
	}
	return mapCommentAreaToDomain(*rec), nil
}

func (r *CommentRepository) FindByID(ctx context.Context, id int64) (*comment.Comment, error) {
	rec, err := r.commentDB.FirstByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, comment.ErrCommentNotFound
		}
		return nil, err
	}
	entity := mapCommentToDomain(*rec)
	return &entity, nil
}

func (r *CommentRepository) ListByAreaID(ctx context.Context, areaID int64) ([]*comment.Comment, error) {
	var recs []model.Comment
	if err := r.db.WithContext(ctx).
		Where("area_id = ?", areaID).
		Order("is_top DESC, created_at ASC").
		Find(&recs).Error; err != nil {
		return nil, err
	}
	out := make([]*comment.Comment, len(recs))
	for i, rec := range recs {
		entity := mapCommentToDomain(rec)
		out[i] = &entity
	}
	return out, nil
}

func (r *CommentRepository) Create(ctx context.Context, commentEntity *comment.Comment) error {
	rec := mapCommentToModel(commentEntity)
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
		return err
	}
	commentEntity.ID = rec.ID
	commentEntity.CreatedAt = rec.CreatedAt
	commentEntity.UpdatedAt = rec.UpdatedAt
	return nil
}

func (r *CommentRepository) Update(ctx context.Context, commentEntity *comment.Comment) error {
	rec := mapCommentToModel(commentEntity)
	return r.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("id = ?", commentEntity.ID).
		Updates(map[string]any{
			"content":    rec.Content,
			"nick_name":  rec.NickName,
			"email":      rec.Email,
			"website":    rec.Website,
			"is_owner":   rec.IsOwner,
			"is_friend":  rec.IsFriend,
			"is_author":  rec.IsAuthor,
			"is_viewed":  rec.IsViewed,
			"is_top":     rec.IsTop,
			"updated_at": time.Now(),
		}).Error
}

func (r *CommentRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Comment{}, id).Error
}

func (r *CommentRepository) SetTopStatus(ctx context.Context, id int64, isTop bool) error {
	return r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ?", id).
		Update("is_top", isTop).Error
}

func mapCommentToDomain(rec model.Comment) comment.Comment {
	return comment.Comment{
		ID:        rec.ID,
		AreaID:    rec.AreaID,
		Content:   rec.Content,
		AuthorID:  rec.AuthorID,
		NickName:  toPtr(rec.NickName),
		IP:        toPtr(rec.IP),
		Location:  toPtr(rec.Location),
		Platform:  toPtr(rec.Platform),
		Browser:   toPtr(rec.Browser),
		Email:     toPtr(rec.Email),
		Website:   toPtr(rec.Website),
		IsOwner:   rec.IsOwner,
		IsFriend:  rec.IsFriend,
		IsAuthor:  rec.IsAuthor,
		IsViewed:  rec.IsViewed,
		IsTop:     rec.IsTop,
		ParentID:  rec.ParentID,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
		DeletedAt: timeToPtr(rec.DeletedAt),
	}
}

func mapCommentToModel(entity *comment.Comment) model.Comment {
	return model.Comment{
		ID:        entity.ID,
		AreaID:    entity.AreaID,
		Content:   strings.TrimSpace(entity.Content),
		AuthorID:  entity.AuthorID,
		NickName:  toValue(entity.NickName),
		IP:        toValue(entity.IP),
		Location:  toValue(entity.Location),
		Platform:  toValue(entity.Platform),
		Browser:   toValue(entity.Browser),
		Email:     toValue(entity.Email),
		Website:   toValue(entity.Website),
		IsOwner:   entity.IsOwner,
		IsFriend:  entity.IsFriend,
		IsAuthor:  entity.IsAuthor,
		IsViewed:  entity.IsViewed,
		IsTop:     entity.IsTop,
		ParentID:  entity.ParentID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: gorm.DeletedAt{Time: timeToValue(entity.DeletedAt), Valid: entity.DeletedAt != nil},
	}
}

func mapCommentAreaToDomain(rec model.CommentArea) *comment.CommentArea {
	return &comment.CommentArea{
		ID:        rec.ID,
		Name:      rec.AreaName,
		Type:      rec.AreaType,
		ContentID: rec.ContentID,
		IsClosed:  rec.IsClosed,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
	}
}

func toPtr(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func toValue(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func timeToPtr(val gorm.DeletedAt) *time.Time {
	if !val.Valid {
		return nil
	}
	t := val.Time
	return &t
}

func timeToValue(val *time.Time) time.Time {
	if val == nil {
		return time.Time{}
	}
	return *val
}
