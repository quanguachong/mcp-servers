package pkg

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/bson"
)

func TestApplyFindDefaults(t *testing.T) {
	limit, bytesLimit := applyFindDefaults(nil, nil)
	if limit != defaultFindLimit {
		t.Fatalf("limit default=%d, want %d", limit, defaultFindLimit)
	}
	if bytesLimit != defaultResponseBytesLimit {
		t.Fatalf("responseBytesLimit default=%d, want %d", bytesLimit, defaultResponseBytesLimit)
	}

	customLimit := int64(0)
	customBytesLimit := int64(2048)
	limit, bytesLimit = applyFindDefaults(&customLimit, &customBytesLimit)
	if limit != 0 {
		t.Fatalf("limit=%d, want 0", limit)
	}
	if bytesLimit != 2048 {
		t.Fatalf("responseBytesLimit=%d, want 2048", bytesLimit)
	}
}

func TestEnsureResponseBytesWithinLimit(t *testing.T) {
	docs := []bson.M{{"name": "alice"}}
	if err := ensureResponseBytesWithinLimit(docs, 1024); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := ensureResponseBytesWithinLimit(docs, 1)
	if err == nil {
		t.Fatal("expected bytes limit error")
	}
}

func TestFindAndListCollectionsInputValidation(t *testing.T) {
	_, _, err := handleFind(context.Background(), &mcp.CallToolRequest{}, findInput{})
	if err == nil {
		t.Fatal("expected error when database/collection missing")
	}

	negative := int64(-1)
	_, _, err = handleFind(context.Background(), &mcp.CallToolRequest{}, findInput{
		Database:   "db",
		Collection: "coll",
		Limit:      &negative,
	})
	if err == nil {
		t.Fatal("expected error when limit < 0")
	}

	zeroBytesLimit := int64(0)
	_, _, err = handleFind(context.Background(), &mcp.CallToolRequest{}, findInput{
		Database:           "db",
		Collection:         "coll",
		ResponseBytesLimit: &zeroBytesLimit,
	})
	if err == nil {
		t.Fatal("expected error when responseBytesLimit <= 0")
	}

	_, _, err = handleListCollections(context.Background(), &mcp.CallToolRequest{}, listCollectionsInput{})
	if err == nil {
		t.Fatal("expected error when database missing")
	}
}
