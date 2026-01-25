package federation

import (
	"crypto"
	"net/http"
	"sync"
	"time"

	"code.superseriousbusiness.org/httpsig"
)

// Signer wraps httpsig signer with protocol defaults.
type Signer struct {
	// Note: httpsig signers are not safe for concurrent use.
	signer httpsig.SignerWithOptions
	mu     sync.Mutex
}

// NewSigner constructs a signer for the given algorithm string.
func NewSigner(algorithm string) (*Signer, error) {
	algo, err := resolveAlgorithm(algorithm)
	if err != nil {
		return nil, err
	}

	headers := []string{
		httpsig.RequestTarget,
		"host",
		"date",
		"digest",
		"content-type",
	}

	// expiresIn=0 means no explicit expires parameter in the signature.
	signer, _, err := httpsig.NewSigner([]httpsig.Algorithm{algo}, httpsig.DigestSha256, headers, httpsig.Signature, 0)
	if err != nil {
		return nil, err
	}

	return &Signer{signer: signer}, nil
}

// SignRequest adds digest/date headers and signs the request.
func (s *Signer) SignRequest(req *http.Request, body []byte, keyID string, privateKey crypto.PrivateKey) error {
	if req.Header.Get("Date") == "" {
		req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	}
	if len(body) > 0 && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if req.Host == "" && req.URL != nil {
		req.Host = req.URL.Host
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.signer.SignRequest(privateKey, keyID, req, body)
}

func resolveAlgorithm(algorithm string) (httpsig.Algorithm, error) {
	switch algorithm {
	case "rsa-sha256", "rsa_sha256", "RSA-SHA256":
		return httpsig.RSA_SHA256, nil
	default:
		// TODO: add ed25519 support once a stable signer strategy is defined.
		return httpsig.Algorithm(0), ErrUnsupportedSignatureAlgorithm
	}
}
