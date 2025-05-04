package egetsecretbyid

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/egetsecretbyid/egetsecretbyidattr"
)

type GetSecretBYID struct {
	attr *egetsecretbyidattr.GetSecretByIDAttr
}

func NewEndpoint(
	attr *egetsecretbyidattr.GetSecretByIDAttr,
) *GetSecretBYID {
	return &GetSecretBYID{attr: attr}
}

func (u *GetSecretBYID) CallEndpoint(
	ctx context.Context,
) (
	*http.Response,
	error,
) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		u.attr.URL,
		bytes.NewReader(*u.attr.Data))
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
