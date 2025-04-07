package loginproc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/dmitrovia/passkeeper/internal/client/auth/authcfg"
	"github.com/dmitrovia/passkeeper/internal/client/proc/loginproc/loginprocattr"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
)

var errSNOK = errors.New("status is not OK")

type LoginProc struct {
	attr *loginprocattr.LoginProcAttr
}

func NewProc(
	attr *loginprocattr.LoginProcAttr,
) *LoginProc {
	return &LoginProc{
		attr: attr,
	}
}

func (lp *LoginProc) RunProcess() error {
	defer lp.attr.Wgroup.Done()
	fmt.Println("LoginProc run")

	defer fmt.Println("LoginProc end")

	ctx, cancel := context.WithTimeout(
		context.Background(), lp.attr.ReqTimeout)

	defer cancel()

	err := lp.Input()
	if err != nil {
		return fmt.Errorf("RP->Input: %w", err)
	}

	resp, err := lp.attr.Login.LoginUser(ctx)
	if err != nil {
		return fmt.Errorf("RP->LoginUser: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("RP->Login: %w", errSNOK)

		return err
	}

	token := resp.Header.Get("Authorization")

	err = authcfg.SaveToken(lp.attr.TokenSavePath, token)
	if err != nil {
		return fmt.Errorf("RP->SaveToken: %w", err)
	}

	lp.attr.AttrClintProc.SetAuth(token)

	fmt.Println("Successfully login")

	return nil
}

func (lp *LoginProc) Input() error {
	var inLogin string

	var inPass string

	fmt.Println("Enter login")

	_, err := fmt.Fscan(os.Stdin, &inLogin)
	if err != nil {
		return fmt.Errorf("RP->Fscan(login): %w", err)
	}

	fmt.Println("Enter password")

	_, err = fmt.Fscan(os.Stdin, &inPass)
	if err != nil {
		return fmt.Errorf("RP->Fscan(pass): %w", err)
	}

	data := apim.InLoginUser{}
	data.Login = inLogin
	data.Password = inPass

	marshal, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("RP->Marshal: %w", err)
	}

	lp.attr.LoginAttr.Data = &marshal

	return nil
}
