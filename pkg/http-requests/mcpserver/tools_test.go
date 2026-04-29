//go:build !http_requests_extended_methods

package mcpserver

import "testing"

func TestMethodToolSpecsDefaultBuild(t *testing.T) {
	specs := methodToolSpecs()
	if len(specs) != 3 {
		t.Fatalf("unexpected specs count in default build: %d", len(specs))
	}
	if specs[0].Name != "get_http_request" || specs[1].Name != "post_http_request" || specs[2].Name != "put_http_request" {
		t.Fatalf("unexpected default tools: %#v", specs)
	}
}
