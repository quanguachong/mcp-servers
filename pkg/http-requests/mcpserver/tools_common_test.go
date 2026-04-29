package mcpserver

import (
	"testing"

	"github.com/quanguachong/mcp-servers/pkg/http-requests/types"
)

func TestBuildSendInput(t *testing.T) {
	in := types.HTTPRequestToolInput{
		URL: "https://example.com",
	}
	got := buildSendInput(in, "post")
	if got.Method != "POST" {
		t.Fatalf("method mismatch, got=%s", got.Method)
	}
	if got.TimeoutMS != 30000 {
		t.Fatalf("timeout mismatch, got=%d", got.TimeoutMS)
	}
}
