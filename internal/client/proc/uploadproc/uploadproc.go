package uploadproc

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/uploader"
	"github.com/dmitrovia/passkeeper/internal/client/models/endpointsattrs/uploaderattrs"
	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/uploadpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

var errSNOK = errors.New("status is not OK")

func RunProcess(attr *uploadpa.UploadProcAttr) error {
	fmt.Println("UploadProc run")
	defer fmt.Println("UploadProc end")

	ctx, cancel := context.WithTimeout(
		context.Background(), attr.ReqTimeout)
	defer cancel()

	cnt := len(attr.Chunks)
	chunkChan := make(chan chunckmeta.ChunkMeta, cnt)
	errChan := make(chan error, cnt)

	for _, chunk := range attr.Chunks {
		attr.Wgroup.Add(1)
		chunkChan <- chunk
	}

	runWorkerPool(ctx, chunkChan, errChan, attr)

	attr.Wgroup.Wait()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func runWorkerPool(ctx context.Context,
	chunkChan chan chunckmeta.ChunkMeta,
	errChan chan error,
	attr *uploadpa.UploadProcAttr,
) {
	client := &http.Client{}
	uploaderAttr := &uploaderattrs.UploaderAttr{}
	uploaderAttr.Init(attr.ServerURL, client)
	upl := uploader.NewUploader(uploaderAttr)

	for range attr.CountWorkers {
		go func() {
			for chunk := range chunkChan {
				defer attr.Wgroup.Done()

				newHash := chunk.Hash

				// attr.Mutex.Lock()
				oldChunk, exists := attr.Metadata[*chunk.FileName]
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

				attr.Mutex.Lock()
				attr.Metadata[*chunk.FileName] = chunk
				attr.Mutex.Unlock()
			}
		}()
	}

	close(chunkChan)
	close(errChan)
}
