package registerprocpattr

import (
	"net/http"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/eregister"
	"github.com/dmitrovia/passkeeper/internal/client/endpoints/eregister/eregisterattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
)

type RegisterProcAttr struct {
	ReqTimeout   time.Duration
	Client       *http.Client
	ServerURL    string
	RegisterAttr *eregisterattr.RegisterAttr
	Register     *eregister.Register
	EncKey       *[]byte
}

func (rpa *RegisterProcAttr) Init(
	attr *clientpa.ClientProcAttr,
) error {
	rpa.Client = &http.Client{}

	rpa.ReqTimeout = attr.ReqTimeout
	rpa.ServerURL = attr.ServerAddr
	rpa.RegisterAttr = &eregisterattr.RegisterAttr{}
	url := rpa.ServerURL + "/api/user/register"
	rpa.RegisterAttr.Init(url, rpa.Client)
	rpa.Register = eregister.NewRegister(rpa.RegisterAttr)
	rpa.EncKey = &attr.PrivateKey

	return nil
}
