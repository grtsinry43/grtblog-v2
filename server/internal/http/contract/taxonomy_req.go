package contract

// CategoryCreateReq 文章分类创建请求。
type CategoryCreateReq struct {
	Name     string  `json:"name"`
	ShortURL *string `json:"shortUrl,omitempty"`
}

// CategoryUpdateReq 文章分类更新请求。
type CategoryUpdateReq struct {
	Name     string  `json:"name"`
	ShortURL *string `json:"shortUrl,omitempty"`
}

// ColumnCreateReq 手记分区创建请求。
type ColumnCreateReq struct {
	Name     string  `json:"name"`
	ShortURL *string `json:"shortUrl,omitempty"`
}

// ColumnUpdateReq 手记分区更新请求。
type ColumnUpdateReq struct {
	Name     string  `json:"name"`
	ShortURL *string `json:"shortUrl,omitempty"`
}

// TagCreateReq 标签创建请求。
type TagCreateReq struct {
	Name string `json:"name"`
}

// TagUpdateReq 标签更新请求。
type TagUpdateReq struct {
	Name string `json:"name"`
}
