package contract

import "github.com/grtsinry43/grtblog-v2/server/internal/http/response"

// OAuthProviderResp 返回可用的 OAuth provider 信息。
type OAuthProviderResp struct {
	Key          string   `json:"key"`
	DisplayName  string   `json:"displayName"`
	Scopes       []string `json:"scopes"`
	PKCERequired bool     `json:"pkceRequired"`
}

// ProviderListRespEnvelope 用于 swagger 展示。
type ProviderListRespEnvelope struct {
	Code   int                 `json:"code"`
	BizErr string              `json:"bizErr"`
	Msg    string              `json:"msg"`
	Data   []OAuthProviderResp `json:"data"`
	Meta   response.Meta       `json:"meta"`
}

// AuthorizeResp 授权信息响应。
type AuthorizeResp struct {
	AuthURL       string `json:"authUrl"`
	State         string `json:"state"`
	CodeChallenge string `json:"codeChallenge,omitempty"`
}

// AuthorizeRespEnvelope 用于 swagger 展示。
type AuthorizeRespEnvelope struct {
	Code   int           `json:"code"`
	BizErr string        `json:"bizErr"`
	Msg    string        `json:"msg"`
	Data   AuthorizeResp `json:"data"`
	Meta   response.Meta `json:"meta"`
}

// OAuthCallbackReq 回调请求。
type OAuthCallbackReq struct {
	Code  string `json:"code"`
	State string `json:"state"`
}
