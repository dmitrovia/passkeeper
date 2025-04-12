package inituploadprocattr

import (
	"net/http"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/einituploader"
	"github.com/dmitrovia/passkeeper/internal/client/endpoints/einituploader/einituploaderattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
)

type InitUploadProcAttr struct {
	ReqTimeout       time.Duration
	Client           *http.Client
	ServerURL        string
	Inituploader     *einituploader.InitUploader
	Inituploaderattr *einituploaderattr.InitUploadAttr
	AuthToken        string
}

func (rpa *InitUploadProcAttr) Init(
	attr *clientpa.ClientProcAttr,
) {
	rpa.AuthToken = attr.AuthToken
	rpa.Client = &http.Client{}
	rpa.ReqTimeout = attr.ReqTimeout
	rpa.ServerURL = attr.ServerAddr
	rpa.Inituploaderattr = &einituploaderattr.InitUploadAttr{}
	url := rpa.ServerURL + "/api/user/initupload"
	rpa.Inituploaderattr.Init(url, rpa.Client, rpa.AuthToken)
	rpa.Inituploader = einituploader.NewUploader(
		rpa.Inituploaderattr)
}
