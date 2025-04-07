package logoutprocattr

import (
	"sync"

	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
)

type LogoutProcAttr struct {
	Wgroup        *sync.WaitGroup
	AttrClintProc *clientpa.ClientProcAttr
}

func (lpa *LogoutProcAttr) Init(
	attr *clientpa.ClientProcAttr,
) error {
	lpa.AttrClintProc = attr
	lpa.Wgroup = attr.WGsubprocess

	return nil
}
