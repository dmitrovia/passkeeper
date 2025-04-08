package upload

import (
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/server/handlers/upload/uploadattr"
	"github.com/dmitrovia/passkeeper/internal/server/service"
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
	_ *http.Request,
) {
	writer.WriteHeader(http.StatusOK)
}
