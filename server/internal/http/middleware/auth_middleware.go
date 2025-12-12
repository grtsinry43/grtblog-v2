package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/jwt"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/rbac"
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

// RequireRoles 要求当前用户拥有任意一个角色。
func RequireRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := getClaims(c)
		if !ok {
			return response.ErrorFromBiz[any](c, response.NotLogin)
		}
		if len(roles) == 0 {
			return c.Next()
		}
		for _, role := range roles {
			if hasString(claims.Roles, role) {
				return c.Next()
			}
		}
		return response.ErrorWithMsg[any](c, response.Unauthorized, "缺少必要角色")
	}
}

// RequirePermission 借助 Casbin 校验权限字符串。
func RequirePermission(enforcer *rbac.Enforcer, permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if permission == "" || enforcer == nil {
			return c.Next()
		}
		claims, ok := getClaims(c)
		if !ok {
			return response.ErrorFromBiz[any](c, response.NotLogin)
		}
		if !enforcer.HasPermission(claims.Roles, permission) {
			return response.ErrorWithMsg[any](c, response.Unauthorized, "缺少必要权限: "+permission)
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

func hasString(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
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
