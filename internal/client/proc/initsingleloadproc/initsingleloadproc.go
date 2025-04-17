package initsingleloadproc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/proc/initsingleloadproc/initsingleloadprocattr"
	"github.com/dmitrovia/passkeeper/internal/general/compress"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

var errSNOK = errors.New("status is not OK")

type InitSingleProc struct {
	attr *initsingleloadprocattr.InitUploadProcAttr
}

func NewProc(
	attr *initsingleloadprocattr.InitUploadProcAttr,
) *InitSingleProc {
	return &InitSingleProc{
		attr: attr,
	}
}

func (isp *InitSingleProc) RunProcess() error {
	fmt.Println("InitSingleProc run")
	defer fmt.Println("InitSingleProc end")

	reqData := &apim.InInitSingleLoad{}
	reqData.FileName = isp.attr.SpecificFileLoadName

	marshal, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("RP->Marshal: %w", err)
	}

	isp.attr.InitSingleLoadAttr.Data = &marshal

	ctx, cancel := context.WithTimeout(
		context.Background(), isp.attr.ReqTimeout)
	defer cancel()

	resp, err := isp.attr.InitSingleLoad.InitSingleLoad(ctx)
	if err != nil {
		return fmt.Errorf("RP->InitSingleLoad: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("RP->InitSingleLoad: %w", errSNOK)

		return err
	}

	err = isp.parseResp(resp)
	if err != nil {
		return fmt.Errorf("RP->parseResp: %w", err)
	}

	return nil
}

func (isp *InitSingleProc) parseResp(
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

	isp.attr.LoadMetadata = metas

	return nil
}
