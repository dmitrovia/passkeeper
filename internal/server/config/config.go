package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
)

const fmd os.FileMode = 0o666

func GetAttrs(path string) (*apim.CfgServer, error) {
	cfg, err := loadCFGServer(path)
	if err != nil {
		return nil, fmt.Errorf("GetAttrs->LCS: %w", err)
	}

	return cfg, nil
}

func loadCFGServer(
	pth string,
) (*apim.CfgServer, error) {
	file, err := os.OpenFile(pth, os.O_RDONLY|os.O_EXCL, fmd)
	if err != nil {
		return nil, fmt.Errorf("LoadConfigServer->OF: %w", err)
	}

	defer file.Close()

	attr := &apim.CfgServer{}

	byteValue, _ := io.ReadAll(file)

	err = json.Unmarshal(byteValue, &attr)
	if err != nil {
		return nil, fmt.Errorf("loadCFGServer->Unmarsha: %w", err)
	}

	return attr, nil
}
