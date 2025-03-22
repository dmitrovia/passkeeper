// Main agent application package.
package main

import (
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc"
	"github.com/dmitrovia/passkeeper/internal/general/functions/loggerf"
)

func main() {
	err := clientproc.RunProcess()
	if err != nil {
		loggerf.Log("clientproc.RunProcess", err)

		return
	}
}
