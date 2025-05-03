package getsecrets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/general/compress"
	"github.com/dmitrovia/passkeeper/internal/general/logger"
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/getsecrets/getsecretsattr"
	"github.com/dmitrovia/passkeeper/internal/server/models/ctxm"
	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	"github.com/dmitrovia/passkeeper/internal/server/service"
)

const (
	statusISE = http.StatusInternalServerError
)

type GetSecrets struct {
	secretService service.SecretService
	attr          *getsecretsattr.GetSecretAttr
}

func NewHandler(
	s service.SecretService,
	inAttr *getsecretsattr.GetSecretAttr,
) *GetSecrets {
	return &GetSecrets{secretService: s, attr: inAttr}
}

func (h *GetSecrets) GetSecretsHandler(
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

func (h *GetSecrets) getResponeBody(
	ctx context.Context,
	user *userm.User,
) (*[]byte, error) {
	secrets,
		_, err := h.secretService.GetSecretByClientOptimized(
		ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("getResponeBody->GSBCO: %w", err)
	}

	marshall, err := json.Marshal(secrets)
	if err != nil {
		return nil, fmt.Errorf("getResponeBody->Marshal: %w", err)
	}

	compress, err := compress.DeflateCompress(marshall)
	if err != nil {
		return nil, fmt.Errorf("getResponeBody->DC: %w", err)
	}

	encrypt, err := rsa.Encrypt(&compress, h.attr.EncKey)
	if err != nil {
		return nil, fmt.Errorf("getResponeBody->Encrypt: %w", err)
	}

	return encrypt, nil
}

func (h *GetSecrets) setErr(writer http.ResponseWriter,
	err error,
	method string,
) {
	writer.WriteHeader(statusISE)
	logger.LogE("initLoad->"+method, err, h.attr.ZapLogger)
}
