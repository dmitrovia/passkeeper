package loginprocattr

import (
	"net/http"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/elogin"
	"github.com/dmitrovia/passkeeper/internal/client/endpoints/elogin/eloginattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
)

type LoginProcAttr struct {
	ReqTimeout    time.Duration
	Client        *http.Client
	ServerURL     string
	TokenSavePath string
	LoginAttr     *eloginattr.LoginAttr
	Login         *elogin.Login
	AttrClintProc *clientpa.ClientProcAttr
	EncKey        *[]byte
}

func (rpa *LoginProcAttr) Init(
	attr *clientpa.ClientProcAttr,
) error {
	httpClient := attr.GetHTTPClient()
	rpa.Client = &httpClient

	rpa.ReqTimeout = attr.ReqTimeout
	rpa.ServerURL = attr.ServerAddr
	rpa.LoginAttr = &eloginattr.LoginAttr{}
	url := rpa.ServerURL + "/api/user/login"
	rpa.LoginAttr.Init(url, rpa.Client)
	rpa.Login = elogin.NewEndpoint(rpa.LoginAttr)
	rpa.TokenSavePath = attr.AuthTokenPath
	rpa.AttrClintProc = attr
	rpa.EncKey = &attr.PublicKey

	return nil
}
