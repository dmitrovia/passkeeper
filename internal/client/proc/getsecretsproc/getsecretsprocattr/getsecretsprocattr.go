package getsecretsprocattr

import (
	"net/http"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/egetsecrets"
	"github.com/dmitrovia/passkeeper/internal/client/endpoints/egetsecrets/egetsecretsattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
)

type GetSecretsProcAttr struct {
	ReqTimeout           time.Duration
	Client               *http.Client
	ServerURL            string
	GetSecrets           *egetsecrets.GetSecrets
	GetSecretsAttr       *egetsecretsattr.GetSecretsAttr
	AuthToken            string
	SpecificFileLoadName string
	DecKey               *[]byte
}

func (rpa *GetSecretsProcAttr) Init(
	attr *clientpa.ClientProcAttr,
) {
	rpa.AuthToken = attr.AuthToken
	httpClient := attr.GetHTTPClient()
	rpa.Client = &httpClient
	rpa.ReqTimeout = attr.ReqTimeout
	rpa.ServerURL = attr.ServerAddr
	sattr := &egetsecretsattr.GetSecretsAttr{}
	rpa.GetSecretsAttr = sattr
	url := rpa.ServerURL + "/api/user/getsecrets"
	rpa.GetSecretsAttr.Init(url, rpa.Client, rpa.AuthToken)
	rpa.GetSecrets = egetsecrets.NewEndpoint(
		rpa.GetSecretsAttr)
	rpa.DecKey = &attr.PrivateKey
}
