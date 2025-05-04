package aaclientoptions_test

import (
	"sync"
	"testing"

	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc/interactionpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/testm"
)

//nolint:cyclop,funlen
func TestMain(t *testing.T) {
	t.Helper()
	t.Parallel()

	clintattr := &clientpa.ClientProcAttr{}

	err := clintattr.Init()
	if err != nil {
		t.Errorf("TestMain->Init: %v", err)

		return
	}

	newAttr := &interactionpa.InteractionProcAttr{}
	newAttr.AttrClintProc = clintattr
	newAttr.WgSubProc = &sync.WaitGroup{}
	testdata := &testm.TestData{}
	newAttr.AttrClintProc.TestData = testdata

	interp := interactionproc.NewProc(newAttr)

	testdata.TestChooseProcInput = "99"

	err = interp.RunProcess()
	if err != nil {
		t.Errorf("TestMain->RP: %v", err)

		return
	}

	testdata.TestChooseProcInput = "1"

	err = interp.RunProcess()
	if err != nil {
		t.Errorf("TestMain->RP: %v", err)

		return
	}

	testdata.TestChooseProcInput = "2"

	err = interp.RunProcess()
	if err != nil {
		t.Errorf("TestMain->RP: %v", err)

		return
	}

	testdata.TestChooseProcInput = "3"

	err = interp.RunProcess()
	if err != nil {
		t.Errorf("TestMain->RP: %v", err)

		return
	}

	testdata.TestChooseProcInput = "4"

	err = interp.RunProcess()
	if err != nil {
		t.Errorf("TestMain->RP: %v", err)

		return
	}

	testdata.TestChooseProcInput = "5"

	err = interp.RunProcess()
	if err != nil {
		t.Errorf("TestMain->RP: %v", err)

		return
	}

	testdata.TestChooseProcInput = "6"

	err = interp.RunProcess()
	if err != nil {
		t.Errorf("TestMain->RP: %v", err)

		return
	}

	testdata.TestChooseProcInput = "7"

	err = interp.RunProcess()
	if err != nil {
		t.Errorf("TestMain->RP: %v", err)

		return
	}

	testdata.TestChooseProcInput = "8"

	err = interp.RunProcess()
	if err != nil {
		t.Errorf("TestMain->RP: %v", err)

		return
	}

	testdata.TestChooseProcInput = "999"

	err = interp.RunProcess()
	if err != nil {
		t.Errorf("TestMain->RP: %v", err)

		return
	}

	testdata.TestChooseProcInput = "34543543"

	err = interp.RunProcess()
	if err != nil {
		t.Errorf("TestMain->RP: %v", err)

		return
	}
}
