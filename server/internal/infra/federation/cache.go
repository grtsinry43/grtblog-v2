package federation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultManifestTTL  = time.Hour
	DefaultPublicKeyTTL = 24 * time.Hour
	DefaultEndpointsTTL = time.Hour
)

// Cache abstracts remote metadata caching.
type Cache interface {
	GetManifest(ctx context.Context, baseURL string) (*Manifest, error)
	SetManifest(ctx context.Context, baseURL string, manifest Manifest, ttl time.Duration) error
	GetPublicKey(ctx context.Context, baseURL string) (*PublicKeyDoc, error)
	SetPublicKey(ctx context.Context, baseURL string, doc PublicKeyDoc, ttl time.Duration) error
	GetEndpoints(ctx context.Context, baseURL string) (*EndpointsDoc, error)
	SetEndpoints(ctx context.Context, baseURL string, doc EndpointsDoc, ttl time.Duration) error
}

// RedisCache stores federation metadata in Redis.
type RedisCache struct {
	client *redis.Client
	prefix string
}

func NewRedisCache(client *redis.Client, prefix string) *RedisCache {
	return &RedisCache{client: client, prefix: prefix}
}

func (c *RedisCache) GetManifest(ctx context.Context, baseURL string) (*Manifest, error) {
	val, err := c.client.Get(ctx, c.key("manifest", baseURL)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var manifest Manifest
	if err := json.Unmarshal([]byte(val), &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func (c *RedisCache) SetManifest(ctx context.Context, baseURL string, manifest Manifest, ttl time.Duration) error {
	payload, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, c.key("manifest", baseURL), payload, ttl).Err()
}

func (c *RedisCache) GetPublicKey(ctx context.Context, baseURL string) (*PublicKeyDoc, error) {
	val, err := c.client.Get(ctx, c.key("pubkey", baseURL)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var doc PublicKeyDoc
	if err := json.Unmarshal([]byte(val), &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (c *RedisCache) SetPublicKey(ctx context.Context, baseURL string, doc PublicKeyDoc, ttl time.Duration) error {
	payload, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, c.key("pubkey", baseURL), payload, ttl).Err()
}

func (c *RedisCache) GetEndpoints(ctx context.Context, baseURL string) (*EndpointsDoc, error) {
	val, err := c.client.Get(ctx, c.key("endpoints", baseURL)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var doc EndpointsDoc
	if err := json.Unmarshal([]byte(val), &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (c *RedisCache) SetEndpoints(ctx context.Context, baseURL string, doc EndpointsDoc, ttl time.Duration) error {
	payload, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, c.key("endpoints", baseURL), payload, ttl).Err()
}

func (c *RedisCache) key(kind string, baseURL string) string {
	escaped := url.PathEscape(baseURL)
	return fmt.Sprintf("%sbfp:%s:%s", c.prefix, kind, escaped)
}
