package chunkerproc

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/uploader"
	"github.com/dmitrovia/passkeeper/internal/client/models/endpointsattrs/uploaderattrs"
	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/chunkerpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

var errSNOK = errors.New("status is not OK")

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

	ctx, cancel := context.WithTimeout(
		context.Background(), cp.attr.ReqTimeout)
	defer cancel()

	file, err := os.Open(cp.attr.FilePath)
	if err != nil {
		return fmt.Errorf("RP->os.Open: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("RP->Stat: %w", err)
	}

	numChunks := int(fileInfo.Size() /
		int64(cp.attr.ChunkSize))
	if fileInfo.Size()%int64(cp.attr.ChunkSize) != 0 {
		numChunks++
	}

	chunkChan := make(chan chunckmeta.ChunkMeta, numChunks)
	errChan := make(chan error, numChunks)

	go cp.runWorkerPoolUpload(ctx, chunkChan, errChan)
	go cp.runWorkerPoolChunker(errChan,
		chunkChan, numChunks, file)

	go func() {
		cp.attr.Wgroup.Wait()
		close(chunkChan)
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (cp *ChunkerProc) runWorkerPoolChunker(
	errChan chan error,
	chunkChan chan chunckmeta.ChunkMeta,
	numChunks int,
	file *os.File,
) {
	indexChan := make(chan int, numChunks)
	defer close(indexChan)

	for i := range numChunks {
		indexChan <- i
	}

	for range cp.attr.CountWorkersChunker {
		cp.attr.Wgroup.Add(1)

		go func() {
			defer cp.attr.Wgroup.Done()

			for index := range indexChan {
				offset := int64(index) * int64(cp.attr.ChunkSize)
				buffer := make([]byte,
					cp.attr.ChunkSize)

				_, err := file.Seek(offset, 0)
				if err != nil {
					errChan <- err

					return
				}

				bytesRead, err := file.Read(buffer)
				if err != nil && errors.Is(err, io.EOF) {
					errChan <- err

					return
				}

				if bytesRead > 0 {
					hash := md5.Sum(buffer[:bytesRead])
					hashString := hex.EncodeToString(hash[:])

					chunkFileName := fmt.Sprintf("%s.chunk.%d",
						cp.attr.FilePath, index)

					chunkFile, err := os.Create(chunkFileName)
					if err != nil {
						errChan <- err

						return
					}

					_, err = chunkFile.Write(buffer[:bytesRead])
					if err != nil {
						errChan <- err

						return
					}

					chunk := chunckmeta.ChunkMeta{
						FileName: &chunkFileName,
						Hash:     &hashString,
						Index:    &index,
					}

					chunkFile.Close()

					chunkChan <- chunk
				}
			}
		}()
	}
}

func (cp *ChunkerProc) runWorkerPoolUpload(
	ctx context.Context,
	chunkChan chan chunckmeta.ChunkMeta,
	errChan chan error,
) {
	client := &http.Client{}
	uploaderAttr := &uploaderattrs.UploaderAttr{}
	uploaderAttr.Init(cp.attr.ServerURL, client)
	upl := uploader.NewUploader(uploaderAttr)

	for range cp.attr.CountWorkersUpload {
		go func() {
			for chunk := range chunkChan {
				defer cp.attr.Wgroup.Done()

				newHash := chunk.Hash

				// attr.Mutex.Lock()
				oldChunk,
					exists := cp.attr.CurrentMetadata[*chunk.FileName]
				// attr.Mutex.Unlock()

				if exists || oldChunk.Hash == newHash {
					return
				}

				resp, err := upl.UploadChunk(ctx, chunk)
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
		}()
	}

	close(chunkChan)
	close(errChan)
}
