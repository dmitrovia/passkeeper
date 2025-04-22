package getsecretbyid

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
	"github.com/dmitrovia/passkeeper/internal/server/handlers/getsecrets/getsecretsattr"
	"github.com/dmitrovia/passkeeper/internal/server/models/ctxm"
	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	"github.com/dmitrovia/passkeeper/internal/server/service"
)

const (
	statusISE = http.StatusInternalServerError
)

var errEmptyData = errors.New("data is empty")

type GetSecretByID struct {
	secretService service.SecretService
	attr          *getsecretsattr.GetSecretAttr
}

func NewHandler(
	s service.SecretService,
	inAttr *getsecretsattr.GetSecretAttr,
) *GetSecretByID {
	return &GetSecretByID{secretService: s, attr: inAttr}
}

func (h *GetSecretByID) GetSecretByIDHadnler(
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

func (h *GetSecretByID) getResponeBody(
	ctx context.Context,
	user *userm.User,
	reqAttr *apim.InGetSecretByID,
) (*[]byte, error) {
	secret, _,
		err := h.secretService.
		GetSecretByClientIdentifierOptimized(
			ctx, user.ID, reqAttr.Identifier)
	if err != nil {
		return nil, fmt.Errorf("getResponeBody->GMBCO: %w", err)
	}

	marshall, err := json.Marshal(secret)
	if err != nil {
		return nil, fmt.Errorf("getResponeBody->Marshal: %w", err)
	}

	compress, err := compress.DeflateCompress(marshall)
	if err != nil {
		return nil, fmt.Errorf("getResponeBody->DC: %w", err)
	}

	return &compress, nil
}

func (h *GetSecretByID) setErr(writer http.ResponseWriter,
	err error,
	method string,
) {
	writer.WriteHeader(statusISE)
	logger.LogE("initLoad->"+method, err, h.attr.ZapLogger)
}

func (h *GetSecretByID) getReqData(
	req *http.Request,
) (*apim.InGetSecretByID, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("getReqData->io.ReadAll: %w", err)
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("getReqData: %w", errEmptyData)
	}

	reqData := &apim.InGetSecretByID{}

	err = json.Unmarshal(body, reqData)
	if err != nil {
		return nil, fmt.Errorf("getReqData->JU: %w", err)
	}

	err = req.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("getReqData->RBC: %w", err)
	}

	return reqData, nil
}
