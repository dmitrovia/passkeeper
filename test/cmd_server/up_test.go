package main_test

import (
	"testing"
	"time"

	"github.com/dmitrovia/passkeeper/internal/server/proc/serverproc"
)

func TestMain(t *testing.T) {
	t.Helper()
	t.Parallel()

	go func() {
		sp := serverproc.NewProc(nil)

		err := sp.RunProcess()
		if err != nil {
			t.Errorf("serverproc.RunProcess %v", err)

			return
		}
	}()

	<-time.After(time.Duration(60) * time.Second)
}
