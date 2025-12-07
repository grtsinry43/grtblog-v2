package response

type AppError struct {
	Biz     BizError // 对应 ErrorCode
	Message string   // 可覆写文案
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Biz.Msg
}

// NewBizError 类似于: throw new BusinessException(ErrorCode)
func NewBizError(b BizError) *AppError {
	return &AppError{
		Biz:     b,
		Message: "",
	}
}

// NewBizErrorWithMsg 类似: throw new BusinessException(ErrorCode, "自定义提示")
func NewBizErrorWithMsg(b BizError, msg string) *AppError {
	return &AppError{
		Biz:     b,
		Message: msg,
	}
}
