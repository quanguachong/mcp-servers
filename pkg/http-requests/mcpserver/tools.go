package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/quanguachong/mcp-servers/pkg/http-requests/httpclient"
	"github.com/quanguachong/mcp-servers/pkg/http-requests/types"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type methodToolSpec struct {
	Name        string
	Method      string
	Description string
}

func RegisterTools(s *server.MCPServer) {
	client := httpclient.New()

	for _, spec := range methodToolSpecs() {
		tool := mcp.NewTool(
			spec.Name,
			mcp.WithDescription(spec.Description),
			mcp.WithString("url", mcp.Description("Request URL"), mcp.Required()),
			mcp.WithObject("headers", mcp.Description("HTTP headers as key-value pairs")),
			mcp.WithObject("query", mcp.Description("Query parameters as key-value pairs")),
			mcp.WithAny("body", mcp.Description("Request body, supports string or JSON object")),
			mcp.WithNumber("timeout_ms", mcp.Description("Request timeout in milliseconds"), mcp.Min(1)),
			authToolProperty(),
		)

		handler := newMethodToolHandler(client, spec.Name, spec.Method)
		s.AddTool(tool, handler)
	}
}

func methodToolSpecs() []methodToolSpec {
	specs := []methodToolSpec{
		{Name: "get_http_request", Method: "GET", Description: "Send an HTTP GET request with optional authentication"},
		{Name: "post_http_request", Method: "POST", Description: "Send an HTTP POST request with optional authentication"},
		{Name: "put_http_request", Method: "PUT", Description: "Send an HTTP PUT request with optional authentication"},
	}
	if enableDeleteTool {
		specs = append(specs, methodToolSpec{
			Name:        "delete_http_request",
			Method:      "DELETE",
			Description: "Send an HTTP DELETE request with optional authentication",
		})
	}
	if enablePatchTool {
		specs = append(specs, methodToolSpec{
			Name:        "patch_http_request",
			Method:      "PATCH",
			Description: "Send an HTTP PATCH request with optional authentication",
		})
	}
	return specs
}

func newMethodToolHandler(client *httpclient.Client, toolName string, method string) server.ToolHandlerFunc {
	return mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args types.HTTPRequestToolInput) (*mcp.CallToolResult, error) {
		sendArgs := buildSendInput(args, method)
		resp, err := client.Send(ctx, sendArgs)
		if err != nil {
			return mcp.NewToolResultErrorf("%s failed: %v", toolName, maskSensitive(err.Error())), nil
		}
		result, err := mcp.NewToolResultJSON(resp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

func buildSendInput(args types.HTTPRequestToolInput, method string) types.SendHTTPRequestInput {
	timeoutMS := args.TimeoutMS
	if timeoutMS <= 0 {
		timeoutMS = 30000
	}
	return types.SendHTTPRequestInput{
		URL:       args.URL,
		Method:    strings.ToUpper(method),
		Headers:   args.Headers,
		Query:     args.Query,
		Body:      args.Body,
		TimeoutMS: timeoutMS,
		Auth:      args.Auth,
	}
}

func authToolProperty() mcp.ToolOption {
	return mcp.WithObject(
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
	)
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
