package euploadsercret

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/euploadsercret/euploadsercretattr"
)

type UploadSecret struct {
	attr *euploadsercretattr.UploadSecretAttr
}

func NewEndpoint(
	attr *euploadsercretattr.UploadSecretAttr,
) *UploadSecret {
	return &UploadSecret{attr: attr}
}

func (u *UploadSecret) CallEndpoint(
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

	req.Header.Set("Authorization", u.attr.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.attr.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("CallEndpoint->Do: %w", err)
	}

	return resp, nil
}
