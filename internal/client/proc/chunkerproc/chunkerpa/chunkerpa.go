package chunkerpa

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

const (
	defChunkSize = 1024 * 1024        // 1MB
	def10MB      = 1024 * 1024 * 10   // 10MB
	def100MB     = 1024 * 1024 * 100  // 100MB
	def1GB       = 1024 * 1024 * 1024 // 1GB
)

type ChunkerProcAttr struct {
	CountWorkersChunker int
	ChunkSize           int
	CntChunks           int
	FileSize            int64
	FilePath            string
	FileName            string
	FileFormat          string
	GzipFormats         string
	WorkerChunkWg       *sync.WaitGroup
	ChFile              *os.File
	UploadChan          chan *chunckmeta.ChunkMeta
	ErrChan             chan error
	CurrentMetadata     map[string]chunckmeta.ChunkMeta
	Aes256keyBytes      []byte
}

func (cpa *ChunkerProcAttr) setChunkSize() {
	fsz := cpa.FileSize

	if fsz <= def10MB {
		cpa.ChunkSize = defChunkSize
		return
	}

	if fsz <= def100MB {
		cpa.ChunkSize = def10MB
		return
	}

	cpa.ChunkSize = def100MB
}

func (cpa *ChunkerProcAttr) Init(
	attr *clientpa.ClientProcAttr,
) error {
	cpa.FilePath = attr.
		SelectFilePath

	file, err := os.Open(cpa.FilePath)
	if err != nil {
		return fmt.Errorf("RP->os.Open: %w", err)
	}

	cpa.ChFile = file

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("RP->Stat: %w", err)
	}

	cpa.FileName = fileInfo.Name()
	cpa.FileFormat = filepath.Ext(cpa.FilePath)

	cpa.FileSize = fileInfo.Size()
	cpa.setChunkSize()
	cz := int64(cpa.ChunkSize)
	cntChunks := int(cpa.FileSize / cz)

	if fileInfo.Size()%cz != 0 {
		cntChunks++
	}

	cpa.GzipFormats = attr.GzipFormats
	cpa.CntChunks = cntChunks

	cpa.CountWorkersChunker = attr.CountWorkersChunker
	cpa.Aes256keyBytes = attr.Aes256keyBytes

	return nil
}
