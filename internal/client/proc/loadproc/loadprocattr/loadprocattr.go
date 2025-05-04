package loadprocattr

import (
	"sync"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type LoadProcAttr struct {
	CountWorkersLoad int
	WorkerChunkWg    *sync.WaitGroup
	ReqTimeout       time.Duration
	LoadChan         chan *chunckmeta.ChunkMeta
	LoadedMetadata   map[string]chunckmeta.ChunkMeta
	ErrChan          chan error
	ServerURL        string
	AuthToken        string
	TempFilesPath    string
	GzipFormats      string
	Aes256keyBytes   []byte
	ClientProcAttr   *clientpa.ClientProcAttr
}

func (lpa *LoadProcAttr) Init(
	attr *clientpa.ClientProcAttr,
) {
	lpa.ClientProcAttr = attr
	lpa.ReqTimeout = attr.ReqTimeout
	lpa.ServerURL = attr.ServerAddr
	lpa.LoadedMetadata = make(
		map[string]chunckmeta.ChunkMeta)
	lpa.CountWorkersLoad = attr.CountWorkersLoad
	lpa.AuthToken = attr.AuthToken
	lpa.TempFilesPath = attr.TempFilesPath
	lpa.GzipFormats = attr.GzipFormats
	lpa.Aes256keyBytes = attr.Aes256keyBytes
}
