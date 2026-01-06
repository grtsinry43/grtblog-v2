package response

type AppError struct {
	Biz     BizError // 对应 ErrorCode
	Message string   // 可覆写文案
	Cause   error    // 方便日志记录的原始错误
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
		Cause:   nil,
	}
}

// NewBizErrorWithMsg 类似: throw new BusinessException(ErrorCode, "自定义提示")
func NewBizErrorWithMsg(b BizError, msg string) *AppError {
	return &AppError{
		Biz:     b,
		Message: msg,
		Cause:   nil,
	}
}

// NewBizErrorWithCause 允许携带原始错误，便于记录日志。
func NewBizErrorWithCause(b BizError, msg string, cause error) *AppError {
	return &AppError{
		Biz:     b,
		Message: msg,
		Cause:   cause,
	}
}
