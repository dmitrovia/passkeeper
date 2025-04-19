package einitload

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/einitload/einitloadattr"
)

type InitLoad struct {
	attr *einitloadattr.InitLoadAttr
}

func NewInitLoad(
	attr *einitloadattr.InitLoadAttr,
) *InitLoad {
	return &InitLoad{attr: attr}
}

func (u *InitLoad) InitLoad(
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
		return nil, fmt.Errorf("InitLoad->NRWC: %w", err)
	}

	req.Header.Set("Authorization", u.attr.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.attr.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("InitLoad->Do: %w", err)
	}

	return resp, nil
}
