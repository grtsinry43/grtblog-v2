package response

import "github.com/gofiber/fiber/v2"

type BizError struct {
	HTTPStatus int
	Code       int
	BizErr     string
	Msg        string
}

var (
	OK = BizError{
		HTTPStatus: fiber.StatusOK,
		Code:       0,
		BizErr:     "OK",
		Msg:        "success",
	}

	NotFound = BizError{
		HTTPStatus: fiber.StatusNotFound,
		Code:       404,
		BizErr:     "NOT_FOUND",
		Msg:        "资源未找到",
	}

	MethodNotAllowed = BizError{
		HTTPStatus: fiber.StatusMethodNotAllowed,
		Code:       405,
		BizErr:     "METHOD_NOT_ALLOWED",
		Msg:        "请求方法不被允许",
	}

	ParamsError = BizError{
		HTTPStatus: fiber.StatusBadRequest,
		Code:       400,
		BizErr:     "PARAMS_ERROR",
		Msg:        "参数错误",
	}

	NotLogin = BizError{
		HTTPStatus: fiber.StatusUnauthorized,
		Code:       401,
		BizErr:     "NOT_LOGIN",
		Msg:        "未登录或登录已过期",
	}

	Unauthorized = BizError{
		HTTPStatus: fiber.StatusForbidden,
		Code:       403,
		BizErr:     "UNAUTHORIZED",
		Msg:        "你没有访问该资源的权限",
	}

	InvalidCredential = BizError{
		HTTPStatus: fiber.StatusUnauthorized,
		Code:       40101,
		BizErr:     "INVALID_CREDENTIAL",
		Msg:        "用户名或密码不正确",
	}

	ServerError = BizError{
		HTTPStatus: fiber.StatusInternalServerError,
		Code:       500,
		BizErr:     "SERVER_ERROR",
		Msg:        "服务器内部错误",
	}
)
