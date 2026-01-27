package contract

type CreateThinkingReq struct {
	Content string `json:"content" validate:"required"`
	Author  string `json:"author"`
}

type UpdateThinkingReq struct {
	ID      int64  `json:"id" validate:"required"`
	Content string `json:"content" validate:"required"`
}
