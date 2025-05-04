package a_b_client_options_fscan_errors_test

import (
	"sync"
	"testing"

	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc/interactionpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/testm"
)

func runNewProc(t *testing.T,
	option string,
	isAuth bool,
	testd *testm.TestData,
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
	newAttr.AttrClintProc.TestData = testd

	interp := interactionproc.NewProc(newAttr)

	interp.SelectOption()
}

func TestMain(t *testing.T) {
	t.Helper()
	t.Parallel()

	const temp string = "test"

	testdata := &testm.TestData{}
	runNewProc(t, "99", true, testdata)

	testdata = &testm.TestData{}
	runNewProc(t, "1", false, testdata)

	testdata = &testm.TestData{}
	testdata.TestLoginInputRegister = temp
	runNewProc(t, "1", false, testdata)

	testdata = &testm.TestData{}
	runNewProc(t, "2", false, testdata)

	testdata = &testm.TestData{}
	testdata.TestLoginInputLogin = temp
	runNewProc(t, "2", false, testdata)

	testdata = &testm.TestData{}
	runNewProc(t, "3", true, testdata)

	testdata = &testm.TestData{}
	testdata.TestUploadAndChunksSelectModeInput = "1"
	runNewProc(t, "3", true, testdata)

	testdata = &testm.TestData{}
	runNewProc(t, "5", true, testdata)

	testdata = &testm.TestData{}
	testdata.TestChooseLoadTypeInput = "1"
	runNewProc(t, "5", true, testdata)

	testdata = &testm.TestData{}
	runNewProc(t, "6", true, testdata)

	testdata = &testm.TestData{}
	testdata.TestInIdentifierInputUploadSecret = temp
	runNewProc(t, "6", true, testdata)

	testdata = &testm.TestData{}
	runNewProc(t, "8", true, testdata)

	testdata = &testm.TestData{}
	runNewProc(t, "999", true, testdata)
	testdata = &testm.TestData{}
	runNewProc(t, "34543543", true, testdata)
}
