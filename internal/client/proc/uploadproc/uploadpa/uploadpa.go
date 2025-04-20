package uploadpa

import (
	"sync"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type UploadProcAttr struct {
	CountWorkersUpload int
	WorkerChunkWg      *sync.WaitGroup
	ReqTimeout         time.Duration
	UploadChan         chan chunckmeta.ChunkMeta
	UploadedMetadata   map[string]chunckmeta.ChunkMeta
	CurrentMetadata    map[string]chunckmeta.ChunkMeta
	Mutex              *sync.Mutex
	ErrChan            chan error
	ServerURL          string
	CountChunk         int
	AuthToken          string
}

func (upa *UploadProcAttr) Init(
	attr *clientpa.ClientProcAttr,
) error {
	upa.Mutex = &sync.Mutex{}

	upa.ReqTimeout = attr.ReqTimeout
	upa.ServerURL = attr.ServerAddr
	upa.UploadedMetadata = make(
		map[string]chunckmeta.ChunkMeta)
	upa.CountWorkersUpload = attr.CountWorkersUpload
	upa.AuthToken = attr.AuthToken

	return nil
}
