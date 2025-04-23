package uploadsecretprocattr

import (
	"net/http"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/euploadsercret"
	"github.com/dmitrovia/passkeeper/internal/client/endpoints/euploadsercret/euploadsercretattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
)

type UploadSecretProcAttr struct {
	ReqTimeout       time.Duration
	Client           *http.Client
	ServerURL        string
	UploadSecret     *euploadsercret.UploadSecret
	UploadSecretAttr *euploadsercretattr.UploadSecretAttr
	AttrClintProc    *clientpa.ClientProcAttr
	EncKey           *[]byte
}

func (rpa *UploadSecretProcAttr) Init(
	attr *clientpa.ClientProcAttr,
) {
	rpa.Client = &http.Client{}

	rpa.ReqTimeout = attr.ReqTimeout
	rpa.ServerURL = attr.ServerAddr
	rpa.UploadSecretAttr = &euploadsercretattr.
		UploadSecretAttr{}
	url := rpa.ServerURL + "/api/user/uploadsecret"
	rpa.UploadSecretAttr.Init(url, rpa.Client, attr.AuthToken)
	rpa.UploadSecret = euploadsercret.NewEndpoint(
		rpa.UploadSecretAttr)
	rpa.AttrClintProc = attr
	rpa.EncKey = &attr.PublicKey
}
