package uploadpa

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/uploader"
	"github.com/dmitrovia/passkeeper/internal/client/metamanager"
	"github.com/dmitrovia/passkeeper/internal/client/models/endpointsattrs/uploaderattrs"
	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/clientpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type UploadProcAttr struct {
	CountWorkersUpload int
	Wgroup             *sync.WaitGroup
	ReqTimeout         time.Duration
	UploadChan         chan chunckmeta.ChunkMeta
	CurrentMetadata    map[string]chunckmeta.ChunkMeta
	UploaderAttr       *uploaderattrs.UploaderAttr
	Uploader           *uploader.Uploader
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

	upa.ServerURL = attr.ServerAddr
	upa.CurrentMetadata = metadata
	upa.UploaderAttr = &uploaderattrs.UploaderAttr{}
	upa.UploaderAttr.Init(upa.ServerURL, upa.Client)
	upa.Uploader = uploader.NewUploader(upa.UploaderAttr)

	upa.CountWorkersUpload = attr.CountWorkersUpload

	return nil
}
