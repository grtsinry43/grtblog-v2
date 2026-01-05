package content

// ArticleListOptions 公开的文章列表查询选项
type ArticleListOptions struct {
	Page       int     `json:"page" validate:"min=1"`
	PageSize   int     `json:"pageSize" validate:"min=1,max=100"`
	CategoryID *int64  `json:"categoryId,omitempty"`
	TagID      *int64  `json:"tagId,omitempty"`
	AuthorID   *int64  `json:"authorId,omitempty"`
	Search     *string `json:"search,omitempty"`
}

// ArticleListOptionsInternal 内部的文章列表查询选项（包含管理功能）
type ArticleListOptionsInternal struct {
	Page       int     `json:"page" validate:"min=1"`
	PageSize   int     `json:"pageSize" validate:"min=1,max=100"`
	CategoryID *int64  `json:"categoryId,omitempty"`
	TagID      *int64  `json:"tagId,omitempty"`
	AuthorID   *int64  `json:"authorId,omitempty"`
	Published  *bool   `json:"published,omitempty"` // 仅管理员可用
	Search     *string `json:"search,omitempty"`
}
