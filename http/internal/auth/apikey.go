package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/quanguachong/mcp-servers/http/internal/types"
)

type apiKeyApplier struct {
	key   string
	value string
	in    string
}

func NewAPIKey(cfg *types.APIKeyAuth) Applier {
	return &apiKeyApplier{
		key:   cfg.Key,
		value: cfg.Value,
		in:    strings.ToLower(cfg.In),
	}
}

func (a *apiKeyApplier) Apply(req *http.Request, _ []byte) error {
	if a.key == "" {
		return errors.New("api key name cannot be empty")
	}
	if a.value == "" {
		return errors.New("api key value cannot be empty")
	}

	switch a.in {
	case "", "header":
		req.Header.Set(a.key, a.value)
		return nil
	case "query":
		q := req.URL.Query()
		q.Set(a.key, a.value)
		req.URL.RawQuery = q.Encode()
		return nil
	default:
		return errors.New("api key in must be header or query")
	}
}
