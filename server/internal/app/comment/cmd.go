package comment

type CreateCommentLoginCmd struct {
	AreaID   int64
	Content  string
	ParentID *int64
}

type CreateCommentVisitorCmd struct {
	AreaID   int64
	Content  string
	ParentID *int64
	NickName string
	Email    string
	Website  *string
}
