package getsecretbyidprocattr

import (
	"net/http"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/egetsecretbyid"
	"github.com/dmitrovia/passkeeper/internal/client/endpoints/egetsecretbyid/egetsecretbyidattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
)

type GetSecretByIDProcAttr struct {
	ReqTimeout           time.Duration
	Client               *http.Client
	ServerURL            string
	GetSecretByID        *egetsecretbyid.GetSecretBYID
	GetSecretByIDAttr    *egetsecretbyidattr.GetSecretByIDAttr
	AuthToken            string
	SpecificFileLoadName string
	EncKey               *[]byte
	DecKey               *[]byte
}

func (rpa *GetSecretByIDProcAttr) Init(
	attr *clientpa.ClientProcAttr,
) {
	rpa.AuthToken = attr.AuthToken
	rpa.Client = &http.Client{}
	rpa.ReqTimeout = attr.ReqTimeout
	rpa.ServerURL = attr.ServerAddr
	sattr := &egetsecretbyidattr.GetSecretByIDAttr{}
	rpa.GetSecretByIDAttr = sattr
	url := rpa.ServerURL + "/api/user/getsecretbyid"
	rpa.GetSecretByIDAttr.Init(url, rpa.Client, rpa.AuthToken)
	rpa.GetSecretByID = egetsecretbyid.NewEndpoint(
		rpa.GetSecretByIDAttr)
	rpa.EncKey = &attr.PublicKey
	rpa.DecKey = &attr.PrivateKey
}
