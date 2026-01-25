package handler

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/federation"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
	fedinfra "github.com/grtsinry43/grtblog-v2/server/internal/infra/federation"
)

func fetchFederationDocs(ctx context.Context, resolver *fedinfra.Resolver, baseURL string) (*fedinfra.Manifest, *fedinfra.EndpointsDoc, *fedinfra.PublicKeyDoc, error) {
	if resolver == nil {
		return nil, nil, nil, errors.New("resolver not configured")
	}
	baseURL = strings.TrimRight(baseURL, "/")
	manifest, err := resolver.FetchManifest(ctx, baseURL)
	if err != nil {
		return nil, nil, nil, response.NewBizErrorWithCause(response.ServerError, "拉取远程 manifest 失败", err)
	}
	endpoints, err := resolver.FetchEndpoints(ctx, baseURL)
	if err != nil {
		return nil, nil, nil, response.NewBizErrorWithCause(response.ServerError, "拉取远程 endpoints 失败", err)
	}
	publicKey, err := resolver.FetchPublicKey(ctx, baseURL)
	if err != nil {
		return nil, nil, nil, response.NewBizErrorWithCause(response.ServerError, "拉取远程公钥失败", err)
	}
	return manifest, endpoints, publicKey, nil
}

func ensureFederationInstance(ctx context.Context, baseURL string, resolver *fedinfra.Resolver, instanceRepo federation.FederationInstanceRepository) (*federation.FederationInstance, error) {
	manifest, endpoints, publicKey, err := fetchFederationDocs(ctx, resolver, baseURL)
	if err != nil {
		return nil, err
	}
	return ensureInstanceFromDocs(ctx, baseURL, manifest, endpoints, publicKey, instanceRepo)
}

func ensureInstanceFromDocs(
	ctx context.Context,
	baseURL string,
	manifest *fedinfra.Manifest,
	endpoints *fedinfra.EndpointsDoc,
	keyDoc *fedinfra.PublicKeyDoc,
	instanceRepo federation.FederationInstanceRepository,
) (*federation.FederationInstance, error) {
	if instanceRepo == nil {
		return nil, errors.New("instance repository not configured")
	}
	baseURL = strings.TrimRight(baseURL, "/")
	features := toJSON(manifest.Features)
	policies := toJSON(manifest.Policies)
	endpointsPayload := toJSON(endpoints)

	instance, err := instanceRepo.GetByBaseURL(ctx, baseURL)
	if err != nil {
		if !errors.Is(err, federation.ErrFederationInstanceNotFound) {
			return nil, err
		}
		newInstance := &federation.FederationInstance{
			BaseURL:         baseURL,
			Name:            toOptionalString(manifest.Instance.Name),
			Description:     toOptionalString(manifest.Instance.Description),
			ProtocolVersion: toOptionalString(manifest.ProtocolVersion),
			PublicKey:       toOptionalString(keyDoc.PublicKey),
			KeyID:           toOptionalString(keyDoc.KeyID),
			Features:        features,
			Policies:        policies,
			Endpoints:       endpointsPayload,
			Status:          "pending",
			LastSeenAt:      timePtr(time.Now().UTC()),
		}
		if err := instanceRepo.Create(ctx, newInstance); err != nil {
			return nil, err
		}
		return newInstance, nil
	}

	instance.Name = toOptionalString(manifest.Instance.Name)
	instance.Description = toOptionalString(manifest.Instance.Description)
	instance.ProtocolVersion = toOptionalString(manifest.ProtocolVersion)
	instance.PublicKey = toOptionalString(keyDoc.PublicKey)
	instance.KeyID = toOptionalString(keyDoc.KeyID)
	instance.Features = features
	instance.Policies = policies
	instance.Endpoints = endpointsPayload
	instance.LastSeenAt = timePtr(time.Now().UTC())
	if err := instanceRepo.Update(ctx, instance); err != nil {
		return nil, err
	}
	return instance, nil
}

func toOptionalString(val string) *string {
	trimmed := strings.TrimSpace(val)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func toJSON(value any) json.RawMessage {
	payload, err := json.Marshal(value)
	if err != nil {
		return json.RawMessage("{}")
	}
	return payload
}

func timePtr(t time.Time) *time.Time {
	return &t
}
