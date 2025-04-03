package uploader

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/models/endpointsattrs/uploaderattrs"
)

type Uploader struct {
	attr *uploaderattrs.UploaderAttr
}

func NewUploader(
	inAttr *uploaderattrs.UploaderAttr,
) *Uploader {
	return &Uploader{attr: inAttr}
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
		u.attr.ServerURL,
		bytes.NewReader(*u.attr.Data))
	if err != nil {
		return nil, fmt.Errorf("UploadChunk->NRWC: %w", err)
	}

	resp, err := u.attr.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("UploadChunk->Do: %w", err)
	}

	return resp, nil
}
