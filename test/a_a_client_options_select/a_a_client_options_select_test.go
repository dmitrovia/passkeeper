package a_a_client_options_select_test

import (
	"sync"
	"testing"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc/interactionpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/testm"
)

func runNewProc(t *testing.T, option string, isAuth bool) {
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
	testdata := &testm.TestData{}
	testdata.TestChooseProcInput = option
	newAttr.AttrClintProc.SelectedProc = &option
	newAttr.AttrClintProc.TestData = testdata

	interp := interactionproc.NewProc(newAttr)

	go func() {
		<-time.After(500 * time.Microsecond)

		exitV := "99"
		newAttr.AttrClintProc.SelectedProc = &exitV
	}()

	err = interp.RunProcess()
	if err != nil {
		t.Errorf("TestMain->RP: %v", err)

		return
	}
}

func TestMain(t *testing.T) {
	t.Helper()
	t.Parallel()

	runNewProc(t, "99", true)
	runNewProc(t, "1", false)
	runNewProc(t, "2", false)
	runNewProc(t, "3", true)
	runNewProc(t, "4", true)
	runNewProc(t, "5", true)
	runNewProc(t, "6", true)
	runNewProc(t, "7", true)
	runNewProc(t, "8", true)
	runNewProc(t, "999", true)
	runNewProc(t, "34543543", true)
}
