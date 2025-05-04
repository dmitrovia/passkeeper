package adnoauth_test

import (
	"errors"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc/interactionpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/testm"
)

var errCntErrors = errors.New("count of errors")

//nolint:funlen,cyclop
func TestNoAuth(t *testing.T) {
	t.Helper()
	t.Parallel()

	time.Sleep(20 * time.Second)

	clintattr := &clientpa.ClientProcAttr{}

	err := clintattr.Init(true)
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

	testdata.
		TestSetRestrictionsInput = "2" // overwrite all files
	// when load

	testdata.TestChooseLoadTypeInput = "2" // load all

	errors := make([]error, 0)

	err = interp.LoadAndChunksSelectMode()
	if err != nil {
		errors = append(errors, err)
	}

	testdata.TestChooseProcInput = "99"

	err = interp.RunProcess()
	if err != nil {
		t.Errorf("TestMain->RP1: %v", err)
	}

	testdata.TestChooseProcInput = "99"

	err = interp.RunProcess()
	if err != nil {
		t.Errorf("TestMain->RP2: %v", err)
	}

	testdata.
		TestUploadAndChunksSelectModeInput = "2" // upload all

	err = interp.UploadAndChunksSelectMode()
	if err != nil {
		errors = append(errors, err)
	}

	testdata.
		TestUploadAndChunksSelectModeInput = "1" // special
	testdata.TestUploadSingleFileInput = "upload_test"

	err = interp.UploadAndChunksSelectMode()
	if err != nil {
		errors = append(errors, err)
	}

	testdata.
		TestSetRestrictionsInput = "2" // overwrite all files
	// when load

	testdata.TestChooseLoadTypeInput = "2" // load all

	err = interp.LoadAndChunksSelectMode()
	if err != nil {
		errors = append(errors, err)
	}

	testdata.TestFileNameInput = "upload_test"
	testdata.TestChooseLoadTypeInput = "1" // load single

	err = interp.LoadAndChunksSelectMode()
	if err != nil {
		errors = append(errors, err)
	}

	testdata.
		TestInIdentifierInputUploadSecret = "test" +
		randomString()
	testdata.
		TestInValueInputUploadSecret = "test" + randomString()

	err = interp.RunUploadSecret()
	if err != nil {
		errors = append(errors, err)
	}

	err = interp.RunGetSecrets()
	if err != nil {
		errors = append(errors, err)
	}

	testdata.TestInIdentifierInput = testdata.
		TestInIdentifierInputUploadSecret

	err = interp.RunGetSecretByID()
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) != 8 {
		t.Errorf("TestMain->Init: %v", errCntErrors)
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
