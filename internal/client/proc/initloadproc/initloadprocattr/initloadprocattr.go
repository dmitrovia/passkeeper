package initloadprocattr

import (
	"net/http"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/einitload"
	"github.com/dmitrovia/passkeeper/internal/client/endpoints/einitload/einitloadattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type InitLoadProcAttr struct {
	ReqTimeout   time.Duration
	Client       *http.Client
	ServerURL    string
	InitLoad     *einitload.InitLoad
	InitLoadAttr *einitloadattr.InitLoadAttr
	AuthToken    string
	LoadMetadata map[string]chunckmeta.ChunkMeta
}

func (ilp *InitLoadProcAttr) Init(
	attr *clientpa.ClientProcAttr,
) {
	ilp.AuthToken = attr.AuthToken
	httpClient := attr.GetHTTPClient()
	ilp.Client = &httpClient
	ilp.ReqTimeout = attr.ReqTimeout
	ilp.ServerURL = attr.ServerAddr
	sattr := &einitloadattr.InitLoadAttr{}
	ilp.InitLoadAttr = sattr
	url := ilp.ServerURL + "/api/user/initload"
	ilp.InitLoadAttr.Init(url, ilp.Client, ilp.AuthToken)
	ilp.InitLoad = einitload.NewEndpoint(
		ilp.InitLoadAttr)
}
