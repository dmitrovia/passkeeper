package chunkerproc

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/chunkerpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

var errSNOK = errors.New("status is not OK")

const wgCount int = 2

type ChunkerProc struct {
	attr *chunkerpa.ChunkerProcAttr
}

func NewProc(attr *chunkerpa.ChunkerProcAttr) *ChunkerProc {
	return &ChunkerProc{
		attr: attr,
	}
}

func (cp *ChunkerProc) RunProcess() error {
	if cp.attr == nil {
		cp.attr = &chunkerpa.ChunkerProcAttr{}
	}

	err := cp.attr.Init()
	if err != nil {
		return fmt.Errorf("RP->attr.Init: %w", err)
	}

	uploadChan := make(chan chunckmeta.ChunkMeta,
		cp.attr.CntChunks)
	errChan := make(chan error, cp.attr.CntChunks)

	cp.attr.Wgroup.Add(cp.attr.CntChunks * wgCount)

	go cp.runWorkerPoolChunker(errChan, uploadChan)
	go cp.runWorkerPoolUpload(uploadChan, errChan)

	cp.attr.Wgroup.Wait()
	close(uploadChan)
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (cp *ChunkerProc) runWorkerPoolChunker(
	errChan chan error,
	uploadChan chan chunckmeta.ChunkMeta,
) {
	defer cp.attr.ChFile.Close()

	indexChan := make(chan int, cp.attr.CntChunks)
	defer close(indexChan)

	for i := range cp.attr.CntChunks {
		indexChan <- i
	}

	for range cp.attr.CountWorkersChunker {
		go cp.toChuck(indexChan, uploadChan, errChan)
	}
}

func (cp *ChunkerProc) toChuck(
	indexChan chan int,
	uploadChan chan chunckmeta.ChunkMeta,
	errChan chan error,
) {
	defer cp.attr.Wgroup.Add(-1)

	for index := range indexChan {
		offset := int64(index) * int64(cp.attr.ChunkSize)
		buffer := make([]byte,
			cp.attr.ChunkSize)

		_, err := cp.attr.ChFile.Seek(offset, 0)
		if err != nil {
			errChan <- err

			cp.attr.Wgroup.Add(-1)

			return
		}

		bytesRead, err := cp.attr.ChFile.Read(buffer)
		if err != nil && errors.Is(err, io.EOF) {
			errChan <- err

			cp.attr.Wgroup.Add(-1)

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

func (cp *ChunkerProc) runWorkerPoolUpload(
	uploadChan chan chunckmeta.ChunkMeta,
	errChan chan error,
) {
	for range cp.attr.CountWorkersUpload {
		go cp.toUpload(uploadChan, errChan)
	}
}

func (cp *ChunkerProc) toUpload(
	uploadChan chan chunckmeta.ChunkMeta,
	errChan chan error,
) {
	for chunk := range uploadChan {
		defer cp.attr.Wgroup.Done()

		ctx, cancel := context.WithTimeout(
			context.Background(), cp.attr.ReqTimeout)
		defer cancel()

		newHash := chunk.Hash

		cp.attr.Mutex.Lock()
		oldChunk,
			exists := cp.attr.CurrentMetadata[*chunk.FileName]
		cp.attr.Mutex.Unlock()

		if exists || oldChunk.Hash == newHash {
			return
		}

		cp.attr.UploaderAttr.Data = chunk.Data
		defer chunk.ClearData()

		resp, err := cp.attr.Uploader.UploadChunk(ctx)
		if err != nil {
			errChan <- err

			return
		}

		if resp.StatusCode != http.StatusOK {
			err := fmt.Errorf("RWP->UploadChunk: %w", errSNOK)
			errChan <- err

			return
		}

		resp.Body.Close()

		cp.attr.Mutex.Lock()
		cp.attr.CurrentMetadata[*chunk.FileName] = chunk
		cp.attr.Mutex.Unlock()
	}
}
