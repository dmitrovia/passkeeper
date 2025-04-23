package egetsecrets

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/egetsecrets/egetsecretsattr"
)

type GetSecrets struct {
	attr *egetsecretsattr.GetSecretsAttr
}

func NewEndpoint(
	attr *egetsecretsattr.GetSecretsAttr,
) *GetSecrets {
	return &GetSecrets{attr: attr}
}

func (u *GetSecrets) CallEndpoint(
	ctx context.Context,
) (
	*http.Response,
	error,
) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		u.attr.URL,
		&bytes.Buffer{})
	if err != nil {
		return nil, fmt.Errorf("CallEndpoint->NRWC: %w", err)
	}

	req.Header.Set("Authorization", u.attr.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.attr.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("CallEndpoint->Do: %w", err)
	}

	return resp, nil
}
