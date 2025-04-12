package initupload

import (
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/server/handlers/initupload/inituploadattr"
	"github.com/dmitrovia/passkeeper/internal/server/models/ctxm"
	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	"github.com/dmitrovia/passkeeper/internal/server/service"
)

type InitUpload struct {
	fileService service.FileService
	attr        *inituploadattr.InitUploadAttr
}

func NewUploadHandler(
	s service.FileService,
	inAttr *inituploadattr.InitUploadAttr,
) *InitUpload {
	return &InitUpload{fileService: s, attr: inAttr}
}

func (h *InitUpload) UploadHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	user, ok := req.Context().Value(ctxm.UserKey).(*userm.User)
	if !ok || user == nil {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	fmt.Println(user)

	writer.WriteHeader(http.StatusOK)
}
