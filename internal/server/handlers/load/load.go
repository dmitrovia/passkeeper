package load

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
	"github.com/dmitrovia/passkeeper/internal/general/validate"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/load/loadattr"
	"github.com/dmitrovia/passkeeper/internal/server/models/ctxm"
	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	"github.com/dmitrovia/passkeeper/internal/server/service"
)

const (
	statusISE = http.StatusInternalServerError
)

var errEmptyData = errors.New("data is empty")

type Load struct {
	metaService service.MetaService
	attr        *loadattr.LoadAttr
}

func NewHandler(
	s service.MetaService,
	inAttr *loadattr.LoadAttr,
) *Load {
	return &Load{metaService: s, attr: inAttr}
}

func (h *Load) InitLoadHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	user, ok := req.Context().Value(ctxm.UserKey).(*userm.User)
	if !ok || user == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	chunkReq, err := getReqData(req)
	if err != nil {
		h.setErr(writer, err, "getReqData")
		return
	}

	isValid := isValid(chunkReq)
	if !isValid {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(
		req.Context(), h.attr.Dbtimeout)
	defer cancel()

	body, err := h.getResponeBody(ctx, user, chunkReq)
	if err != nil {
		h.setErr(writer, err, "getResponeBody")
		return
	}

	_, err = writer.Write(*body)
	if err != nil {
		h.setErr(writer, err, "Write")
		return
	}

	writer.WriteHeader(http.StatusOK)
}

//nolint:cyclop
func isValid(chunk *chunckmeta.ChunkMeta,
) bool {
	notUsesFields := chunk.ID != nil ||
		chunk.FilePath != nil ||
		chunk.User != nil ||
		chunk.Createddate != nil ||
		chunk.Index != nil ||
		chunk.Hash != nil

	if notUsesFields {
		return false
	}

	isNil := chunk.FileName == nil || chunk.OrigFileName == nil
	if isNil {
		return false
	}

	res := validate.IsValidOrigFileName(*chunk.OrigFileName)
	if !res {
		return false
	}

	isEmpty := *chunk.FileName == "" ||
		*chunk.OrigFileName == ""
	if isEmpty {
		return false
	}

	res = chunk.FNIsValid()

	return res
}

func (h *Load) getResponeBody(
	ctx context.Context,
	user *userm.User,
	chunk *chunckmeta.ChunkMeta,
) (*[]byte, error) {
	meta, _,
		err := h.metaService.GetMetaByClientFileNameOptimized(
		ctx, user.ID, *chunk.FileName)
	if err != nil {
		return nil, fmt.Errorf("GRB->GMBCO: %w", err)
	}

	data, err := os.ReadFile(*meta.FilePath)
	if err != nil {
		return nil, fmt.Errorf("GRB->ReadFile: %w", err)
	}

	meta.Data = &data

	marshall, err := json.Marshal(meta)
	if err != nil {
		return nil, fmt.Errorf("GRB->Marshal: %w", err)
	}

	return &marshall, nil
}

func (h *Load) setErr(writer http.ResponseWriter,
	err error,
	method string,
) {
	writer.WriteHeader(statusISE)
	logger.LogE("initLoad->"+method, err, h.attr.ZapLogger)
}

func getReqData(
	req *http.Request,
) (*chunckmeta.ChunkMeta, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("GRD->io.ReadAll: %w", err)
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("GRD: %w", errEmptyData)
	}

	chunk := &chunckmeta.ChunkMeta{}

	err = json.Unmarshal(body, &chunk)
	if err != nil {
		return nil, fmt.Errorf("GRD->json.Unmarshal: %w", err)
	}

	err = req.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("GRD->req.Body.Close: %w", err)
	}

	return chunk, nil
}
