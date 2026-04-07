package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/quanguachong/mcp-servers/http/internal/types"
)

type Applier interface {
	Apply(req *http.Request, body []byte) error
}

func NewApplier(cfg *types.AuthConfig) (Applier, error) {
	if cfg == nil {
		return nil, nil
	}

	authType := strings.TrimSpace(cfg.Type)
	if authType == "" {
		// 兼容 inspector 传入空 auth 对象，或仅传认证子对象不传 type 的场景
		switch {
		case cfg.Bearer != nil:
			authType = "bearer"
		case cfg.APIKey != nil:
			authType = "api_key"
		case cfg.AKSKHMAC != nil:
			authType = "aksk_hmac"
		default:
			return nil, nil
		}
	}

	switch authType {
	case "bearer":
		if cfg.Bearer == nil {
			return nil, errors.New("auth.bearer is required when auth.type=bearer")
		}
		return NewBearer(cfg.Bearer), nil
	case "api_key":
		if cfg.APIKey == nil {
			return nil, errors.New("auth.api_key is required when auth.type=api_key")
		}
		return NewAPIKey(cfg.APIKey), nil
	case "aksk_hmac":
		if cfg.AKSKHMAC == nil {
			return nil, errors.New("auth.aksk_hmac is required when auth.type=aksk_hmac")
		}
		return NewAKSKHMAC(cfg.AKSKHMAC), nil
	default:
		return nil, fmt.Errorf("unsupported auth.type: %s", cfg.Type)
	}
}
