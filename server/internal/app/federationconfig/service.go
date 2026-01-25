package federationconfig

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"strconv"
	"strings"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/sysconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	domainconfig "github.com/grtsinry43/grtblog-v2/server/internal/domain/config"
)

type Service struct {
	core *sysconfig.Service
	repo domainconfig.SysConfigRepository
}

func NewService(repo domainconfig.SysConfigRepository) *Service {
	return &Service{
		core: sysconfig.NewService(repo, config.TurnstileConfig{}),
		repo: repo,
	}
}

func (s *Service) ListConfigs(ctx context.Context, keys []string) ([]domainconfig.SysConfig, error) {
	return s.core.ListConfigs(ctx, keys)
}

func (s *Service) UpdateConfigs(ctx context.Context, items []sysconfig.UpdateItem) ([]domainconfig.SysConfig, error) {
	updated, err := s.core.UpdateConfigs(ctx, items)
	if err != nil {
		return nil, err
	}

	settings, err := s.Settings(ctx)
	if err != nil {
		return updated, err
	}
	if !settings.Enabled {
		return updated, nil
	}

	keyUpdates, err := s.ensureKeyPairUpdates(settings)
	if err != nil {
		return updated, err
	}
	if len(keyUpdates) == 0 {
		return updated, nil
	}
	keyUpdated, err := s.core.UpdateConfigs(ctx, keyUpdates)
	if err != nil {
		return updated, err
	}
	updated = append(updated, keyUpdated...)
	return updated, nil
}

// Settings aggregates federation configuration values from kv entries.
type Settings struct {
	Enabled         bool
	InstanceName    string
	InstanceURL     string
	PublicKey       string
	PrivateKey      string
	SignatureAlg    string
	RequireHTTPS    bool
	AllowInbound    bool
	AllowOutbound   bool
	DefaultPolicies json.RawMessage
	RateLimits      json.RawMessage
}

func (s *Service) Settings(ctx context.Context) (Settings, error) {
	keys := []string{
		"federation.enabled",
		"federation.instanceName",
		"federation.instanceURL",
		"federation.publicKey",
		"federation.privateKey",
		"federation.signatureAlg",
		"federation.requireHTTPS",
		"federation.allowInbound",
		"federation.allowOutbound",
		"federation.defaultPolicies",
		"federation.rateLimits",
	}
	items, err := s.repo.List(ctx, keys)
	if err != nil {
		return Settings{}, err
	}
	lookup := make(map[string]domainconfig.SysConfig, len(items))
	for _, item := range items {
		lookup[item.Key] = item
	}

	return Settings{
		Enabled:         parseBool(lookup["federation.enabled"], false),
		InstanceName:    parseString(lookup["federation.instanceName"], ""),
		InstanceURL:     parseString(lookup["federation.instanceURL"], ""),
		PublicKey:       parseString(lookup["federation.publicKey"], ""),
		PrivateKey:      parseString(lookup["federation.privateKey"], ""),
		SignatureAlg:    parseString(lookup["federation.signatureAlg"], "rsa-sha256"),
		RequireHTTPS:    parseBool(lookup["federation.requireHTTPS"], true),
		AllowInbound:    parseBool(lookup["federation.allowInbound"], true),
		AllowOutbound:   parseBool(lookup["federation.allowOutbound"], true),
		DefaultPolicies: parseJSON(lookup["federation.defaultPolicies"], json.RawMessage("{}")),
		RateLimits:      parseJSON(lookup["federation.rateLimits"], json.RawMessage("{}")),
	}, nil
}

func parseString(cfg domainconfig.SysConfig, fallback string) string {
	val := valueOrDefault(cfg)
	if strings.TrimSpace(val) == "" {
		return fallback
	}
	return val
}

func parseBool(cfg domainconfig.SysConfig, fallback bool) bool {
	val := valueOrDefault(cfg)
	if strings.TrimSpace(val) == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseJSON(cfg domainconfig.SysConfig, fallback json.RawMessage) json.RawMessage {
	val := valueOrDefault(cfg)
	if strings.TrimSpace(val) == "" {
		return fallback
	}
	return json.RawMessage(val)
}

func valueOrDefault(cfg domainconfig.SysConfig) string {
	if strings.TrimSpace(cfg.Value) != "" {
		return cfg.Value
	}
	if cfg.DefaultValue != nil {
		return *cfg.DefaultValue
	}
	return ""
}

func (s *Service) ensureKeyPairUpdates(settings Settings) ([]sysconfig.UpdateItem, error) {
	if strings.TrimSpace(settings.PublicKey) != "" && strings.TrimSpace(settings.PrivateKey) != "" {
		return nil, nil
	}

	if settings.SignatureAlg != "" && settings.SignatureAlg != "rsa-sha256" {
		return nil, errors.New("仅支持 rsa-sha256")
	}

	pub, priv, err := generateRSAKeyPair()
	if err != nil {
		return nil, err
	}

	pubRaw, err := json.Marshal(pub)
	if err != nil {
		return nil, err
	}
	privRaw, err := json.Marshal(priv)
	if err != nil {
		return nil, err
	}

	return []sysconfig.UpdateItem{
		{
			Key:   "federation.publicKey",
			Value: toRaw(pubRaw),
		},
		{
			Key:   "federation.privateKey",
			Value: toRaw(privRaw),
		},
	}, nil
}

func generateRSAKeyPair() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}
	privateBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privatePEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateBytes})

	publicBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	publicPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicBytes})

	return string(publicPEM), string(privatePEM), nil
}

func toRaw(raw []byte) *json.RawMessage {
	msg := json.RawMessage(raw)
	return &msg
}
