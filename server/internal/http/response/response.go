package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// TimeLayout 用于对外时间字段格式化。
const TimeLayout = time.RFC3339

// Meta 附加元信息
type Meta struct {
	RequestID string `json:"requestId"`
	Timestamp string `json:"timestamp"`
}

// Envelope 定义统一响应结构
type Envelope[T any] struct {
	Code   int    `json:"code"`
	BizErr string `json:"bizErr"`
	Msg    string `json:"msg"`
	Data   T      `json:"data"`
	Meta   Meta   `json:"meta"`
}

// 从上下文构造 Meta（可根据项目实际调整）
func metaFromCtx(c *fiber.Ctx) Meta {
	// 优先从 Locals 取（可以在中间件里提前写进去）
	reqID, _ := c.Locals("requestId").(string)
	if reqID == "" {
		// 兜底：从 header 取，或者自己生成
		reqID = c.Get("X-Request-ID", "")
	}

	return Meta{
		RequestID: reqID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// Success 返回 200 + 业务成功码的响应。
func Success[T any](c *fiber.Ctx, data T) error {
	return SuccessWithMessage(c, data, OK.Msg)
}

// SuccessWithMessage 允许自定义成功文案。
func SuccessWithMessage[T any](c *fiber.Ctx, data T, msg string) error {
	return respond(c, OK.HTTPStatus, OK.Code, OK.BizErr, msg, data)
}

// ErrorFromBiz 根据预定义 BizError 返回错误响应（最常用）
func ErrorFromBiz[T any](c *fiber.Ctx, be BizError) error {
	var zero T
	return respond(c, be.HTTPStatus, be.Code, be.BizErr, be.Msg, zero)
}

// ErrorWithMsg 允许覆盖默认错误文案
func ErrorWithMsg[T any](c *fiber.Ctx, be BizError, msg string) error {
	if msg == "" {
		msg = be.Msg
	}
	var zero T
	return respond(c, be.HTTPStatus, be.Code, be.BizErr, msg, zero)
}

// 低层封装：真正写出 JSON 的地方
func respond[T any](c *fiber.Ctx, status int, code int, bizErr string, msg string, data T) error {
	return c.Status(status).JSON(Envelope[T]{
		Code:   code,
		BizErr: bizErr,
		Msg:    msg,
		Data:   data,
		Meta:   metaFromCtx(c),
	})
}
