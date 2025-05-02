package interactionproc_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc/interactionpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/testm"
)

//nolint:funlen,cyclop
func TestMain(t *testing.T) {
	t.Helper()
	t.Parallel()

	time.Sleep(5 * time.Second)

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
	testdata.TestLoginInputRegister = "test" +
		randomString()
	testdata.TestPassInputRegister = "test" +
		randomString()
	newAttr.AttrClintProc.TestData = testdata

	interp := interactionproc.NewProc(newAttr)

	testdata.TestChooseProcInput = "99"

	err = interp.RunProcess()
	if err != nil {
		t.Errorf("TestMain->RP: %v", err)

		return
	}

	err = interp.RunRegister()
	if err != nil {
		t.Errorf("TestMain->RunRegister: %v", err)

		return
	}

	testdata.TestLoginInputLogin = testdata.
		TestLoginInputRegister
	testdata.TestPassInputLogin = testdata.
		TestPassInputRegister

	err = interp.RunLogin()
	if err != nil {
		t.Errorf("TestMain->RunLogin: %v", err)

		return
	}

	err = interp.RunLogout()
	if err != nil {
		t.Errorf("TestMain->RunLogin: %v", err)

		return
	}

	err = interp.RunLogin()
	if err != nil {
		t.Errorf("TestMain->RunLogin: %v", err)

		return
	}

	testdata.
		TestUploadAndChunksSelectModeInput = "2" // upload all

	err = interp.UploadAndChunksSelectMode()
	if err != nil {
		t.Errorf("TestMain->UploadAndChunksSelectMode1: %v", err)

		return
	}

	testdata.
		TestUploadAndChunksSelectModeInput = "1" // special
	testdata.TestUploadSingleFileInput = "upload_test"

	err = interp.UploadAndChunksSelectMode()
	if err != nil {
		t.Errorf("TestMain->UploadAndChunksSelectMode2: %v", err)

		return
	}

	testdata.
		TestSetRestrictionsInput = "2" // overwrite all files
	// when load

	testdata.TestChooseLoadTypeInput = "2" // load all

	err = interp.LoadAndChunksSelectMode()
	if err != nil {
		t.Errorf("TestMain->UploadAndChunksSelectMode1: %v", err)

		return
	}

	testdata.TestFileNameInput = "upload_test"
	testdata.TestChooseLoadTypeInput = "1" // load single

	err = interp.LoadAndChunksSelectMode()
	if err != nil {
		t.Errorf("TestMain->UploadAndChunksSelectMode1: %v", err)

		return
	}

	testdata.
		TestInIdentifierInputUploadSecret = "test" +
		randomString()
	testdata.
		TestInValueInputUploadSecret = "test" + randomString()

	err = interp.RunUploadSecret()
	if err != nil {
		t.Errorf("TestMain->RunUploadSecret: %v", err)

		return
	}

	err = interp.RunGetSecrets()
	if err != nil {
		t.Errorf("TestMain->RunGetSecrets: %v", err)

		return
	}

	testdata.TestInIdentifierInput = testdata.
		TestInIdentifierInputUploadSecret

	err = interp.RunGetSecretByID()
	if err != nil {
		t.Errorf("TestMain->RunGetSecrets: %v", err)

		return
	}
}

func randomString() string {
	letters := []rune(
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
