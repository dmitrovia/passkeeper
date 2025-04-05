package uploadproc

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadproc/uploadpa"
)

var errSNOK = errors.New("status is not OK")

type UploadProc struct {
	attr *uploadpa.UploadProcAttr
}

func NewProc(attr *uploadpa.UploadProcAttr,
) *UploadProc {
	return &UploadProc{
		attr: attr,
	}
}

func (up *UploadProc) RunProcess() error {
	fmt.Println("UploadProc run")
	defer fmt.Println("UploadProc end")

	up.runWorkerPoolUpload()

	return nil
}

func (up *UploadProc) runWorkerPoolUpload() {
	for range up.attr.CountWorkersUpload {
		go up.toUpload()
	}
}

func (up *UploadProc) toUpload() {
	defer up.attr.Wgroup.Done()

	for chunk := range up.attr.UploadChan {
		ctx, cancel := context.WithTimeout(
			context.Background(), up.attr.ReqTimeout)
		defer cancel()

		newHash := chunk.Hash

		up.attr.Mutex.Lock()
		oldChunk,
			exists := up.attr.CurrentMetadata[*chunk.FileName]
		up.attr.Mutex.Unlock()

		if exists || oldChunk.Hash == newHash {
			return
		}

		up.attr.UploaderAttr.Data = chunk.Data
		defer chunk.ClearData()

		resp, err := up.attr.Uploader.UploadChunk(ctx)
		if err != nil {
			up.attr.ErrChan <- err

			return
		}

		if resp.StatusCode != http.StatusOK {
			err := fmt.Errorf("RWP->UploadChunk: %w", errSNOK)
			up.attr.ErrChan <- err

			return
		}

		resp.Body.Close()

		up.attr.Mutex.Lock()
		up.attr.CurrentMetadata[*chunk.FileName] = chunk
		up.attr.Mutex.Unlock()
	}
}
