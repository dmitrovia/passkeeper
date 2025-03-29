package interaction

import (
	"fmt"

	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/clientpa"
)

func RunProcess(clientAttr *clientpa.ClientProcAttr) error {
	fmt.Println("InteractionProc run")
	defer fmt.Println("InteractionProc end")
	fmt.Println(clientAttr)

	return nil
}
