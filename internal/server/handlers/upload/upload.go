package upload

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dmitrovia/passkeeper/internal/general/logger"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/upload/uploadattr"
	"github.com/dmitrovia/passkeeper/internal/server/models/ctxm"
	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	"github.com/dmitrovia/passkeeper/internal/server/service"
)

var errEmptyData = errors.New("data is empty")

const (
	statusISE = http.StatusInternalServerError
)

type Upload struct {
	fileService service.FileService
	metaService service.MetaService
	attr        *uploadattr.UploadAttr
}

func NewHandler(
	s service.FileService,
	metaService service.MetaService,
	inAttr *uploadattr.UploadAttr,
) *Upload {
	return &Upload{
		fileService: s,
		attr:        inAttr,
		metaService: metaService,
	}
}

func (h *Upload) UploadHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	user, ok := req.Context().Value(ctxm.UserKey).(*userm.User)
	if !ok || user == nil {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	chunk, err := h.getReqData(req)
	if err != nil {
		h.setErr(writer, err, "getReqData")

		return
	}

	ctx, cancel := context.WithTimeout(
		req.Context(), h.attr.Dbtimeout)
	defer cancel()

	chunk.User = user
	newPath := fmt.Sprintf("%s/%s/%s",
		h.attr.FilesStoragePath, *user.Login, *chunk.FileName)
	chunk.FilePath = &newPath

	err = h.createChunkFile(chunk)
	if err != nil {
		h.setErr(writer, err, "createChunkFile")

		return
	}

	err = h.metaService.CreateMeta(ctx, chunk)
	if err != nil {
		h.setErr(writer, err, "CreateMeta")

		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (h *Upload) createChunkFile(
	chunk *chunckmeta.ChunkMeta,
) error {
	chunkFile, err := os.Create(*chunk.FilePath)
	if err != nil {
		return fmt.Errorf("createChunkFile->Create: %w", err)
	}

	_, err = chunkFile.Write(*chunk.Data)
	if err != nil {
		return fmt.Errorf("createChunkFile->Write: %w", err)
	}

	return nil
}

func (h *Upload) setErr(writer http.ResponseWriter,
	err error,
	method string,
) {
	writer.WriteHeader(statusISE)
	logger.LogE("upload->"+method, err, h.attr.ZapLogger)
}

func (h *Upload) getReqData(
	req *http.Request,
) (*chunckmeta.ChunkMeta, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("getReqData->io.ReadAll: %w", err)
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("getReqData: %w", errEmptyData)
	}

	chunk := &chunckmeta.ChunkMeta{}

	err = json.Unmarshal(body, &chunk)
	if err != nil {
		return nil, fmt.Errorf("getReqData->JU: %w", err)
	}

	err = req.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("getReqData->RBC: %w", err)
	}

	return chunk, nil
}
