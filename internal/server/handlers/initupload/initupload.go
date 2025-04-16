package initupload

import (
	"errors"
	"net/http"
	"os"

	"github.com/dmitrovia/passkeeper/internal/general/logger"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/initupload/inituploadattr"
	"github.com/dmitrovia/passkeeper/internal/server/models/ctxm"
	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	"github.com/dmitrovia/passkeeper/internal/server/service"
)

const (
	statusISE = http.StatusInternalServerError
)

type InitUpload struct {
	fileService service.FileService
	attr        *inituploadattr.InitUploadAttr
}

func NewHandler(
	s service.FileService,
	inAttr *inituploadattr.InitUploadAttr,
) *InitUpload {
	return &InitUpload{fileService: s, attr: inAttr}
}

func (h *InitUpload) InitUploadHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	user, ok := req.Context().Value(ctxm.UserKey).(*userm.User)
	if !ok || user == nil {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	path := h.attr.SaveFilesPath + "/" + *user.Login
	if _, err := os.Stat(path); errors.Is(
		err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			h.setErr(writer, err, "Mkdir")

			return
		}
	}

	writer.WriteHeader(http.StatusOK)
}

func (h *InitUpload) setErr(writer http.ResponseWriter,
	err error,
	method string,
) {
	writer.WriteHeader(statusISE)
	logger.LogE("initupload->"+method, err, h.attr.ZapLogger)
}
