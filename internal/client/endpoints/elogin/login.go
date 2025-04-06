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

func NewLogin(
	attr *eloginattr.LoginAttr,
) *Login {
	return &Login{attr: attr}
}

func (u *Login) LoginUser(
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
		return nil, fmt.Errorf("LoginUser->NRWC: %w", err)
	}

	resp, err := u.attr.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("LoginUser->Do: %w", err)
	}

	return resp, nil
}
