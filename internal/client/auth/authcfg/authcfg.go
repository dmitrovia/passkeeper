package authcfg

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
)

const fmd os.FileMode = 0o666

func GetToken(path string) (string, error) {
	cfg, err := loadToken(path)
	if err != nil {
		return "", fmt.Errorf("GetToken->loadToken: %w", err)
	}

	return cfg.Token, nil
}

func SaveToken(path string, token string) error {
	err := uploadToken(path, token)
	if err != nil {
		return fmt.Errorf("SaveToken->loadToken: %w", err)
	}

	return nil
}

func uploadToken(
	path string,
	token string,
) error {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_EXCL, fmd)
	if err != nil {
		return fmt.Errorf("uploadToken->OF: %w", err)
	}

	defer file.Close()

	cfg := apim.CfgToken{}
	cfg.Token = token

	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("uploadToken->Marshal: %w", err)
	}

	err = os.WriteFile(path, data, fmd)
	if err != nil {
		return fmt.Errorf("uploadToken->WriteFile: %w", err)
	}

	return nil
}

func loadToken(
	pth string,
) (*apim.CfgToken, error) {
	file, err := os.OpenFile(pth, os.O_RDONLY|os.O_EXCL, fmd)
	if err != nil {
		return nil, fmt.Errorf("loadToken->OF: %w", err)
	}

	defer file.Close()

	attr := &apim.CfgToken{}

	bytes, _ := io.ReadAll(file)

	err = json.Unmarshal(bytes, &attr)
	if err != nil {
		return nil, fmt.Errorf("loadToken->Unmarshal: %w", err)
	}

	return attr, nil
}
