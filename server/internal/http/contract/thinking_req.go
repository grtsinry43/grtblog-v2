package contract

type CreateThinkingReq struct {
	Content string `json:"content" validate:"required"`
}

type UpdateThinkingReq struct {
	Content string `json:"content" validate:"required"`
}
