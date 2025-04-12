package upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

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
	attr        *uploadattr.UploadAttr
}

func NewUploadHandler(
	s service.FileService,
	inAttr *uploadattr.UploadAttr,
) *Upload {
	return &Upload{fileService: s, attr: inAttr}
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

	fmt.Println(chunk)

	writer.WriteHeader(http.StatusOK)
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

	err = json.Unmarshal(body, chunk)
	if err != nil {
		return nil, fmt.Errorf("getReqData->JU: %w", err)
	}

	err = req.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("getReqData->RBC: %w", err)
	}

	return chunk, nil
}
