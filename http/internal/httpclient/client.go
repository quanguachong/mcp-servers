package httpclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/quanguachong/mcp-servers/http/internal/auth"
	"github.com/quanguachong/mcp-servers/http/internal/types"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) Send(ctx context.Context, input types.SendHTTPRequestInput) (*types.SendHTTPRequestResult, error) {
	if strings.TrimSpace(input.URL) == "" {
		return nil, errors.New("url is required")
	}
	if strings.TrimSpace(input.Method) == "" {
		return nil, errors.New("method is required")
	}

	parsedURL, err := url.Parse(input.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	query := parsedURL.Query()
	for k, v := range input.Query {
		query.Set(k, v)
	}
	parsedURL.RawQuery = query.Encode()

	bodyBytes, err := buildBodyBytes(input.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(input.Method), parsedURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	for k, v := range input.Headers {
		req.Header.Set(k, v)
	}

	applier, err := auth.NewApplier(input.Auth)
	if err != nil {
		return nil, err
	}
	if applier != nil {
		if err := applier.Apply(req, bodyBytes); err != nil {
			return nil, err
		}
	}

	timeout := 30 * time.Second
	if input.TimeoutMS > 0 {
		timeout = time.Duration(input.TimeoutMS) * time.Millisecond
	}
	httpClient := &http.Client{Timeout: timeout}

	start := time.Now()
	resp, err := httpClient.Do(req)
	latency := time.Since(start)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	body, isBase64 := normalizeBody(respBody)

	return &types.SendHTTPRequestResult{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
		BodyBase64: isBase64,
		LatencyMS:  latency.Milliseconds(),
		FinalURL:   resp.Request.URL.String(),
	}, nil
}

func buildBodyBytes(body any) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	switch v := body.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		return b, nil
	}
}

func normalizeBody(b []byte) (string, bool) {
	if len(b) == 0 {
		return "", false
	}
	if utf8.Valid(b) {
		return string(b), false
	}
	return base64.StdEncoding.EncodeToString(b), true
}
