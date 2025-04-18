package logoutproc

import (
	"fmt"

	"github.com/dmitrovia/passkeeper/internal/client/auth/authcfg"
	logoutattr "github.com/dmitrovia/passkeeper/internal/client/proc/logoutproc/logoutprocattr"
)

type LogoutProc struct {
	attr *logoutattr.LogoutProcAttr
}

func NewProc(
	attr *logoutattr.LogoutProcAttr,
) *LogoutProc {
	return &LogoutProc{
		attr: attr,
	}
}

func (lp *LogoutProc) RunProcess() error {
	defer lp.attr.WgSubProc.Done()
	fmt.Println("LogoutProc run")

	defer fmt.Println("LogoutProc end")

	lp.attr.AttrClintProc.SetAuth("")

	err := authcfg.SaveToken(
		lp.attr.AttrClintProc.AuthTokenPath, "")
	if err != nil {
		return fmt.Errorf("RP->SaveToken: %w", err)
	}

	fmt.Println("Successfully logout")

	return nil
}
