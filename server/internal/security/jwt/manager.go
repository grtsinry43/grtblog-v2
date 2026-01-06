package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
)

var (
	// ErrInvalidToken 表示 token 结构或签名非法。
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken 表示 token 已超出过期时间。
	ErrExpiredToken = errors.New("token expired")
)

// Claims 表示 JWT 的载荷内容。
type Claims struct {
	Subject   string `json:"sub,omitempty"`
	UserID    int64  `json:"uid"`
	IsAdmin   bool   `json:"isAdmin"`
	Issuer    string `json:"iss"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

// Manager 负责签发、解析与校验 JWT。
type Manager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewManager(cfg config.AuthConfig) *Manager {
	return &Manager{
		secret: []byte(cfg.Secret),
		issuer: cfg.Issuer,
		ttl:    cfg.AccessTTL,
	}
}

// Generate 针对指定用户签发 token。
func (m *Manager) Generate(userID int64, isAdmin bool) (string, *Claims, error) {
	now := time.Now()
	claims := &Claims{
		UserID:    userID,
		IsAdmin:   isAdmin,
		Issuer:    m.issuer,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(m.ttl).Unix(),
	}
	token, err := m.sign(claims)
	if err != nil {
		return "", nil, err
	}
	return token, claims, nil
}

// Parse 验证 token 并返回 claims。
func (m *Manager) Parse(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	header := parts[0]
	payload := parts[1]
	signature := parts[2]

	if !m.verifySignature(header, payload, signature) {
		return nil, ErrInvalidToken
	}

	data, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(data, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if err := m.validateClaims(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

func (m *Manager) sign(claims *Claims) (string, error) {
	headerJSON := `{"alg":"HS256","typ":"JWT"}`
	header := base64.RawURLEncoding.EncodeToString([]byte(headerJSON))

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal claims: %w", err)
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	signature := m.computeSignature(header, payload)
	return strings.Join([]string{header, payload, signature}, "."), nil
}

func (m *Manager) computeSignature(header, payload string) string {
	signingInput := header + "." + payload
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(signingInput))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (m *Manager) verifySignature(header, payload, signature string) bool {
	expected := m.computeSignature(header, payload)
	return hmac.Equal([]byte(expected), []byte(signature))
}

func (m *Manager) validateClaims(claims *Claims) error {
	if claims.Issuer != "" && claims.Issuer != m.issuer {
		return ErrInvalidToken
	}
	if claims.ExpiresAt == 0 {
		return ErrInvalidToken
	}
	if time.Now().Unix() > claims.ExpiresAt {
		return ErrExpiredToken
	}
	return nil
}
