package a_b_client_options_fscan_errors_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc/interactionpa"
)

var errCntErrors = errors.New("count of errors")

func runNewProc(t *testing.T,
	option string,
	isAuth bool,
	errors *[]error,
) {
	t.Helper()

	clintattr := &clientpa.ClientProcAttr{}

	err := clintattr.Init(false)
	if err != nil {
		t.Errorf("TestMain->Init: %v", err)

		return
	}

	newAttr := &interactionpa.InteractionProcAttr{}
	newAttr.AttrClintProc = clintattr
	newAttr.AttrClintProc.IsAuth = isAuth
	newAttr.WgSubProc = &sync.WaitGroup{}
	newAttr.AttrClintProc.SelectedProc = &option

	interp := interactionproc.NewProc(newAttr)

	err = interp.RunProcess()
	if err != nil {
		*errors = append(*errors, err)

		return
	}
}

func TestMain(t *testing.T) {
	t.Helper()
	t.Parallel()

	errors := make([]error, 0)

	runNewProc(t, "99", true, &errors)
	runNewProc(t, "1", false, &errors)
	runNewProc(t, "2", false, &errors)
	runNewProc(t, "3", true, &errors)
	runNewProc(t, "4", true, &errors)
	runNewProc(t, "5", true, &errors)
	runNewProc(t, "6", true, &errors)
	runNewProc(t, "7", true, &errors)
	runNewProc(t, "8", true, &errors)
	runNewProc(t, "999", true, &errors)
	runNewProc(t, "34543543", true, &errors)

	if len(errors) != 10 {
		t.Errorf("TestMain->Init: %v", errCntErrors)
	}
}
