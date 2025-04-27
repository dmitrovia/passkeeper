package chunkerproc

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/dmitrovia/passkeeper/internal/client/proc/chunkerproc/chunkerpa"
	"github.com/dmitrovia/passkeeper/internal/general/aes256"
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

func (proc *ChunkerProc) RunProcess() error {
	fmt.Println("ChunkerProc run")
	defer fmt.Println("ChunkerProc end")

	proc.runWorkerPoolChunker()

	return nil
}

func (proc *ChunkerProc) runWorkerPoolChunker() {
	indexChan := make(chan int, proc.attr.CntChunks)

	for i := range proc.attr.CntChunks {
		indexChan <- i
	}

	close(indexChan)

	for range proc.attr.CountWorkersChunker {
		go proc.runWorker(indexChan)
	}
}

func (proc *ChunkerProc) runWorker(
	indexChan chan int,
) {
	for index := range indexChan {
		err := proc.toChunk(index)
		if err != nil {
			proc.attr.ErrChan <- err

			proc.attr.WorkerChunkWg.Done()
		}
	}
}

func (proc *ChunkerProc) toChunk(
	index int,
) error {
	offset := int64(index) * int64(proc.attr.ChunkSize)
	buffer := make([]byte, proc.attr.ChunkSize)

	_, err := proc.attr.ChFile.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("toChunk->Seek: %w", err)
	}

	bytesRead, err := proc.attr.ChFile.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("toChunk->Read: %w", err)
	}

	if bytesRead == 0 {
		proc.attr.WorkerChunkWg.Done()

		return nil
	}

	chBytes := buffer[:bytesRead]
	hash := md5.Sum(chBytes)
	fileName := fmt.Sprintf("%s.chunk.%d",
		proc.attr.FileName, index)
	encodeHash := hex.EncodeToString(hash[:])

	oldChunk,
		exists := proc.attr.CurrentMetadata[fileName]

	if exists || oldChunk.Hash == &encodeHash {
		proc.attr.WorkerChunkWg.Done()

		return nil
	}

	chunk := chunckmeta.NewChunkMeta(
		fileName, proc.attr.FileName, encodeHash, index, nil,
	)

	err = proc.compressAndEncrypt(chunk, &chBytes)
	if err != nil {
		return fmt.Errorf("toChunk->compressAndEncrypt: %w", err)
	}

	proc.attr.UploadChan <- chunk

	return nil
}

func (proc *ChunkerProc) compressAndEncrypt(
	chunk *chunckmeta.ChunkMeta,
	chBytes *[]byte,
) error {
	isCompress := strings.Contains(proc.attr.GzipFormats,
		proc.attr.FileFormat)
	if isCompress {
		compressData, err := compress.DeflateCompress(
			*chBytes)
		if err != nil {
			return fmt.Errorf("toChunk->DC: %w", err)
		}

		chunk.Data = &compressData
	} else {
		chunk.Data = chBytes
	}

	dec, err := aes256.Encrypt(chunk.Data,
		&proc.attr.Aes256keyBytes)
	if err != nil {
		return fmt.Errorf("PRASF->aes256Decrypt: %w", err)
	}

	chunk.Data = dec

	return nil
}
