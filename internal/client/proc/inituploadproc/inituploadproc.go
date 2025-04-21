package inituploadproc

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/proc/inituploadproc/inituploadprocattr.go"
)

var errSNOK = errors.New("status is not OK")

type InitUploadProc struct {
	attr *inituploadprocattr.InitUploadProcAttr
}

func NewProc(
	attr *inituploadprocattr.InitUploadProcAttr,
) *InitUploadProc {
	return &InitUploadProc{
		attr: attr,
	}
}

func (proc *InitUploadProc) RunProcess() error {
	fmt.Println("InitUploadProc run")
	defer fmt.Println("InitUploadProc end")

	ctx, cancel := context.WithTimeout(
		context.Background(), proc.attr.ReqTimeout)
	defer cancel()

	resp, err := proc.attr.Inituploader.InitUpload(ctx)
	if err != nil {
		return fmt.Errorf("RP->InitUpload: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("RP->InitUpload: %w", errSNOK)
	}

	return nil
}
