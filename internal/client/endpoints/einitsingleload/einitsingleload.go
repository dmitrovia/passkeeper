package einitsingleload

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/einitsingleload/einitsingleloadattr"
)

type InitSingleLoad struct {
	attr *einitsingleloadattr.InitSingleLoadAttr
}

func NewInitSingleLoad(
	attr *einitsingleloadattr.InitSingleLoadAttr,
) *InitSingleLoad {
	return &InitSingleLoad{attr: attr}
}

func (u *InitSingleLoad) InitSingleLoad(
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
		return nil, fmt.Errorf("InitSingleLoad->NRWC: %w", err)
	}

	req.Header.Set("Authorization", u.attr.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.attr.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("InitSingleLoad->Do: %w", err)
	}

	return resp, nil
}
