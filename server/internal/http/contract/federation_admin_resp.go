package contract

// FederationAdminProxyResp 返回远端响应。
type FederationAdminProxyResp struct {
	StatusCode int    `json:"status_code"`
	Body       string `json:"body"`
}

// FederationAdminRemoteCheckResp 返回远端 well-known 信息（仅用于文档与测试展示）。
type FederationAdminRemoteCheckResp struct {
	Manifest  any `json:"manifest,omitempty" swaggertype:"object"`
	PublicKey any `json:"public_key,omitempty" swaggertype:"object"`
	Endpoints any `json:"endpoints,omitempty" swaggertype:"object"`
}
