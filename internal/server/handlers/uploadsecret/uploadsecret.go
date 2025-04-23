package uploadsecret

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/general/logger"
	"github.com/dmitrovia/passkeeper/internal/general/models/secret"
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/uploadsecret/uploadsecretattr"
	"github.com/dmitrovia/passkeeper/internal/server/models/ctxm"
	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	"github.com/dmitrovia/passkeeper/internal/server/service"
)

var errEmptyData = errors.New("data is empty")

const (
	statusISE = http.StatusInternalServerError
)

type UploadSecret struct {
	secretService service.SecretService
	attr          *uploadsecretattr.UploadSecretAttr
}

func NewHandler(
	secretService service.SecretService,
	inAttr *uploadsecretattr.UploadSecretAttr,
) *UploadSecret {
	return &UploadSecret{
		attr:          inAttr,
		secretService: secretService,
	}
}

func (h *UploadSecret) UploadSecretHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	user, ok := req.Context().Value(ctxm.UserKey).(*userm.User)
	if !ok || user == nil {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	secret, err := h.getReqData(req)
	if err != nil {
		h.setErr(writer, err, "getReqData")

		return
	}

	ctx, cancel := context.WithTimeout(
		req.Context(), h.attr.Dbtimeout)
	defer cancel()

	secret.User = user

	err = h.secretService.CreateSecret(ctx, secret)
	if err != nil {
		h.setErr(writer, err, "CreateSecret")

		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (h *UploadSecret) setErr(writer http.ResponseWriter,
	err error,
	method string,
) {
	writer.WriteHeader(statusISE)
	logger.LogE("upload->"+method, err, h.attr.ZapLogger)
}

func (h *UploadSecret) getReqData(
	req *http.Request,
) (*secret.Secret, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("getReqData->io.ReadAll: %w", err)
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("getReqData: %w", errEmptyData)
	}

	dec, err := rsa.Decrypt(&body, h.attr.DecKey)
	if err != nil {
		return nil, fmt.Errorf("getReqData->Decrypt: %w", err)
	}

	secret := &secret.Secret{}

	err = json.Unmarshal(*dec, &secret)
	if err != nil {
		return nil, fmt.Errorf("getReqData->JU: %w", err)
	}

	err = req.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("getReqData->RBC: %w", err)
	}

	return secret, nil
}
