package initload

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/general/compress"
	"github.com/dmitrovia/passkeeper/internal/general/logger"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/initload/initloadattr"
	"github.com/dmitrovia/passkeeper/internal/server/models/ctxm"
	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	"github.com/dmitrovia/passkeeper/internal/server/service"
)

const (
	statusISE = http.StatusInternalServerError
)

type InitLoad struct {
	metaService service.MetaService
	attr        *initloadattr.InitLoadAttr
}

func NewLoadHandler(
	s service.MetaService,
	inAttr *initloadattr.InitLoadAttr,
) *InitLoad {
	return &InitLoad{metaService: s, attr: inAttr}
}

func (h *InitLoad) InitLoadHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	user, ok := req.Context().Value(ctxm.UserKey).(*userm.User)
	if !ok || user == nil {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	ctx, cancel := context.WithTimeout(
		req.Context(), h.attr.Dbtimeout)
	defer cancel()

	body, err := h.getResponeBody(ctx, user)
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

func (h *InitLoad) getResponeBody(
	ctx context.Context,
	user *userm.User,
) (*[]byte, error) {
	metas, _, err := h.metaService.GetMetaByClientOptimized(
		ctx, user.ID)
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

func (h *InitLoad) setErr(writer http.ResponseWriter,
	err error,
	method string,
) {
	writer.WriteHeader(statusISE)
	logger.LogE("initLoad->"+method, err, h.attr.ZapLogger)
}
