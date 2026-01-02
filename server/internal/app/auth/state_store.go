package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// OAuthState 保存 state 与 PKCE 信息，存储在 Redis。
type OAuthState struct {
	Provider     string    `json:"provider"`
	Redirect     string    `json:"redirect,omitempty"`
	CodeVerifier string    `json:"code_verifier,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type StateStore interface {
	Save(ctx context.Context, state string, data OAuthState, ttl time.Duration) error
	Load(ctx context.Context, state string) (*OAuthState, error)
	Delete(ctx context.Context, state string) error
}

type redisStateStore struct {
	client *redis.Client
	prefix string
}

func NewRedisStateStore(client *redis.Client, prefix string) StateStore {
	return &redisStateStore{
		client: client,
		prefix: prefix,
	}
}

func (s *redisStateStore) Save(ctx context.Context, state string, data OAuthState, ttl time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}
	return s.client.Set(ctx, s.key(state), b, ttl).Err()
}

func (s *redisStateStore) Load(ctx context.Context, state string) (*OAuthState, error) {
	val, err := s.client.Get(ctx, s.key(state)).Bytes()
	if err != nil {
		return nil, err
	}
	var data OAuthState
	if err := json.Unmarshal(val, &data); err != nil {
		return nil, fmt.Errorf("unmarshal state: %w", err)
	}
	return &data, nil
}

func (s *redisStateStore) Delete(ctx context.Context, state string) error {
	return s.client.Del(ctx, s.key(state)).Err()
}

func (s *redisStateStore) key(state string) string {
	return s.prefix + "oauth_state:" + state
}

func randomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// GenerateState 返回随机 state。
func GenerateState() (string, error) {
	return randomString(24)
}

// GenerateCodeVerifier 生成 PKCE code verifier。
func GenerateCodeVerifier() (string, string, error) {
	verifier, err := randomString(32)
	if err != nil {
		return "", "", err
	}
	challenge := base64.RawURLEncoding.EncodeToString(hashSHA256([]byte(verifier)))
	return verifier, challenge, nil
}

func hashSHA256(data []byte) []byte {
	h := sha256Sum(data)
	return h[:]
}

func sha256Sum(data []byte) [32]byte {
	return sha256.Sum256(data)
}
