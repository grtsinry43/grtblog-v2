package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
)

// Audit 记录简单的审计日志，携带 requestId、用户与动作。
func Audit(c *fiber.Ctx, action string, fields map[string]any) {
	reqID, _ := c.Locals("requestId").(string)
	userID := "-"
	roleLabel := "匿名"
	if claims, ok := middleware.GetClaims(c); ok {
		userID = stringFromInt64(claims.UserID)
		if claims.IsAdmin {
			roleLabel = "管理员"
		} else {
			roleLabel = "用户"
		}
	} else if fields != nil {
		if v, ok := fields["userId"]; ok {
			userID = fmt.Sprintf("%v", v)
		}
		if v, ok := fields["isAdmin"].(bool); ok && v {
			roleLabel = "管理员"
		} else if userID != "-" {
			roleLabel = "用户"
		}
	}
	message := buildAuditMessage(action, roleLabel, fields)
	log.Printf("[审计] %s req=%s user=%s fields=%v", message, reqID, userID, fields)
}

func stringFromInt64(v int64) string {
	return fmt.Sprintf("%d", v)
}

func buildAuditMessage(action, role string, fields map[string]any) string {
	switch action {
	case "auth.login":
		if role == "管理员" {
			return "管理员登录"
		}
		return "用户登录"
	case "auth.register":
		return "新用户注册"
	case "auth.update_profile":
		return "更新用户资料"
	case "auth.change_password":
		return "修改用户密码"
	case "article.create":
		return "创建文章" + suffixByTitle(fields)
	case "article.update":
		return "更新文章" + suffixByTitle(fields)
	case "article.delete":
		return "删除文章" + suffixByTitle(fields)
	case "admin.oauth.create":
		return "创建 OAuth 提供方" + suffixByKey(fields)
	case "admin.oauth.update":
		return "更新 OAuth 提供方" + suffixByKey(fields)
	case "admin.oauth.delete":
		return "删除 OAuth 提供方" + suffixByKey(fields)
	case "friend-link.submit":
		return "提交友链申请"
	default:
		return action
	}
}

func suffixByTitle(fields map[string]any) string {
	if fields == nil {
		return ""
	}
	if title, ok := fields["title"].(string); ok && strings.TrimSpace(title) != "" {
		return "：" + title
	}
	return ""
}

func suffixByKey(fields map[string]any) string {
	if fields == nil {
		return ""
	}
	if key, ok := fields["key"].(string); ok && strings.TrimSpace(key) != "" {
		return "：" + key
	}
	return ""
}
