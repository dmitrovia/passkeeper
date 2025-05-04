package logoutprocattr

import (
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
)

type LogoutProcAttr struct {
	AttrClintProc *clientpa.ClientProcAttr
}

func (lpa *LogoutProcAttr) Init(
	attr *clientpa.ClientProcAttr,
) error {
	lpa.AttrClintProc = attr

	return nil
}
