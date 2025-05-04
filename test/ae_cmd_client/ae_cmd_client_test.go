package main_test

import (
	"testing"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc"
)

func TestMain(t *testing.T) {
	t.Helper()
	t.Parallel()

	go func() {
		cp := clientproc.NewProc(nil)

		err := cp.RunProcess()
		if err != nil {
			t.Errorf("clientproc.RunProcess %v", err)

			return
		}
	}()

	<-time.After(time.Duration(10) * time.Second)
}
