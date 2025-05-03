package main

import (
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc"
	"github.com/dmitrovia/passkeeper/internal/general/functions/loggerf"
)

func main() {
	cp := clientproc.NewProc(nil)

	err := cp.RunProcess()
	if err != nil {
		loggerf.Log("clientproc.RunProcess", err)
		return
	}
}
