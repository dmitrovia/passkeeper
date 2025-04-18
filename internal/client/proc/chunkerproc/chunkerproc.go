package chunkerproc

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/dmitrovia/passkeeper/internal/client/proc/chunkerproc/chunkerpa"
	"github.com/dmitrovia/passkeeper/internal/general/compress"
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
		go cp.runWorker(indexChan)
	}
}

func (cp *ChunkerProc) runWorker(
	indexChan chan int,
) {
	defer cp.attr.WgSubProc.Done()

	for index := range indexChan {
		err := cp.toChunk(index)
		if err != nil {
			cp.attr.ErrChan <- err

			cp.attr.WorkerChunkWg.Done()
		}
	}
}

func (cp *ChunkerProc) toChunk(
	index int,
) error {
	offset := int64(index) * int64(cp.attr.ChunkSize)
	buffer := make([]byte, cp.attr.ChunkSize)

	_, err := cp.attr.ChFile.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("toChunk->Seek: %w", err)
	}

	bytesRead, err := cp.attr.ChFile.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("toChunk->Read: %w", err)
	}

	if bytesRead == 0 {
		cp.attr.WorkerChunkWg.Done()

		return nil
	}

	chBytes := buffer[:bytesRead]
	hash := md5.Sum(chBytes)
	fileName := fmt.Sprintf("%s.chunk.%d",
		cp.attr.FileName, index)
	encodeHash := hex.EncodeToString(hash[:])

	oldChunk,
		exists := cp.attr.CurrentMetadata[fileName]

	if exists || oldChunk.Hash == &encodeHash {
		cp.attr.WorkerChunkWg.Done()

		return nil
	}

	chunk := chunckmeta.NewChunkMeta(
		fileName, cp.attr.FileName, encodeHash, index, nil,
	)

	isCompress := strings.Contains(cp.attr.GzipFormats,
		cp.attr.FileFormat)
	if isCompress {
		compressData, err := compress.DeflateCompress(
			chBytes)
		if err != nil {
			return fmt.Errorf("toChunk->DC: %w", err)
		}

		chunk.Data = &compressData
	} else {
		chunk.Data = &chBytes
	}

	cp.attr.UploadChan <- *chunk

	return nil
}
