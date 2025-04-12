package euploader

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/euploader/euploaderattr"
)

type Uploader struct {
	attr *euploaderattr.UploaderAttr
}

func NewUploader(
	attr *euploaderattr.UploaderAttr,
) *Uploader {
	return &Uploader{attr: attr}
}

func (u *Uploader) UploadChunk(
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
		return nil, fmt.Errorf("UploadChunk->NRWC: %w", err)
	}

	req.Header.Set("Authorization", u.attr.Token)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.attr.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("UploadChunk->Do: %w", err)
	}

	return resp, nil
}
