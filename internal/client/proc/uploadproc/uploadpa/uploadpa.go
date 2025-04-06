package uploadpa

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/euploader"
	"github.com/dmitrovia/passkeeper/internal/client/endpoints/euploader/euploaderattr"
	"github.com/dmitrovia/passkeeper/internal/client/metamanager"
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type UploadProcAttr struct {
	CountWorkersUpload int
	Wgroup             *sync.WaitGroup
	ReqTimeout         time.Duration
	UploadChan         chan chunckmeta.ChunkMeta
	CurrentMetadata    map[string]chunckmeta.ChunkMeta
	UploaderAttr       *euploaderattr.UploaderAttr
	Uploader           *euploader.Uploader
	Mutex              *sync.Mutex
	ErrChan            chan error
	Client             *http.Client
	ServerURL          string
}

func (upa *UploadProcAttr) Init(
	attr *clientpa.ClientProcAttr,
) error {
	upa.Client = &http.Client{}
	upa.Mutex = &sync.Mutex{}

	metaManager := metamanager.NewMetaManager(
		attr.MetaPath)

	metadata, err := metaManager.LoadMetadata()
	if err != nil {
		return fmt.Errorf("Init->LoadMetadata: %w", err)
	}

	upa.ReqTimeout = attr.ReqTimeout
	upa.ServerURL = attr.ServerAddr
	upa.CurrentMetadata = metadata
	upa.UploaderAttr = &euploaderattr.UploaderAttr{}
	url := upa.ServerURL + "/upload"
	upa.UploaderAttr.Init(url, upa.Client)
	upa.Uploader = euploader.NewUploader(upa.UploaderAttr)

	upa.CountWorkersUpload = attr.CountWorkersUpload

	return nil
}
