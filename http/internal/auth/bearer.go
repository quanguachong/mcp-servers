package auth

import (
	"errors"
	"fmt"
	"net/http"

	"http-requests/internal/types"
)

type bearerApplier struct {
	token string
}

func NewBearer(cfg *types.BearerAuth) Applier {
	return &bearerApplier{token: cfg.Token}
}

func (b *bearerApplier) Apply(req *http.Request, _ []byte) error {
	if b.token == "" {
		return errors.New("bearer token cannot be empty")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", b.token))
	return nil
}
