package turnstile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
)

var (
	// ErrVerificationFailed 表示 Turnstile 校验未通过。
	ErrVerificationFailed = errors.New("turnstile verification failed")
	// ErrMissingSecret 表示启用 Turnstile 但未提供 Secret。
	ErrMissingSecret = errors.New("turnstile secret not configured")
)

// Settings 描述一次校验需要的配置。
type Settings struct {
	Enabled   bool
	Secret    string
	SiteKey   string
	VerifyURL string
	Timeout   time.Duration
}

// Client 封装 Cloudflare Turnstile 校验。
type Client struct {
	defaultVerifyURL string
	defaultTimeout   time.Duration
	httpClient       *http.Client
}

type verifyResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
	Challenge  string   `json:"challenge_ts"`
	Hostname   string   `json:"hostname"`
	Action     string   `json:"action"`
	CData      string   `json:"cdata"`
}

// NewClient 构造 Turnstile 校验客户端，默认值来自环境配置。
func NewClient(cfg config.TurnstileConfig) *Client {
	httpClient := &http.Client{Timeout: cfg.Timeout}
	return &Client{
		defaultVerifyURL: strings.TrimSpace(cfg.VerifyURL),
		defaultTimeout:   cfg.Timeout,
		httpClient:       httpClient,
	}
}

// Verify 对前端传入的 token 执行 Turnstile 校验（支持运行时动态配置）。
// remoteIP 可选；如果可解析为 IP，将透传给 Cloudflare 提升风控准确度。
func (c *Client) Verify(ctx context.Context, token, remoteIP string, cfg Settings) error {
	if !cfg.Enabled {
		return nil
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return fmt.Errorf("%w: missing token", ErrVerificationFailed)
	}

	secret := strings.TrimSpace(cfg.Secret)
	if secret == "" {
		return ErrMissingSecret
	}

	verifyURL := strings.TrimSpace(cfg.VerifyURL)
	if verifyURL == "" {
		verifyURL = c.defaultVerifyURL
	}

	client := c.httpClient
	if cfg.Timeout > 0 && cfg.Timeout != c.defaultTimeout {
		cp := *c.httpClient
		cp.Timeout = cfg.Timeout
		client = &cp
	}

	form := url.Values{}
	form.Set("secret", secret)
	form.Set("response", token)
	if ip := net.ParseIP(strings.TrimSpace(remoteIP)); ip != nil {
		form.Set("remoteip", ip.String())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, verifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var vr verifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&vr); err != nil {
		return fmt.Errorf("decode turnstile response: %w", err)
	}
	if !vr.Success {
		return fmt.Errorf("%w: %v", ErrVerificationFailed, vr.ErrorCodes)
	}
	return nil
}
