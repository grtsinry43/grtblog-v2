package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/jwt"
)

const authContextKey = "authUser"

// RequireAuth 校验 Authorization header 并解析 JWT。
func RequireAuth(manager *jwt.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := extractToken(c.Get("Authorization"))
		if token == "" {
			return response.ErrorWithMsg[any](c, response.NotLogin, "用户未登录，请提供有效的 token")
		}

		claims, err := manager.Parse(token)
		if err != nil {
			if errors.Is(err, jwt.ErrExpiredToken) {
				return response.ErrorWithMsg[any](c, response.NotLogin, "登录已过期，请重新获取 token")
			}
			return response.ErrorWithMsg[any](c, response.NotLogin, "token 无效")
		}

		c.Locals(authContextKey, claims)
		return c.Next()
	}
}

// RequireAdmin 要求当前用户是管理员。
func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := getClaims(c)
		if !ok {
			return response.ErrorFromBiz[any](c, response.NotLogin)
		}
		if !claims.IsAdmin {
			return response.ErrorWithMsg[any](c, response.Unauthorized, "需要管理员权限")
		}
		return c.Next()
	}
}

// GetClaims 从上下文中获取 JWT claims。
func GetClaims(c *fiber.Ctx) (*jwt.Claims, bool) {
	return getClaims(c)
}

func getClaims(c *fiber.Ctx) (*jwt.Claims, bool) {
	val := c.Locals(authContextKey)
	if val == nil {
		return nil, false
	}
	claims, ok := val.(*jwt.Claims)
	return claims, ok
}

func extractToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
