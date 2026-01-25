package handler

import (
	"fmt"
	"net/url"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/federationconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	fedinfra "github.com/grtsinry43/grtblog-v2/server/internal/infra/federation"
)

type FederationWellKnownHandler struct {
	cfgSvc *federationconfig.Service
	appCfg config.AppConfig
}

func NewFederationWellKnownHandler(cfgSvc *federationconfig.Service, appCfg config.AppConfig) *FederationWellKnownHandler {
	return &FederationWellKnownHandler{cfgSvc: cfgSvc, appCfg: appCfg}
}

func (h *FederationWellKnownHandler) Manifest(c *fiber.Ctx) error {
	settings, err := h.cfgSvc.Settings(c.Context())
	if err != nil || !settings.Enabled {
		return c.SendStatus(fiber.StatusNotFound)
	}
	policy := parseFederationPolicy(settings)
	instanceURL := h.resolveInstanceURL(c, settings)
	features := []string{"friendlink-timeline", "related-posts"}
	if policyBool(policy.AllowCitation, true) {
		features = append(features, "cross-citation")
	}
	if policyBool(policy.AllowMention, true) {
		features = append(features, "cross-mention")
	}
	manifest := fedinfra.Manifest{
		ProtocolVersion: "1.0.0",
		Instance: fedinfra.ManifestNode{
			Name:        h.resolveInstanceName(settings),
			URL:         instanceURL,
			Description: "",
		},
		Software: fedinfra.ManifestSoftware{
			Name:    h.appCfg.Name,
			Version: buildVersion(),
		},
		Features: features,
		Policies: fedinfra.ManifestPolicy{
			AllowCitation:                 policyBool(policy.AllowCitation, true),
			AllowMention:                  policyBool(policy.AllowMention, true),
			AutoApproveFriendlinkCitation: policyBool(policy.AutoApproveFriendlinkCitation, false),
			RequireHTTPS:                  settings.RequireHTTPS,
			MaxCacheAge:                   86400,
		},
		RateLimits: fedinfra.ManifestRate{},
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	return c.JSON(manifest)
}

func (h *FederationWellKnownHandler) PublicKey(c *fiber.Ctx) error {
	settings, err := h.cfgSvc.Settings(c.Context())
	if err != nil || !settings.Enabled || strings.TrimSpace(settings.PublicKey) == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	keyID := h.publicKeyID(c, settings)
	doc := fedinfra.PublicKeyDoc{
		KeyID:     keyID,
		Algorithm: settings.SignatureAlg,
		PublicKey: settings.PublicKey,
		CreatedAt: time.Now().UTC(),
	}
	return c.JSON(doc)
}

func (h *FederationWellKnownHandler) Endpoints(c *fiber.Ctx) error {
	settings, err := h.cfgSvc.Settings(c.Context())
	if err != nil || !settings.Enabled {
		return c.SendStatus(fiber.StatusNotFound)
	}
	baseURL := strings.TrimRight(h.resolveInstanceURL(c, settings), "/") + "/api/federation"
	doc := fedinfra.EndpointsDoc{
		BaseURL: baseURL,
		Endpoints: map[string]string{
			"friendlink_request": "/friendlinks/request",
			"timeline":           "/timeline/posts",
			"post_detail":        "/posts/{id}",
			"citation_request":   "/citations/request",
			"mention_notify":     "/mentions/notify",
		},
	}
	return c.JSON(doc)
}

func (h *FederationWellKnownHandler) resolveInstanceName(settings federationconfig.Settings) string {
	if strings.TrimSpace(settings.InstanceName) != "" {
		return strings.TrimSpace(settings.InstanceName)
	}
	if h.appCfg.Name != "" {
		return h.appCfg.Name
	}
	return "federation-instance"
}

func (h *FederationWellKnownHandler) resolveInstanceURL(c *fiber.Ctx, settings federationconfig.Settings) string {
	if strings.TrimSpace(settings.InstanceURL) != "" {
		return normalizeInstanceURL(settings.InstanceURL)
	}
	scheme := "https"
	if c.Protocol() != "" {
		scheme = c.Protocol()
	}
	host := c.Hostname()
	return fmt.Sprintf("%s://%s", scheme, host)
}

func (h *FederationWellKnownHandler) publicKeyID(c *fiber.Ctx, settings federationconfig.Settings) string {
	baseURL := h.resolveInstanceURL(c, settings)
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}
	parsed.Path = "/.well-known/blog-federation/public-key.json"
	return parsed.String()
}

func buildVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info == nil {
		return "dev"
	}
	if info.Main.Version != "" {
		return info.Main.Version
	}
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			return setting.Value
		}
	}
	return "dev"
}
