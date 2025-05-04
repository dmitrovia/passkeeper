package uploadsecretproc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadsecretproc/uploadsecretprocattr"
	"github.com/dmitrovia/passkeeper/internal/general/models/secret"
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
)

var errSNOK = errors.New("status is not OK")

type UploadSecretProc struct {
	attr *uploadsecretprocattr.UploadSecretProcAttr
}

func NewProc(
	attr *uploadsecretprocattr.UploadSecretProcAttr,
) *UploadSecretProc {
	return &UploadSecretProc{
		attr: attr,
	}
}

func (lp *UploadSecretProc) RunProcess() error {
	fmt.Println("UploadSecretProc run")

	defer fmt.Println("UploadSecretProc end")

	ctx, cancel := context.WithTimeout(
		context.Background(), lp.attr.ReqTimeout)

	defer cancel()

	err := lp.input()
	if err != nil {
		return fmt.Errorf("RP->Input: %w", err)
	}

	resp, err := lp.attr.UploadSecret.CallEndpoint(ctx)
	if err != nil {
		return fmt.Errorf("RP->CallEndpoint: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("RP->UploadSecret: %w", errSNOK)
		return err
	}

	fmt.Println("Successfully upload")

	return nil
}

func (lp *UploadSecretProc) input() error {
	var inIdentifier string

	var inValue string

	fmt.Println("Enter identifier of secret")

	td := lp.attr.AttrClintProc.TestData != nil &&
		lp.attr.AttrClintProc.TestData.
			TestInIdentifierInputUploadSecret != ""

	if !td {
		_, err := fmt.Fscan(os.Stdin, &inIdentifier)
		if err != nil {
			return fmt.Errorf("Input->Fscan(identifier): %w", err)
		}
	} else {
		inIdentifier = lp.attr.AttrClintProc.TestData.
			TestInIdentifierInputUploadSecret
	}

	fmt.Println("Enter value")

	td1 := lp.attr.AttrClintProc.TestData != nil &&
		lp.attr.AttrClintProc.TestData.
			TestInValueInputUploadSecret != ""

	if !td1 {
		_, err := fmt.Fscan(os.Stdin, &inValue)
		if err != nil {
			return fmt.Errorf("Input->Fscan(value secret): %w", err)
		}
	} else {
		inValue = lp.attr.AttrClintProc.TestData.
			TestInValueInputUploadSecret
	}

	secret := &secret.Secret{}
	secret.Identifier = &inIdentifier
	secret.Value = &inValue

	marshal, err := json.Marshal(secret)
	if err != nil {
		return fmt.Errorf("Input->Marshal: %w", err)
	}

	encrypt, err := rsa.Encrypt(&marshal, lp.attr.EncKey)
	if err != nil {
		return fmt.Errorf("Input->Encrypt: %w", err)
	}

	lp.attr.UploadSecretAttr.Data = encrypt

	return nil
}
