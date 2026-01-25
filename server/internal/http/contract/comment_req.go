package contract

type CreateCommentLoginReq struct {
	Content  string `json:"content" validate:"required"`
	ParentID *int64 `json:"parentId"`
}

type CreateCommentVisitorReq struct {
	Content  string  `json:"content" validate:"required"`
	NickName *string `json:"nickName" validate:"required,max=255"`
	Email    *string `json:"email" validate:"required,max=255"`
	Website  *string `json:"website" validate:"max=255"`
	ParentID *int64  `json:"parentId"`
}

type UpdateCommentReq struct {
	Content string `json:"content" validate:"required"`
}
