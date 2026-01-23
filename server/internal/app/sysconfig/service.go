package sysconfig

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	domainconfig "github.com/grtsinry43/grtblog-v2/server/internal/domain/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/turnstile"
)

// Service 负责从数据库读取系统配置并做类型转换。
type Service struct {
	repo             domainconfig.SysConfigRepository
	defaultTurnstile config.TurnstileConfig
}

func NewService(repo domainconfig.SysConfigRepository, defaults config.TurnstileConfig) *Service {
	return &Service{
		repo:             repo,
		defaultTurnstile: defaults,
	}
}

// Turnstile 返回实时的 Turnstile 配置，优先读取 sys_config，未配置时回退到 env 默认值。
// 约定 key：
// - turnstile.enabled: bool 字符串
// - turnstile.secret: Turnstile Secret
// - turnstile.siteKey: Turnstile Site Key（给前端）
// - turnstile.verifyURL: 覆盖校验端点
// - turnstile.timeoutSeconds: 请求超时秒数
func (s *Service) Turnstile(ctx context.Context) (turnstile.Settings, error) {
	settings := turnstile.Settings{
		Enabled:   s.defaultTurnstile.Enabled,
		Secret:    strings.TrimSpace(s.defaultTurnstile.Secret),
		SiteKey:   "",
		VerifyURL: strings.TrimSpace(s.defaultTurnstile.VerifyURL),
		Timeout:   s.defaultTurnstile.Timeout,
	}

	applyString := func(key string, apply func(string) error) error {
		cfg, err := s.repo.GetByKey(ctx, key)
		if err != nil {
			if err == domainconfig.ErrSysConfigNotFound {
				return nil
			}
			return fmt.Errorf("load %s: %w", key, err)
		}
		val := strings.TrimSpace(cfg.Value)
		if val == "" {
			return nil
		}
		return apply(val)
	}

	if err := applyString("turnstile.enabled", func(val string) error {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("parse bool: %w", err)
		}
		settings.Enabled = b
		return nil
	}); err != nil {
		return settings, err
	}

	_ = applyString("turnstile.secret", func(val string) error {
		settings.Secret = val
		return nil
	})
	_ = applyString("turnstile.siteKey", func(val string) error {
		settings.SiteKey = val
		return nil
	})
	_ = applyString("turnstile.verifyURL", func(val string) error {
		settings.VerifyURL = val
		return nil
	})
	if err := applyString("turnstile.timeoutSeconds", func(val string) error {
		sec, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("parse timeoutSeconds: %w", err)
		}
		if sec > 0 {
			settings.Timeout = time.Duration(sec) * time.Second
		}
		return nil
	}); err != nil {
		return settings, err
	}

	// 如果开启但未配置 Secret，视为关闭以避免空 Secret 造成误判。
	if settings.Enabled && strings.TrimSpace(settings.Secret) == "" {
		settings.Enabled = false
	}
	return settings, nil
}

// UploadMaxSizeBytes 返回上传文件的最大大小（字节），范围 1MB~50MB，默认 50MB。
func (s *Service) UploadMaxSizeBytes(ctx context.Context) int {
	const (
		uploadKey     = "upload.maxSizeMB"
		defaultSizeMB = 50
		minSizeMB     = 1
		maxSizeMB     = 50
	)

	sizeMB := defaultSizeMB
	cfg, err := s.repo.GetByKey(ctx, uploadKey)
	if err == nil {
		val := strings.TrimSpace(cfg.Value)
		if parsed, parseErr := strconv.Atoi(val); parseErr == nil {
			sizeMB = parsed
		}
	}

	if sizeMB < minSizeMB {
		sizeMB = minSizeMB
	}
	if sizeMB > maxSizeMB {
		sizeMB = maxSizeMB
	}
	return sizeMB * 1024 * 1024
}

type WebhookSettings struct {
	Timeout   time.Duration
	Workers   int
	QueueSize int
}

// WebhookSettings 返回 Webhook 发送配置，优先读取 sys_config，未配置时回退默认值。
// 约定 key：
// - webhook.timeoutSeconds: 请求超时秒数
// - webhook.workers: 并发 worker 数
// - webhook.queueSize: 队列长度
func (s *Service) WebhookSettings(ctx context.Context) (WebhookSettings, error) {
	const (
		timeoutKey  = "webhook.timeoutSeconds"
		workersKey  = "webhook.workers"
		queueKey    = "webhook.queueSize"
		defaultSec  = 30
		defaultWork = 4
		defaultQ    = 200
	)

	settings := WebhookSettings{
		Timeout:   time.Duration(defaultSec) * time.Second,
		Workers:   defaultWork,
		QueueSize: defaultQ,
	}

	applyInt := func(key string, apply func(int) error) error {
		cfg, err := s.repo.GetByKey(ctx, key)
		if err != nil {
			if err == domainconfig.ErrSysConfigNotFound {
				return nil
			}
			return fmt.Errorf("load %s: %w", key, err)
		}
		val := strings.TrimSpace(cfg.Value)
		if val == "" {
			return nil
		}
		parsed, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("parse %s: %w", key, err)
		}
		return apply(parsed)
	}

	if err := applyInt(timeoutKey, func(val int) error {
		if val > 0 {
			settings.Timeout = time.Duration(val) * time.Second
		}
		return nil
	}); err != nil {
		return settings, err
	}
	if err := applyInt(workersKey, func(val int) error {
		if val > 0 {
			settings.Workers = val
		}
		return nil
	}); err != nil {
		return settings, err
	}
	if err := applyInt(queueKey, func(val int) error {
		if val > 0 {
			settings.QueueSize = val
		}
		return nil
	}); err != nil {
		return settings, err
	}

	return settings, nil
}
