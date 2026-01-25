package federation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// Resolver fetches well-known federation metadata.
type Resolver struct {
	client *http.Client
	cache  Cache
}

func NewResolver(client *http.Client, cache Cache) *Resolver {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &Resolver{client: client, cache: cache}
}

func (r *Resolver) FetchManifest(ctx context.Context, baseURL string) (*Manifest, error) {
	baseURL = normalizeBaseURL(baseURL)
	if r.cache != nil {
		if cached, err := r.cache.GetManifest(ctx, baseURL); err != nil {
			return nil, err
		} else if cached != nil {
			return cached, nil
		}
	}
	var manifest Manifest
	if err := r.fetchJSON(ctx, baseURL, "manifest.json", &manifest); err != nil {
		return nil, err
	}
	if r.cache != nil {
		_ = r.cache.SetManifest(ctx, baseURL, manifest, DefaultManifestTTL)
	}
	return &manifest, nil
}

func (r *Resolver) FetchPublicKey(ctx context.Context, baseURL string) (*PublicKeyDoc, error) {
	baseURL = normalizeBaseURL(baseURL)
	if r.cache != nil {
		if cached, err := r.cache.GetPublicKey(ctx, baseURL); err != nil {
			return nil, err
		} else if cached != nil {
			return cached, nil
		}
	}
	var doc PublicKeyDoc
	if err := r.fetchJSON(ctx, baseURL, "public-key.json", &doc); err != nil {
		return nil, err
	}
	if r.cache != nil {
		_ = r.cache.SetPublicKey(ctx, baseURL, doc, DefaultPublicKeyTTL)
	}
	return &doc, nil
}

func (r *Resolver) FetchEndpoints(ctx context.Context, baseURL string) (*EndpointsDoc, error) {
	baseURL = normalizeBaseURL(baseURL)
	if r.cache != nil {
		if cached, err := r.cache.GetEndpoints(ctx, baseURL); err != nil {
			return nil, err
		} else if cached != nil {
			return cached, nil
		}
	}
	var doc EndpointsDoc
	if err := r.fetchJSON(ctx, baseURL, "endpoints.json", &doc); err != nil {
		return nil, err
	}
	if r.cache != nil {
		_ = r.cache.SetEndpoints(ctx, baseURL, doc, DefaultEndpointsTTL)
	}
	return &doc, nil
}

func (r *Resolver) fetchJSON(ctx context.Context, baseURL string, filename string, target any) error {
	wellKnownURL, err := buildWellKnownURL(baseURL, filename)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, wellKnownURL, nil)
	if err != nil {
		return err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("well-known request failed: %s", resp.Status)
	}
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(target)
}

func buildWellKnownURL(baseURL string, filename string) (string, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	parsed.Path = path.Join(parsed.Path, "/.well-known/blog-federation/", filename)
	return parsed.String(), nil
}

func normalizeBaseURL(raw string) string {
	return strings.TrimRight(raw, "/")
}
