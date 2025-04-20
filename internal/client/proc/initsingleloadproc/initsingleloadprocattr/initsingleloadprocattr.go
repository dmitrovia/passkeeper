package initsingleloadprocattr

import (
	"net/http"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/einitsingleload"
	"github.com/dmitrovia/passkeeper/internal/client/endpoints/einitsingleload/einitsingleloadattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type InitSingleLoadProcAttr struct {
	ReqTimeout         time.Duration
	Client             *http.Client
	ServerURL          string
	InitSingleLoad     *einitsingleload.InitSingleLoad
	InitSingleLoadAttr *einitsingleloadattr.
				InitSingleLoadAttr
	AuthToken            string
	SpecificFileLoadName string
	LoadMetadata         map[string]chunckmeta.ChunkMeta
}

func (rpa *InitSingleLoadProcAttr) Init(
	attr *clientpa.ClientProcAttr,
) {
	rpa.AuthToken = attr.AuthToken
	rpa.Client = &http.Client{}
	rpa.ReqTimeout = attr.ReqTimeout
	rpa.ServerURL = attr.ServerAddr
	sattr := &einitsingleloadattr.InitSingleLoadAttr{}
	rpa.InitSingleLoadAttr = sattr
	url := rpa.ServerURL + "/api/user/initsingleload"
	rpa.InitSingleLoadAttr.Init(url, rpa.Client, rpa.AuthToken)
	rpa.InitSingleLoad = einitsingleload.NewInitSingleLoad(
		rpa.InitSingleLoadAttr)
}
