package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"testing"

	"github.com/quanguachong/mcp-servers/pkg/http-requests/types"
)

func TestAPIKeyQueryApply(t *testing.T) {
	u, _ := url.Parse("https://example.com/path")
	req := &http.Request{URL: u, Header: http.Header{}}
	applier := NewAPIKey(&types.APIKeyAuth{
		Key:   "k",
		Value: "v",
		In:    "query",
	})
	if err := applier.Apply(req, nil); err != nil {
		t.Fatalf("apply failed: %v", err)
	}
	if got := req.URL.Query().Get("k"); got != "v" {
		t.Fatalf("unexpected query value: %s", got)
	}
}

func TestAKSKHMACDeterministic(t *testing.T) {
	u, _ := url.Parse("https://example.com/v1/resource?a=1")
	req := &http.Request{
		Method: "POST",
		URL:    u,
		Header: http.Header{},
	}
	body := []byte(`{"x":1}`)
	applier := NewAKSKHMAC(&types.AKSKHMACAuth{
		AccessKey: "ak",
		SecretKey: "sk",
		Timestamp: "2026-01-01T00:00:00Z",
	})
	if err := applier.Apply(req, body); err != nil {
		t.Fatalf("apply failed: %v", err)
	}

	sum := sha256.Sum256(body)
	canonical := "POST\n/v1/resource\na=1\n" + hex.EncodeToString(sum[:]) + "\n2026-01-01T00:00:00Z"
	mac := hmac.New(sha256.New, []byte("sk"))
	_, _ = mac.Write([]byte(canonical))
	expected := hex.EncodeToString(mac.Sum(nil))

	if got := req.Header.Get("X-Signature"); got != expected {
		t.Fatalf("signature mismatch, expected=%s got=%s", expected, got)
	}
}

func TestNewApplierEmptyAuthObject(t *testing.T) {
	applier, err := NewApplier(&types.AuthConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if applier != nil {
		t.Fatalf("expected nil applier for empty auth object")
	}
}

func TestNewApplierInferTypeFromPayload(t *testing.T) {
	applier, err := NewApplier(&types.AuthConfig{
		Bearer: &types.BearerAuth{Token: "abc"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if applier == nil {
		t.Fatalf("expected inferred bearer applier")
	}
}
