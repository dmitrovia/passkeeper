package getsecretsproc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dmitrovia/passkeeper/internal/client/proc/getsecretsproc/getsecretsprocattr"
	"github.com/dmitrovia/passkeeper/internal/general/compress"
	"github.com/dmitrovia/passkeeper/internal/general/models/secret"
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
)

var errSNOK = errors.New("status is not OK")

type GetSecrets struct {
	attr *getsecretsprocattr.GetSecretsProcAttr
}

func NewProc(
	attr *getsecretsprocattr.GetSecretsProcAttr,
) *GetSecrets {
	return &GetSecrets{
		attr: attr,
	}
}

func (lp *GetSecrets) RunProcess() error {
	fmt.Println("GetSecretsProc run")

	defer fmt.Println("GetSecretsProc end")

	ctx, cancel := context.WithTimeout(
		context.Background(), lp.attr.ReqTimeout)

	defer cancel()

	resp, err := lp.attr.GetSecrets.CallEndpoint(ctx)
	if err != nil {
		return fmt.Errorf("RP->CallEndpoint: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("RP->GetSecrets: %w", errSNOK)

		return err
	}

	err = lp.parseResp(resp)
	if err != nil {
		return fmt.Errorf("RP->CallEndpoint: %w", err)
	}

	fmt.Println("Successfully get secrets")

	return nil
}

func (lp *GetSecrets) parseResp(
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
