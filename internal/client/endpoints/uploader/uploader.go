package uploader

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/dmitrovia/passkeeper/internal/client/models/endpointsattrs/uploaderattrs"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
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
	chunk chunckmeta.ChunkMeta,
) (
	*http.Response,
	error,
) {
	data, err := os.ReadFile(*chunk.FileName)
	if err != nil {
		return nil, fmt.Errorf("UploadChunk->ReadFile: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		u.attr.ServerURL,
		bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("UploadChunk->NRWC: %w", err)
	}

	resp, err := u.attr.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("UploadChunk->Do: %w", err)
	}

	return resp, nil
}
