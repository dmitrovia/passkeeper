package loadproc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dmitrovia/passkeeper/internal/client/endpoints/eload"
	"github.com/dmitrovia/passkeeper/internal/client/endpoints/eload/eloadattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/loadproc/loadprocattr"
	"github.com/dmitrovia/passkeeper/internal/general/compress"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

var errSNOK = errors.New("status is not OK")

type LoadProc struct {
	attr *loadprocattr.LoadProcAttr
}

func NewProc(attr *loadprocattr.LoadProcAttr,
) *LoadProc {
	return &LoadProc{
		attr: attr,
	}
}

func (proc *LoadProc) RunProcess() error {
	fmt.Println("LoadProc run")
	defer fmt.Println("LoadProc end")

	proc.runWorkerPoolLoad()

	return nil
}

func (proc *LoadProc) runWorkerPoolLoad() {
	for range proc.attr.CountWorkersLoad {
		go proc.runWorker()
	}
}

func (proc *LoadProc) runWorker() {
	for chunk := range proc.attr.LoadChan {
		proc.loadChunk(&chunk)
	}
}

func (proc *LoadProc) toJSON(
	chunk *chunckmeta.ChunkMeta,
) (*[]byte, error) {
	metricMarshall, err := json.Marshal(chunk)
	if err != nil {
		return nil, fmt.Errorf("toJSON->Marshal: %w", err)
	}

	return &metricMarshall, nil
}

func (proc *LoadProc) loadChunk(
	chunk *chunckmeta.ChunkMeta,
) {
	defer proc.attr.WorkerChunkWg.Done()

	client := &http.Client{}
	loadattr := &eloadattr.LoadAttr{}
	loadattr.URL = proc.attr.ServerURL + "/api/user/load"

	data, err := proc.toJSON(chunk)
	if err != nil {
		proc.attr.ErrChan <- err

		return
	}

	loadattr.Init(loadattr.URL, client,
		proc.attr.AuthToken, data)

	loader := eload.NewLoader(loadattr)

	ctx, cancel := context.WithTimeout(
		context.Background(), proc.attr.ReqTimeout)
	defer cancel()

	resp, err := loader.LoadChunk(ctx)
	if err != nil {
		proc.attr.ErrChan <- err

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		err := fmt.Errorf("RWP->LoadChunk: %w", errSNOK)
		proc.attr.ErrChan <- err

		return
	}

	err = proc.parseRespAndSaveFile(chunk, resp)
	if err != nil {
		proc.attr.ErrChan <- err

		return
	}
}

func (proc *LoadProc) parseRespAndSaveFile(
	orig *chunckmeta.ChunkMeta,
	resp *http.Response,
) error {
	respChunk := &chunckmeta.ChunkMeta{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("PRASF->io.ReadAll: %w", err)
	}

	err = json.Unmarshal(body, &respChunk)
	if err != nil {
		return fmt.Errorf("PRASF->Unmarshal: %w", err)
	}

	decompress, err := compress.DeflateDecompress(
		bytes.NewReader(*respChunk.Data),
	)
	if err != nil {
		return fmt.Errorf("parseRespAndSaveFile->DD: %w", err)
	}

	orig.Data = &decompress
	newPath := fmt.Sprintf("%s/%s", proc.attr.TempFilesPath,
		*respChunk.FileName)
	orig.FilePath = &newPath
	orig.Hash = respChunk.Hash
	orig.Index = respChunk.Index
	orig.FileName = respChunk.FileName
	orig.OrigFileName = respChunk.OrigFileName

	err = proc.createChunkFile(respChunk)
	if err != nil {
		return fmt.Errorf("parseRespAndSaveFile->CCF: %w", err)
	}

	defer respChunk.ClearData()

	return nil
}

func (proc *LoadProc) createChunkFile(
	chunk *chunckmeta.ChunkMeta,
) error {
	chunkFile, err := os.Create(*chunk.FilePath)
	if err != nil {
		return fmt.Errorf("createChunkFile->Create: %w", err)
	}

	_, err = chunkFile.Write(*chunk.Data)
	if err != nil {
		return fmt.Errorf("createChunkFile->Write: %w", err)
	}

	return nil
}
