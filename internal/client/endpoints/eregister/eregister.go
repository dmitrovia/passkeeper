package eregister

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/eregister/eregisterattr"
)

type Register struct {
	attr *eregisterattr.RegisterAttr
}

func NewRegister(
	attr *eregisterattr.RegisterAttr,
) *Register {
	return &Register{attr: attr}
}

func (u *Register) RegisterUser(
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
		return nil, fmt.Errorf("RegisterUser->NRWC: %w", err)
	}

	resp, err := u.attr.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("RegisterUser->Do: %w", err)
	}

	return resp, nil
}
