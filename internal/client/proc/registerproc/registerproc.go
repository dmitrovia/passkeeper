package registerproc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/dmitrovia/passkeeper/internal/client/proc/registerproc/registerprocpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
)

var errSNOK = errors.New("status is not OK")

type RegisterProc struct {
	attr *registerprocpa.RegisterProcAttr
}

func NewProc(
	attr *registerprocpa.RegisterProcAttr,
) *RegisterProc {
	return &RegisterProc{
		attr: attr,
	}
}

func (rp *RegisterProc) RunProcess() error {
	defer rp.attr.Wgroup.Done()
	fmt.Println("RegisterProc run")

	defer fmt.Println("RegisterProc end")

	ctx, cancel := context.WithTimeout(
		context.Background(), rp.attr.ReqTimeout)
	defer cancel()

	err := rp.Input()
	if err != nil {
		return fmt.Errorf("RP->Input: %w", err)
	}

	resp, err := rp.attr.Register.RegisterUser(ctx)
	if err != nil {
		return fmt.Errorf("RP->RegisterUser: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("RP->Register: %w", errSNOK)

		return err
	}

	fmt.Println("Successfully registered")

	return nil
}

func (rp *RegisterProc) Input() error {
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

	rp.attr.RegisterAttr.Data = &marshal

	return nil
}
