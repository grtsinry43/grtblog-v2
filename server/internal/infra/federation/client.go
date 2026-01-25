package federation

import (
	"bytes"
	"context"
	"crypto"
	"fmt"
	"net/http"
	"time"
)

// Client sends signed federation requests.
type Client struct {
	httpClient *http.Client
	signer     *Signer
}

func NewClient(httpClient *http.Client, signer *Signer) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &Client{httpClient: httpClient, signer: signer}
}

// DoSigned sends a signed HTTP request with optional JSON body.
func (c *Client) DoSigned(ctx context.Context, method string, url string, body []byte, keyID string, privateKey crypto.PrivateKey) (*http.Response, error) {
	if c.signer == nil {
		return nil, fmt.Errorf("signer not configured")
	}
	var reader *bytes.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	} else {
		reader = bytes.NewReader([]byte{})
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return nil, err
	}
	if err := c.signer.SignRequest(req, body, keyID, privateKey); err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}
