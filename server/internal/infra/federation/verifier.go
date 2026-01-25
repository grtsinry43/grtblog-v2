package federation

import (
	"context"
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"code.superseriousbusiness.org/httpsig"
)

// Verifier validates signed federation requests.
type Verifier struct {
	resolver    *Resolver
	allowedSkew time.Duration
}

func NewVerifier(resolver *Resolver, allowedSkew time.Duration) *Verifier {
	if allowedSkew <= 0 {
		allowedSkew = 5 * time.Minute
	}
	return &Verifier{resolver: resolver, allowedSkew: allowedSkew}
}

// VerifyRequest validates digest, date window, and signature.
func (v *Verifier) VerifyRequest(ctx context.Context, req *http.Request, body []byte) (*VerifiedSignature, error) {
	if req.Header.Get("Signature") == "" {
		return nil, ErrMissingSignatureHeader
	}
	if len(body) > 0 {
		if err := verifyDigest(req.Header.Get("Digest"), body); err != nil {
			return nil, err
		}
	}
	requestTime, err := parseRequestTime(req.Header.Get("Date"))
	if err != nil {
		return nil, err
	}
	if skew := time.Since(requestTime); skew > v.allowedSkew || skew < -v.allowedSkew {
		return nil, ErrSignatureExpired
	}

	verifier, err := httpsig.NewVerifier(req)
	if err != nil {
		return nil, err
	}
	keyID := verifier.KeyId()
	if keyID == "" {
		return nil, ErrMissingSignatureHeader
	}

	baseURL, err := baseURLFromKeyID(keyID)
	if err != nil {
		return nil, err
	}
	if v.resolver == nil {
		return nil, fmt.Errorf("resolver not configured")
	}
	pubDoc, err := v.resolver.FetchPublicKey(ctx, baseURL)
	if err != nil {
		return nil, err
	}
	if pubDoc == nil || pubDoc.PublicKey == "" {
		return nil, fmt.Errorf("public key not found")
	}

	pubKey, err := parsePublicKey(pubDoc.PublicKey)
	if err != nil {
		return nil, err
	}

	algo := httpsig.RSA_SHA256
	if pubDoc.Algorithm != "" {
		resolved, err := resolveAlgorithm(pubDoc.Algorithm)
		if err != nil {
			return nil, err
		}
		algo = resolved
	}
	if err := verifier.Verify(pubKey, algo); err != nil {
		return nil, err
	}

	return &VerifiedSignature{
		KeyID:    keyID,
		BaseURL:  baseURL,
		DateTime: requestTime,
	}, nil
}

func verifyDigest(digestHeader string, body []byte) error {
	if digestHeader == "" {
		return ErrInvalidDigest
	}
	parts := strings.SplitN(digestHeader, "=", 2)
	if len(parts) != 2 {
		return ErrInvalidDigest
	}
	algo := strings.TrimSpace(parts[0])
	if !strings.EqualFold(algo, "SHA-256") {
		return ErrInvalidDigest
	}
	expected := strings.TrimSpace(parts[1])
	sum := sha256.Sum256(body)
	actual := base64.StdEncoding.EncodeToString(sum[:])
	if expected != actual {
		return ErrInvalidDigest
	}
	return nil
}

func parseRequestTime(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, fmt.Errorf("missing Date header")
	}
	return time.Parse(http.TimeFormat, raw)
}

func baseURLFromKeyID(keyID string) (string, error) {
	parsed, err := url.Parse(keyID)
	if err != nil {
		return "", err
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("invalid keyId URL")
	}
	return fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host), nil
}

func parsePublicKey(pemData string) (crypto.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, fmt.Errorf("invalid public key PEM")
	}
	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
		return key, nil
	}
	return nil, fmt.Errorf("unsupported public key format")
}

// TODO: support Ed25519 verification once a stable key format is agreed.
