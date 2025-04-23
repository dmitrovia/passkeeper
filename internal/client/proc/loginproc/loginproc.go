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
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
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
	fmt.Println("LoginProc run")

	defer fmt.Println("LoginProc end")

	ctx, cancel := context.WithTimeout(
		context.Background(), lp.attr.ReqTimeout)

	defer cancel()

	err := lp.input()
	if err != nil {
		return fmt.Errorf("RP->Input: %w", err)
	}

	resp, err := lp.attr.Login.CallEndpoint(ctx)
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

func (lp *LoginProc) input() error {
	var inLogin string

	var inPass string

	fmt.Println("Enter login")

	_, err := fmt.Fscan(os.Stdin, &inLogin)
	if err != nil {
		return fmt.Errorf("Input->Fscan(login): %w", err)
	}

	fmt.Println("Enter password")

	_, err = fmt.Fscan(os.Stdin, &inPass)
	if err != nil {
		return fmt.Errorf("Input->Fscan(pass): %w", err)
	}

	data := apim.InLoginUser{}
	data.Login = inLogin
	data.Password = inPass

	marshal, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Input->Marshal: %w", err)
	}

	encrypt, err := rsa.Encrypt(&marshal, lp.attr.EncKey)
	if err != nil {
		return fmt.Errorf("Input->Encrypt: %w", err)
	}

	lp.attr.LoginAttr.Data = encrypt

	return nil
}
