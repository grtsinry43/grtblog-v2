package thinking

type CreateThinkingCmd struct {
	Content string
	Author  string
}

type UpdateThinkingCmd struct {
	ID      int64
	Content string
}
