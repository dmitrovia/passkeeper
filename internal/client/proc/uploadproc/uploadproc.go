package uploadproc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/euploader"
	"github.com/dmitrovia/passkeeper/internal/client/endpoints/euploader/euploaderattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadproc/uploadpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
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

	go up.awaitClose()

	up.runWorkerPoolUpload()

	return nil
}

func (up *UploadProc) runWorkerPoolUpload() {
	for range up.attr.CountWorkersUpload {
		go up.runWorker()
	}
}

func (up *UploadProc) awaitClose() {
	up.attr.WorkerChunkWg.Wait()
	close(up.attr.UploadChan)
}

func (up *UploadProc) runWorker() {
	for chunk := range up.attr.UploadChan {
		up.uploadChunk(&chunk)
	}
}

func (up *UploadProc) toJSON(
	chunk *chunckmeta.ChunkMeta,
) (*[]byte, error) {
	metricMarshall, err := json.Marshal(chunk)
	if err != nil {
		return nil, fmt.Errorf("toCompressJSON->Marshal: %w", err)
	}

	return &metricMarshall, nil
}

func (up *UploadProc) uploadChunk(
	chunk *chunckmeta.ChunkMeta,
) {
	defer up.attr.WorkerChunkWg.Done()
	defer chunk.ClearData()

	client := &http.Client{}
	uplattr := &euploaderattr.UploaderAttr{}
	uplattr.URL = up.attr.ServerURL + "/api/user/upload"

	data, err := up.toJSON(chunk)
	if err != nil {
		up.attr.ErrChan <- err

		return
	}

	uplattr.Init(uplattr.URL, client,
		up.attr.AuthToken, data)

	uploader := euploader.NewUploader(uplattr)

	ctx, cancel := context.WithTimeout(
		context.Background(), up.attr.ReqTimeout)
	defer cancel()

	resp, err := uploader.UploadChunk(ctx)
	if err != nil {
		up.attr.ErrChan <- err

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		err := fmt.Errorf("RWP->UploadChunk: %w", errSNOK)
		up.attr.ErrChan <- err

		return
	}

	up.attr.Mutex.Lock()
	up.attr.CurrentMetadata[*chunk.FileName] = *chunk
	up.attr.UploadedMetadata[*chunk.FileName] = *chunk
	up.attr.Mutex.Unlock()
}
