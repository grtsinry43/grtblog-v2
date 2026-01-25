package federation

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/federationconfig"
	domainfed "github.com/grtsinry43/grtblog-v2/server/internal/domain/federation"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	fedinfra "github.com/grtsinry43/grtblog-v2/server/internal/infra/federation"
)

type OutboundService struct {
	cfgSvc       *federationconfig.Service
	resolver     *fedinfra.Resolver
	instanceRepo domainfed.FederationInstanceRepository
	client       *http.Client
}

func NewOutboundService(cfgSvc *federationconfig.Service, resolver *fedinfra.Resolver, instanceRepo domainfed.FederationInstanceRepository) *OutboundService {
	return &OutboundService{
		cfgSvc:       cfgSvc,
		resolver:     resolver,
		instanceRepo: instanceRepo,
		client:       &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *OutboundService) SendFriendLinkRequest(ctx context.Context, target string, message string, rssURL string) (*http.Response, []byte, error) {
	endpoint, err := s.resolveEndpoint(ctx, target, "friendlink_request", "/api/federation/friendlinks/request")
	if err != nil {
		return nil, nil, err
	}
	settings, keyID, privKey, client, err := s.signedClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	payload := contract.FederationFriendLinkRequestReq{
		RequesterURL: settings.InstanceURL,
		Message:      message,
		RSSURL:       rssURL,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}
	resp, err := client.DoSigned(ctx, http.MethodPost, endpoint, body, keyID, privKey)
	if err != nil {
		log.Printf("[federation] 出站 友链申请 target=%s endpoint=%s err=%v", target, endpoint, err)
		return nil, nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	log.Printf("[federation] 出站 友链申请 target=%s endpoint=%s status=%d", target, endpoint, resp.StatusCode)
	return resp, raw, nil
}

func (s *OutboundService) SendCitation(ctx context.Context, ev CitationDetected) (*http.Response, []byte, error) {
	endpoint, err := s.resolveEndpoint(ctx, ev.TargetInstance, "citation_request", "/api/federation/citations/request")
	if err != nil {
		return nil, nil, err
	}
	settings, keyID, privKey, client, err := s.signedClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	sourceURL := strings.TrimRight(settings.InstanceURL, "/") + "/posts/" + ev.ShortURL
	payload := contract.FederationCitationRequestReq{
		SourceInstanceURL: settings.InstanceURL,
		SourcePost: contract.FederationCitationSourcePost{
			ID:    ev.ShortURL,
			URL:   sourceURL,
			Title: ev.Title,
		},
		TargetPostID:    ev.TargetPostID,
		CitationContext: ev.Context,
		CitationType:    firstNonEmpty(ev.CitationType, "reference"),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}
	resp, err := client.DoSigned(ctx, http.MethodPost, endpoint, body, keyID, privKey)
	if err != nil {
		log.Printf("[federation] 出站 引用申请 target=%s endpoint=%s err=%v", ev.TargetInstance, endpoint, err)
		return nil, nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	log.Printf("[federation] 出站 引用申请 target=%s endpoint=%s status=%d", ev.TargetInstance, endpoint, resp.StatusCode)
	return resp, raw, nil
}

func (s *OutboundService) SendMention(ctx context.Context, ev MentionDetected) (*http.Response, []byte, error) {
	endpoint, err := s.resolveEndpoint(ctx, ev.TargetInstance, "mention_notify", "/api/federation/mentions/notify")
	if err != nil {
		return nil, nil, err
	}
	settings, keyID, privKey, client, err := s.signedClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	sourceURL := strings.TrimRight(settings.InstanceURL, "/") + "/posts/" + ev.ShortURL
	payload := contract.FederationMentionNotifyReq{
		SourceInstanceURL: settings.InstanceURL,
		SourcePost: contract.FederationMentionSourcePost{
			URL:   sourceURL,
			Title: ev.Title,
		},
		MentionedUser:  ev.TargetUser,
		MentionContext: ev.Context,
		MentionType:    firstNonEmpty(ev.MentionType, "discussion"),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}
	resp, err := client.DoSigned(ctx, http.MethodPost, endpoint, body, keyID, privKey)
	if err != nil {
		log.Printf("[federation] 出站 提及通知 target=%s endpoint=%s err=%v", ev.TargetInstance, endpoint, err)
		return nil, nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	log.Printf("[federation] 出站 提及通知 target=%s endpoint=%s status=%d", ev.TargetInstance, endpoint, resp.StatusCode)
	return resp, raw, nil
}

func (s *OutboundService) resolveEndpoint(ctx context.Context, target string, key string, fallbackPath string) (string, error) {
	baseURL := s.resolveTargetBaseURL(ctx, target)
	if baseURL == "" {
		return "", errors.New("target instance is empty")
	}
	if s.resolver == nil {
		return "", errors.New("resolver not configured")
	}
	endpoints, err := s.resolver.FetchEndpoints(ctx, baseURL)
	if err != nil {
		return "", err
	}
	if endpoints == nil {
		return "", errors.New("endpoints not found")
	}
	path := endpoints.Endpoints[key]
	if path == "" {
		path = fallbackPath
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path, nil
	}
	base := strings.TrimSpace(endpoints.BaseURL)
	if base == "" {
		return "", errors.New("endpoints.base_url is empty")
	}
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		return "", fmt.Errorf("endpoints.base_url must include scheme: %s", base)
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return strings.TrimRight(base, "/") + path, nil
}

type signingSettings struct {
	InstanceURL   string
	SignatureAlg  string
	PrivateKey    string
	AllowOutbound bool
	Enabled       bool
}

func (s *OutboundService) signedClient(ctx context.Context) (signingSettings, string, crypto.PrivateKey, *fedinfra.Client, error) {
	settings, keyID, privKey, err := s.signingContext(ctx)
	if err != nil {
		return signingSettings{}, "", nil, nil, err
	}
	signer, err := fedinfra.NewSigner(settings.SignatureAlg)
	if err != nil {
		return signingSettings{}, "", nil, nil, err
	}
	client := fedinfra.NewClient(s.client, signer)
	return settings, keyID, privKey, client, nil
}

func (s *OutboundService) signingContext(ctx context.Context) (signingSettings, string, crypto.PrivateKey, error) {
	if s.cfgSvc == nil {
		return signingSettings{}, "", nil, errors.New("config service not configured")
	}
	settings, err := s.cfgSvc.Settings(ctx)
	if err != nil {
		return signingSettings{}, "", nil, err
	}
	if !settings.Enabled || !settings.AllowOutbound {
		return signingSettings{}, "", nil, errors.New("federation outbound disabled")
	}
	if strings.TrimSpace(settings.InstanceURL) == "" {
		return signingSettings{}, "", nil, errors.New("instanceURL not configured")
	}
	if strings.TrimSpace(settings.PrivateKey) == "" {
		return signingSettings{}, "", nil, errors.New("private key not configured")
	}
	privKey, err := parsePrivateKey(settings.PrivateKey)
	if err != nil {
		return signingSettings{}, "", nil, err
	}
	keyID := strings.TrimRight(settings.InstanceURL, "/") + "/.well-known/blog-federation/public-key.json"
	return signingSettings{
		InstanceURL:   strings.TrimRight(settings.InstanceURL, "/"),
		SignatureAlg:  settings.SignatureAlg,
		PrivateKey:    settings.PrivateKey,
		AllowOutbound: settings.AllowOutbound,
		Enabled:       settings.Enabled,
	}, keyID, privKey, nil
}

func parsePrivateKey(pemData string) (crypto.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, errors.New("invalid private key")
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if rsaKey, ok := key.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
		return nil, fmt.Errorf("unsupported private key type")
	}
	return nil, errors.New("unsupported private key format")
}

func (s *OutboundService) resolveTargetBaseURL(ctx context.Context, raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return strings.TrimRight(trimmed, "/")
	}
	host, port := parseHostPort(trimmed)
	if host != "" && s.instanceRepo != nil {
		if instances, err := s.instanceRepo.ListActive(ctx); err == nil {
			for _, instance := range instances {
				base := strings.TrimRight(instance.BaseURL, "/")
				parsed, err := url.Parse(base)
				if err != nil {
					continue
				}
				if !strings.EqualFold(parsed.Hostname(), host) {
					continue
				}
				if port != "" && parsed.Port() != port {
					continue
				}
				return base
			}
		}
	}
	return "https://" + strings.TrimRight(trimmed, "/")
}

func parseHostPort(raw string) (string, string) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", ""
	}
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		parsed, err := url.Parse(trimmed)
		if err != nil {
			return "", ""
		}
		return parsed.Hostname(), parsed.Port()
	}
	parsed, err := url.Parse("http://" + trimmed)
	if err == nil {
		return parsed.Hostname(), parsed.Port()
	}
	hostPart := trimmed
	if idx := strings.Index(hostPart, "/"); idx >= 0 {
		hostPart = hostPart[:idx]
	}
	if host, port, err := net.SplitHostPort(hostPart); err == nil {
		return host, port
	}
	return hostPart, ""
}

func firstNonEmpty(values ...string) string {
	for _, val := range values {
		if strings.TrimSpace(val) != "" {
			return val
		}
	}
	return ""
}
