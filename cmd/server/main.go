package main

import (
	"github.com/dmitrovia/passkeeper/internal/general/functions/loggerf"
	"github.com/dmitrovia/passkeeper/internal/server/proc/serverproc"
)

func main() {
	sp := serverproc.NewProc(nil)

	err := sp.RunProcess()
	if err != nil {
		loggerf.Log("serverproc.RunProcess", err)

		return
	}
}
