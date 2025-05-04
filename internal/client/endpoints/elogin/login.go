package elogin

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/elogin/eloginattr"
)

type Login struct {
	attr *eloginattr.LoginAttr
}

func NewEndpoint(
	attr *eloginattr.LoginAttr,
) *Login {
	return &Login{attr: attr}
}

func (u *Login) CallEndpoint(
	ctx context.Context,
) (
	*http.Response,
	error,
) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		u.attr.URL,
		bytes.NewReader(*u.attr.Data))
	if err != nil {
		return nil, fmt.Errorf("CallEndpoint->NRWC: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := u.attr.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("CallEndpoint->Do: %w", err)
	}

	return resp, nil
}
