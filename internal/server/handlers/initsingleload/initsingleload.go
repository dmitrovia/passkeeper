package initsingleload

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/general/compress"
	"github.com/dmitrovia/passkeeper/internal/general/logger"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/general/validate"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/initsingleload/initsingleloadattr"
	"github.com/dmitrovia/passkeeper/internal/server/models/ctxm"
	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	"github.com/dmitrovia/passkeeper/internal/server/service"
)

const (
	statusISE = http.StatusInternalServerError
)

var errEmptyData = errors.New("data is empty")

type InitSingleLoad struct {
	metaService service.MetaService
	attr        *initsingleloadattr.InitSingleLoadAttr
}

func NewHandler(
	s service.MetaService,
	inAttr *initsingleloadattr.InitSingleLoadAttr,
) *InitSingleLoad {
	return &InitSingleLoad{metaService: s, attr: inAttr}
}

func (h *InitSingleLoad) InitSingleLoadHadnler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	user, ok := req.Context().Value(ctxm.UserKey).(*userm.User)
	if !ok || user == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	reqAttr, err := h.getReqData(req)
	if err != nil {
		h.setErr(writer, err, "getReqData")
		return
	}

	isValid := isValid(reqAttr)
	if !isValid {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(
		req.Context(), h.attr.Dbtimeout)
	defer cancel()

	body, err := h.getResponeBody(ctx, user, reqAttr)
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

func isValid(reqAttr *apim.InInitSingleLoad,
) bool {
	if reqAttr.FileName == "" {
		return false
	}

	res := validate.IsValidFileName(reqAttr.FileName)

	return res
}

func (h *InitSingleLoad) getResponeBody(
	ctx context.Context,
	user *userm.User,
	reqAttr *apim.InInitSingleLoad,
) (*[]byte, error) {
	metas, _,
		err := h.metaService.GetMetaByClientOrigFileNameOptimized(
		ctx, user.ID, reqAttr.FileName)
	if err != nil {
		return nil, fmt.Errorf("getResponeBody->GMBCO: %w", err)
	}

	marshall, err := json.Marshal(metas)
	if err != nil {
		return nil, fmt.Errorf("getResponeBody->Marshal: %w", err)
	}

	compress, err := compress.DeflateCompress(marshall)
	if err != nil {
		return nil, fmt.Errorf("getResponeBody->DC: %w", err)
	}

	return &compress, nil
}

func (h *InitSingleLoad) setErr(writer http.ResponseWriter,
	err error,
	method string,
) {
	writer.WriteHeader(statusISE)
	logger.LogE("initLoad->"+method, err, h.attr.ZapLogger)
}

func (h *InitSingleLoad) getReqData(
	req *http.Request,
) (*apim.InInitSingleLoad, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("getReqData->io.ReadAll: %w", err)
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("getReqData: %w", errEmptyData)
	}

	reqData := &apim.InInitSingleLoad{}

	err = json.Unmarshal(body, &reqData)
	if err != nil {
		return nil, fmt.Errorf("getReqData->JU: %w", err)
	}

	err = req.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("getReqData->RBC: %w", err)
	}

	return reqData, nil
}
