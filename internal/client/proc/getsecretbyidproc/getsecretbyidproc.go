package getsecretbyidproc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dmitrovia/passkeeper/internal/client/proc/getsecretbyidproc/getsecretbyidprocattr"
	"github.com/dmitrovia/passkeeper/internal/general/compress"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/general/models/secret"
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
)

var errSNOK = errors.New("status is not OK")

type GetSecretByID struct {
	attr *getsecretbyidprocattr.GetSecretByIDProcAttr
}

func NewProc(
	attr *getsecretbyidprocattr.GetSecretByIDProcAttr,
) *GetSecretByID {
	return &GetSecretByID{
		attr: attr,
	}
}

func (lp *GetSecretByID) RunProcess() error {
	fmt.Println("GetSecretByIDProc run")

	defer fmt.Println("GetSecretByIDProc end")

	ctx, cancel := context.WithTimeout(
		context.Background(), lp.attr.ReqTimeout)

	defer cancel()

	err := lp.input()
	if err != nil {
		return fmt.Errorf("RP->Input: %w", err)
	}

	resp, err := lp.attr.GetSecretByID.CallEndpoint(ctx)
	if err != nil {
		return fmt.Errorf("RP->CallEndpoint: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("RP->GetSecretByID: %w", errSNOK)
		return err
	}

	err = lp.parseResp(resp)
	if err != nil {
		return fmt.Errorf("RP->CallEndpoint: %w", err)
	}

	fmt.Println("Successfully get secret")

	return nil
}

func (lp *GetSecretByID) parseResp(
	resp *http.Response,
) error {
	secrets := []secret.Secret{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("PR->io.ReadAll: %w", err)
	}

	dec, err := rsa.Decrypt(&body, lp.attr.DecKey)
	if err != nil {
		return fmt.Errorf("PR->Decrypt: %w", err)
	}

	decompress, err := compress.DeflateDecompress(
		bytes.NewReader(*dec),
	)
	if err != nil {
		return fmt.Errorf("PR->DeflateDecompress: %w", err)
	}

	err = json.Unmarshal(decompress, &secrets)
	if err != nil {
		return fmt.Errorf("PR->Unmarshal: %w", err)
	}

	for _, secret := range secrets {
		fmt.Println(*secret.Identifier)
		fmt.Println(*secret.Value)
		fmt.Println("-------------------------------------------")
	}

	return nil
}

func (lp *GetSecretByID) input() error {
	var inIdentifier string

	fmt.Println("Enter identifier of secret")

	if lp.attr.AttrClintProc.TestData == nil {
		_, err := fmt.Fscan(os.Stdin, &inIdentifier)
		if err != nil {
			return fmt.Errorf("Input->Fscan(identifier): %w", err)
		}
	} else {
		inIdentifier = lp.attr.AttrClintProc.TestData.
			TestInIdentifierInput
	}

	outAttr := &apim.InGetSecretByID{}
	outAttr.Identifier = inIdentifier

	marshal, err := json.Marshal(outAttr)
	if err != nil {
		return fmt.Errorf("Input->Marshal: %w", err)
	}

	encrypt, err := rsa.Encrypt(&marshal, lp.attr.EncKey)
	if err != nil {
		return fmt.Errorf("Input->Encrypt: %w", err)
	}

	lp.attr.GetSecretByIDAttr.Data = encrypt

	return nil
}
