package httpclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"http-requests/internal/types"
)

func TestSendCustomMethodPassthrough(t *testing.T) {
	var gotMethod string
	var gotBody string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := New()
	_, err := c.Send(context.Background(), types.SendHTTPRequestInput{
		URL:    ts.URL + "/hello",
		Method: "PURGE",
		Body:   `{"ok":true}`,
	})
	if err != nil {
		t.Fatalf("send failed: %v", err)
	}
	if gotMethod != "PURGE" {
		t.Fatalf("method not passthrough, got=%s", gotMethod)
	}
	if gotBody != `{"ok":true}` {
		t.Fatalf("body mismatch, got=%s", gotBody)
	}
}
