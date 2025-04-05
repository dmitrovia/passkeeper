package chunkerproc

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/chunkerpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type ChunkerProc struct {
	attr *chunkerpa.ChunkerProcAttr
}

func NewProc(attr *chunkerpa.ChunkerProcAttr) *ChunkerProc {
	return &ChunkerProc{
		attr: attr,
	}
}

func (cp *ChunkerProc) RunProcess() error {
	fmt.Println("ChunkerProc run")
	defer fmt.Println("ChunkerProc end")

	cp.runWorkerPoolChunker()

	return nil
}

func (cp *ChunkerProc) runWorkerPoolChunker() {
	indexChan := make(chan int, cp.attr.CntChunks)

	for i := range cp.attr.CntChunks {
		indexChan <- i
	}

	close(indexChan)

	for range cp.attr.CountWorkersChunker {
		go cp.toChuck(indexChan,
			cp.attr.UploadChan,
			cp.attr.ErrChan)
	}
}

func (cp *ChunkerProc) toChuck(
	indexChan chan int,
	uploadChan chan chunckmeta.ChunkMeta,
	errChan chan error,
) {
	defer cp.attr.Wgroup.Done()

	for index := range indexChan {
		offset := int64(index) * int64(cp.attr.ChunkSize)
		buffer := make([]byte,
			cp.attr.ChunkSize)

		_, err := cp.attr.ChFile.Seek(offset, 0)
		if err != nil {
			errChan <- err

			return
		}

		bytesRead, err := cp.attr.ChFile.Read(buffer)
		if err != nil && errors.Is(err, io.EOF) {
			errChan <- err

			return
		}

		if bytesRead == 0 {
			return
		}

		chBytes := buffer[:bytesRead]
		hash := md5.Sum(chBytes)

		chunk := chunckmeta.NewChunkMeta(
			fmt.Sprintf("%s.chunk.%d", cp.attr.FilePath, index),
			hex.EncodeToString(hash[:]),
			index,
			&chBytes,
		)
		uploadChan <- *chunk
	}
}
