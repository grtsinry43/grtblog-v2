package persistence

import (
	"context"
	"errors"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/discovery"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
	"gorm.io/gorm"
)

type DiscoveryRepository struct{ db *gorm.DB }

func NewDiscoveryRepository(db *gorm.DB) *DiscoveryRepository { return &DiscoveryRepository{db: db} }

func (r *DiscoveryRepository) query(ctx context.Context, kind string) (*gorm.DB, error) {
	const fields = "short_url AS slug, title, created_at, content_updated_at AS modified_at"
	switch kind {
	case "posts":
		return r.db.WithContext(ctx).Model(&model.Article{}).Where("is_published = ?", true).Select(fields + ", summary, author_id"), nil
	case "moments":
		return r.db.WithContext(ctx).Model(&model.Moment{}).Where("is_published = ?", true).Select(fields + ", summary, author_id"), nil
	case "pages":
		return r.db.WithContext(ctx).Model(&model.Page{}).Where("is_enabled = ?", true).Select(fields + ", description AS summary"), nil
	default:
		return nil, discovery.ErrNotFound
	}
}

func (r *DiscoveryRepository) List(ctx context.Context) ([]discovery.Record, error) {
	result := make([]discovery.Record, 0)
	for _, kind := range []string{"posts", "moments", "pages"} {
		query, err := r.query(ctx, kind)
		if err != nil {
			return nil, err
		}
		var rows []discovery.Record
		if err := query.Order("created_at DESC, id DESC").Find(&rows).Error; err != nil {
			return nil, err
		}
		for _, row := range rows {
			row.Kind = kind
			result = append(result, row)
		}
	}
	for _, kind := range []string{"categories", "columns"} {
		var source interface{} = &model.ArticleCategory{}
		if kind == "columns" {
			source = &model.MomentColumn{}
		}
		var rows []discovery.Record
		if err := r.db.WithContext(ctx).Model(source).Select("short_url AS slug, name AS title").Order("id").Find(&rows).Error; err != nil {
			return nil, err
		}
		for _, row := range rows {
			row.Kind = kind
			result = append(result, row)
		}
	}
	var albums []discovery.Record
	if err := r.db.WithContext(ctx).Model(&model.Album{}).Where("is_published = ?", true).
		Select("id, short_url AS slug, title, description AS summary").Order("id").Find(&albums).Error; err != nil {
		return nil, err
	}
	for _, row := range albums {
		row.Kind = "albums"
		result = append(result, row)
	}
	var photos []discovery.Record
	if err := r.db.WithContext(ctx).Model(&model.Photo{}).
		Joins("JOIN album ON album.id = photo.album_id AND album.is_published = TRUE AND album.deleted_at IS NULL").
		Select("photo.id, album.short_url AS slug, COALESCE(NULLIF(photo.caption, ''), album.title) AS title").
		Order("photo.id").Find(&photos).Error; err != nil {
		return nil, err
	}
	for _, row := range photos {
		row.Kind = "photos"
		result = append(result, row)
	}
	return result, nil
}

func (r *DiscoveryRepository) Document(ctx context.Context, kind, slug string) (discovery.Record, error) {
	query, err := r.query(ctx, kind)
	if err != nil {
		return discovery.Record{}, err
	}
	// Keep the public predicate and GORM soft-delete scope identical to List.
	fields := append(append([]string{}, query.Statement.Selects...), "content")
	var row discovery.Record
	err = query.Select(fields).Where("short_url = ?", slug).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return row, discovery.ErrNotFound
	}
	if err != nil {
		return row, err
	}
	row.Kind = kind
	if row.AuthorID > 0 {
		var author model.User
		err = r.db.WithContext(ctx).Select("id", "nickname", "username").Where("id = ?", row.AuthorID).Take(&author).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return row, err
		}
		row.Author = author.Nickname
		if row.Author == "" {
			row.Author = author.Username
		}
	}
	return row, nil
}
