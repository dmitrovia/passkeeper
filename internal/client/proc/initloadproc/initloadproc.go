package initloadproc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/proc/initloadproc/initloadprocattr"
	"github.com/dmitrovia/passkeeper/internal/general/compress"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

var errSNOK = errors.New("status is not OK")

type InitProc struct {
	attr *initloadprocattr.InitLoadProcAttr
}

func NewProc(
	attr *initloadprocattr.InitLoadProcAttr,
) *InitProc {
	return &InitProc{
		attr: attr,
	}
}

func (proc *InitProc) RunProcess() error {
	fmt.Println("InitProc run")
	defer fmt.Println("InitProc end")

	ctx, cancel := context.WithTimeout(
		context.Background(), proc.attr.ReqTimeout)
	defer cancel()

	resp, err := proc.attr.InitLoad.InitLoad(ctx)
	if err != nil {
		return fmt.Errorf("RP->InitLoad: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("RP->InitLoad: %w", errSNOK)

		return err
	}

	err = proc.parseResp(resp)
	if err != nil {
		return fmt.Errorf("RP->parseResp: %w", err)
	}

	return nil
}

func (proc *InitProc) parseResp(
	response *http.Response,
) error {
	out, err := compress.DeflateDecompress(response.Body)
	if err != nil {
		return fmt.Errorf("parseResp->DeflateDecompress: %w", err)
	}

	metas := make(map[string]chunckmeta.ChunkMeta)

	err = json.Unmarshal(out, &metas)
	if err != nil {
		return fmt.Errorf("parseResp->Unmarshal: %w", err)
	}

	proc.attr.LoadMetadata = metas

	return nil
}
