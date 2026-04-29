//go:build http_requests_extended_methods

package mcpserver

import "testing"

func TestMethodToolSpecsExtendedBuild(t *testing.T) {
	specs := methodToolSpecs()
	if len(specs) != 5 {
		t.Fatalf("unexpected specs count in extended build: %d", len(specs))
	}
	if specs[3].Name != "delete_http_request" {
		t.Fatalf("fourth tool should be delete_http_request, got=%s", specs[3].Name)
	}
	if specs[4].Name != "patch_http_request" {
		t.Fatalf("fifth tool should be patch_http_request, got=%s", specs[4].Name)
	}
}
