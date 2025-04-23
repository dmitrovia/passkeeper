package registerproc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/dmitrovia/passkeeper/internal/client/proc/registerproc/registerprocpattr"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
)

var errSNOK = errors.New("status is not OK")

type RegisterProc struct {
	attr *registerprocpattr.RegisterProcAttr
}

func NewProc(
	attr *registerprocpattr.RegisterProcAttr,
) *RegisterProc {
	return &RegisterProc{
		attr: attr,
	}
}

func (rp *RegisterProc) RunProcess() error {
	fmt.Println("RegisterProc run")

	defer fmt.Println("RegisterProc end")

	ctx, cancel := context.WithTimeout(
		context.Background(), rp.attr.ReqTimeout)
	defer cancel()

	err := rp.input()
	if err != nil {
		return fmt.Errorf("RP->Input: %w", err)
	}

	resp, err := rp.attr.Register.CallEndpoint(ctx)
	if err != nil {
		return fmt.Errorf("RP->RegisterUser: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("RP->Register: %w", errSNOK)
	}

	fmt.Println("Successfully registered")

	return nil
}

func (rp *RegisterProc) input() error {
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

	encrypt, err := rsa.Encrypt(&marshal, rp.attr.EncKey)
	if err != nil {
		return fmt.Errorf("Input->Encrypt: %w", err)
	}

	rp.attr.RegisterAttr.Data = encrypt

	return nil
}
