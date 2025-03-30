package interaction

import (
	"fmt"

	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/uploadpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadproc"
)

func RunProcess(clientAttr *clientpa.ClientProcAttr) error {
	fmt.Println("InteractionProc run")
	defer fmt.Println("InteractionProc end")
	fmt.Println(clientAttr)

	attr := &uploadpa.UploadProcAttr{}

	err := uploadproc.RunProcess(attr)
	if err != nil {
		return fmt.Errorf("RP->uploadproc.RP: %w", err)
	}

	return nil
}
