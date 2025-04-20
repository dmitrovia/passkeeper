package initsingleloadproc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/dmitrovia/passkeeper/internal/client/proc/initsingleloadproc/initsingleloadprocattr"
	"github.com/dmitrovia/passkeeper/internal/general/compress"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

var errSNOK = errors.New("status is not OK")

type InitSingleProc struct {
	attr *initsingleloadprocattr.InitSingleLoadProcAttr
}

func NewProc(
	attr *initsingleloadprocattr.InitSingleLoadProcAttr,
) *InitSingleProc {
	return &InitSingleProc{
		attr: attr,
	}
}

func (proc *InitSingleProc) RunProcess() error {
	fmt.Println("InitSingleProc run")
	defer fmt.Println("InitSingleProc end")

	fmt.Println("Enter file name")

	var fileName string

	_, err1 := fmt.Fscan(os.Stdin, &fileName)
	if err1 != nil {
		return fmt.Errorf("RP->Fscan: %w", err1)
	}

	proc.attr.SpecificFileLoadName = fileName

	reqData := &apim.InInitSingleLoad{}
	reqData.FileName = proc.attr.SpecificFileLoadName

	marshal, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("RP->Marshal: %w", err)
	}

	proc.attr.InitSingleLoadAttr.Data = &marshal

	ctx, cancel := context.WithTimeout(
		context.Background(), proc.attr.ReqTimeout)
	defer cancel()

	resp, err := proc.attr.InitSingleLoad.InitSingleLoad(ctx)
	if err != nil {
		return fmt.Errorf("RP->InitSingleLoad: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("RP->InitSingleLoad: %w", errSNOK)

		return err
	}

	err = proc.parseResp(resp)
	if err != nil {
		return fmt.Errorf("RP->parseResp: %w", err)
	}

	return nil
}

func (proc *InitSingleProc) parseResp(
	response *http.Response,
) error {
	out, err := compress.DeflateDecompress(response.Body)
	if err != nil {
		fmt.Println("parseResp->DeflateDecompress: %w", err)
	}

	metas := make(map[string]chunckmeta.ChunkMeta)

	err = json.Unmarshal(out, &metas)
	if err != nil {
		return fmt.Errorf("parseResp->Unmarshal: %w", err)
	}

	proc.attr.LoadMetadata = metas

	return nil
}
