package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"http-requests/internal/httpclient"
	"http-requests/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterTools(s *server.MCPServer) {
	client := httpclient.New()

	tool := mcp.NewTool(
		"send_http_request",
		mcp.WithDescription("Send an HTTP request with optional authentication"),
		mcp.WithString("url", mcp.Description("Request URL"), mcp.Required()),
		mcp.WithString("method", mcp.Description("HTTP method, e.g. GET/POST/PUT/PATCH/DELETE"), mcp.Required()),
		mcp.WithObject("headers", mcp.Description("HTTP headers as key-value pairs")),
		mcp.WithObject("query", mcp.Description("Query parameters as key-value pairs")),
		mcp.WithAny("body", mcp.Description("Request body, supports string or JSON object")),
		mcp.WithNumber("timeout_ms", mcp.Description("Request timeout in milliseconds"), mcp.Min(1)),
		mcp.WithObject(
			"auth",
			mcp.Description("Authentication config"),
			mcp.Properties(map[string]any{
				"type": map[string]any{
					"type":        "string",
					"description": "Auth type",
					"enum":        []string{"bearer", "api_key", "aksk_hmac"},
				},
				"bearer": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"token": map[string]any{"type": "string"},
					},
				},
				"api_key": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"key":   map[string]any{"type": "string"},
						"value": map[string]any{"type": "string"},
						"in":    map[string]any{"type": "string", "enum": []string{"header", "query"}},
					},
				},
				"aksk_hmac": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"access_key":        map[string]any{"type": "string"},
						"secret_key":        map[string]any{"type": "string"},
						"timestamp":         map[string]any{"type": "string"},
						"access_key_header": map[string]any{"type": "string"},
						"signature_header":  map[string]any{"type": "string"},
						"timestamp_header":  map[string]any{"type": "string"},
					},
				},
			}),
		),
	)

	handler := mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args types.SendHTTPRequestInput) (*mcp.CallToolResult, error) {
		if args.TimeoutMS <= 0 {
			args.TimeoutMS = 30000
		}
		resp, err := client.Send(ctx, args)
		if err != nil {
			return mcp.NewToolResultErrorf("send_http_request failed: %v", maskSensitive(err.Error())), nil
		}
		result, err := mcp.NewToolResultJSON(resp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	s.AddTool(tool, handler)
}

func maskSensitive(msg string) string {
	type token struct {
		Key string `json:"key"`
	}
	_ = token{}
	return msg
}

func ParseAuth(raw map[string]any) (*types.AuthConfig, error) {
	if raw == nil {
		return nil, nil
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal auth config failed: %w", err)
	}
	var cfg types.AuthConfig
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("parse auth config failed: %w", err)
	}
	cfg.RawPayload = b
	return &cfg, nil
}
