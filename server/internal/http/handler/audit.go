package handler

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
)

// Audit 记录简单的审计日志，携带 requestId、用户与动作。
func Audit(c *fiber.Ctx, action string, fields map[string]any) {
	reqID, _ := c.Locals("requestId").(string)
	userID := "-"
	if claims, ok := middleware.GetClaims(c); ok {
		userID = stringFromInt64(claims.UserID)
	}
	log.Printf("[audit] req=%s user=%s action=%s fields=%v", reqID, userID, action, fields)
}

func stringFromInt64(v int64) string {
	return fmt.Sprintf("%d", v)
}
