package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/quanguachong/mcp-servers/http/internal/types"
)

type akskHMACApplier struct {
	accessKey       string
	secretKey       string
	timestamp       string
	accessKeyHeader string
	signatureHeader string
	timestampHeader string
}

func NewAKSKHMAC(cfg *types.AKSKHMACAuth) Applier {
	return &akskHMACApplier{
		accessKey:       cfg.AccessKey,
		secretKey:       cfg.SecretKey,
		timestamp:       cfg.Timestamp,
		accessKeyHeader: defaultIfEmpty(cfg.AccessKeyHeader, "ak"),
		signatureHeader: defaultIfEmpty(cfg.SignatureHeader, "sk"),
		timestampHeader: defaultIfEmpty(cfg.TimestampHeader, "X-Timestamp"),
	}
}

func (a *akskHMACApplier) Apply(req *http.Request, body []byte) error {
	if a.accessKey == "" || a.secretKey == "" {
		return errors.New("aksk access_key and secret_key cannot be empty")
	}

	ts := a.timestamp
	if ts == "" {
		ts = time.Now().UTC().Format(time.RFC3339)
	}

	sum := sha256.Sum256(body)
	bodyHash := hex.EncodeToString(sum[:])
	canonical := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s",
		strings.ToUpper(req.Method),
		req.URL.Path,
		req.URL.RawQuery,
		bodyHash,
		ts,
	)
	mac := hmac.New(sha256.New, []byte(a.secretKey))
	if _, err := mac.Write([]byte(canonical)); err != nil {
		return err
	}
	signature := hex.EncodeToString(mac.Sum(nil))

	req.Header.Set(a.accessKeyHeader, a.accessKey)
	req.Header.Set(a.timestampHeader, ts)
	req.Header.Set(a.signatureHeader, signature)
	return nil
}

func defaultIfEmpty(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}
