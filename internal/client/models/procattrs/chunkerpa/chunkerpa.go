package chunkerpa

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/uploader"
	"github.com/dmitrovia/passkeeper/internal/client/models/endpointsattrs/uploaderattrs"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type ChunkerProcAttr struct {
	CountWorkersUpload  int
	CountWorkersChunker int
	ChunkSize           int
	CntChunks           int
	FilePath            string
	ServerURL           string
	CurrentMetadata     map[string]chunckmeta.ChunkMeta
	ReqTimeout          time.Duration
	Wgroup              *sync.WaitGroup
	Mutex               *sync.Mutex
	ChFile              *os.File
	Client              *http.Client
	Uploader            *uploader.Uploader
	UploaderAttr        *uploaderattrs.UploaderAttr
}

func (cpa *ChunkerProcAttr) Init() error {
	cpa.Client = &http.Client{}

	cpa.UploaderAttr = &uploaderattrs.UploaderAttr{}
	cpa.UploaderAttr.Init(cpa.ServerURL, cpa.Client)
	cpa.Uploader = uploader.NewUploader(cpa.UploaderAttr)

	file, err := os.Open(cpa.FilePath)
	if err != nil {
		return fmt.Errorf("RP->os.Open: %w", err)
	}

	cpa.ChFile = file

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("RP->Stat: %w", err)
	}

	cz := int64(cpa.ChunkSize)
	cntChunks := int(fileInfo.Size() / cz)

	if fileInfo.Size()%cz != 0 {
		cntChunks++
	}

	cpa.CntChunks = cntChunks

	return nil
}
