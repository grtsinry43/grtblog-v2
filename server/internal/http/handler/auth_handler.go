package handler

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/auth"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/sysconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/identity"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/turnstile"
)

type AuthHandler struct {
	svc       *auth.Service
	sysCfg    *sysconfig.Service
	turnstile TurnstileVerifier
}

// TurnstileVerifier 便于替换实现/注入 mock。
type TurnstileVerifier interface {
	Verify(ctx context.Context, token, remoteIP string, settings turnstile.Settings) error
}

func NewAuthHandler(svc *auth.Service, sysCfg *sysconfig.Service, verifier TurnstileVerifier) *AuthHandler {
	return &AuthHandler{svc: svc, sysCfg: sysCfg, turnstile: verifier}
}

type RegisterRequest struct {
	Username       string `json:"username"`
	Nickname       string `json:"nickname"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	TurnstileToken string `json:"turnstileToken"`
}

type LoginRequest struct {
	Credential     string `json:"credential"` // username or email
	Password       string `json:"password"`
	TurnstileToken string `json:"turnstileToken"`
}

// Register godoc
// @Summary 用户注册
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册参数"
// @Success 200 {object} RegisterResponseEnvelope
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体解析失败")
	}
	if err := h.verifyTurnstile(c, req.TurnstileToken); err != nil {
		return err
	}
	user, err := h.svc.Register(c.Context(), auth.RegisterCommand{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, identity.ErrUserExists) {
			return response.NewBizErrorWithMsg(response.ParamsError, "用户名或邮箱已存在")
		}
		return err
	}
	return response.SuccessWithMessage(c, toUserResponse(*user), "注册成功")
}

// Login godoc
// @Summary 用户登录
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录参数"
// @Success 200 {object} LoginResponseEnvelope
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体解析失败")
	}
	if err := h.verifyTurnstile(c, req.TurnstileToken); err != nil {
		return err
	}
	result, err := h.svc.Login(c.Context(), auth.LoginCommand{
		Credential: req.Credential,
		Password:   req.Password,
	})
	if err != nil {
		if errors.Is(err, identity.ErrInvalidCredentials) {
			return response.NewBizError(response.InvalidCredential)
		}
		return err
	}
	return response.Success(c, LoginResponse{
		Token: result.Token,
		User:  toUserResponse(result.User),
	})
}

func (h *AuthHandler) verifyTurnstile(c *fiber.Ctx, token string) error {
	if h.turnstile == nil || h.sysCfg == nil {
		return nil
	}
	settings, err := h.sysCfg.Turnstile(c.Context())
	if err != nil {
		return response.NewBizErrorWithMsg(response.ServerError, "获取系统配置失败")
	}
	if !settings.Enabled {
		return nil
	}
	if err := h.turnstile.Verify(c.Context(), token, c.IP(), settings); err != nil {
		if errors.Is(err, turnstile.ErrVerificationFailed) {
			return response.NewBizErrorWithMsg(response.ParamsError, "人机校验未通过")
		}
		if errors.Is(err, turnstile.ErrMissingSecret) {
			return response.NewBizErrorWithMsg(response.ServerError, "人机校验未配置，请联系管理员")
		}
		return response.NewBizErrorWithMsg(response.ServerError, "人机校验失败，请稍后重试")
	}
	return nil
}
