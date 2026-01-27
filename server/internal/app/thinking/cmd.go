package thinking

type CreateThinkingCmd struct {
	Content  string
	AuthorID int64
}

type UpdateThinkingCmd struct {
	ID      int64
	Content string
}
