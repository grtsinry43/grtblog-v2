package thinking

import "errors"

var (
	ErrThinkingNotFound     = errors.New("回想不存在或已被删除")
	ErrThinkingContentEmpty = errors.New("回想内容不能为空")
)
