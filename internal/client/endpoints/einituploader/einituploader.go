package einituploader

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/einituploader/einituploaderattr"
)

type InitUploader struct {
	attr *einituploaderattr.InitUploadAttr
}

func NewEndpoint(
	attr *einituploaderattr.InitUploadAttr,
) *InitUploader {
	return &InitUploader{attr: attr}
}

func (u *InitUploader) CallEndpoint(
	ctx context.Context,
) (
	*http.Response,
	error,
) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
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
