package chunkerpa

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type ChunkerProcAttr struct {
	CountWorkersChunker int
	ChunkSize           int
	CntChunks           int
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

func (cpa *ChunkerProcAttr) Init(
	attr *clientpa.ClientProcAttr,
) error {
	cpa.ChunkSize = attr.DefChunkSize
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

	cz := int64(cpa.ChunkSize)
	cntChunks := int(fileInfo.Size() / cz)

	if fileInfo.Size()%cz != 0 {
		cntChunks++
	}

	cpa.GzipFormats = attr.GzipFormats
	cpa.CntChunks = cntChunks

	cpa.CountWorkersChunker = attr.CountWorkersChunker
	cpa.Aes256keyBytes = attr.Aes256keyBytes

	return nil
}
