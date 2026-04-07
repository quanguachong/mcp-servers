package types

import "encoding/json"

type AuthConfig struct {
	Type       string          `json:"type"`
	Bearer     *BearerAuth     `json:"bearer,omitempty"`
	APIKey     *APIKeyAuth     `json:"api_key,omitempty"`
	AKSKHMAC   *AKSKHMACAuth   `json:"aksk_hmac,omitempty"`
	RawPayload json.RawMessage `json:"-"`
}

type BearerAuth struct {
	Token string `json:"token"`
}

type APIKeyAuth struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	In    string `json:"in"` // header | query
}

type AKSKHMACAuth struct {
	AccessKey       string `json:"access_key"`
	SecretKey       string `json:"secret_key"`
	Timestamp       string `json:"timestamp,omitempty"`
	AccessKeyHeader string `json:"access_key_header,omitempty"`
	SignatureHeader string `json:"signature_header,omitempty"`
	TimestampHeader string `json:"timestamp_header,omitempty"`
}

type SendHTTPRequestInput struct {
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers,omitempty"`
	Query     map[string]string `json:"query,omitempty"`
	Body      any               `json:"body,omitempty"`
	TimeoutMS int               `json:"timeout_ms,omitempty"`
	Auth      *AuthConfig       `json:"auth,omitempty"`
}

type SendHTTPRequestResult struct {
	StatusCode int                 `json:"status_code"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
	BodyBase64 bool                `json:"body_base64"`
	LatencyMS  int64               `json:"latency_ms"`
	FinalURL   string              `json:"final_url"`
}
