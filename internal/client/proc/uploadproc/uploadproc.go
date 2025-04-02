package uploadproc

import (
	"fmt"

	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/uploadpa"
)

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

	return nil
}
