package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
)

const fmd os.FileMode = 0o666

func GetAttrsS(path string) (*apim.CfgServer, error) {
	cfg, err := loadCFGServerS(path)
	if err != nil {
		return nil, fmt.Errorf("GetAttrs->LCS: %w", err)
	}

	return cfg, nil
}

func loadCFGServerS(
	pth string,
) (*apim.CfgServer, error) {
	file, err := os.OpenFile(pth, os.O_RDONLY|os.O_EXCL, fmd)
	if err != nil {
		return nil, fmt.Errorf("loadCFGServerS->OF: %w", err)
	}

	defer file.Close()

	attr := &apim.CfgServer{}

	byteValue, _ := io.ReadAll(file)

	err = json.Unmarshal(byteValue, &attr)
	if err != nil {
		return nil, fmt.Errorf("loadCFGServerS->Unmarsh: %w", err)
	}

	return attr, nil
}

func GetAttrsC(path string) (*apim.CfgClient, error) {
	cfg, err := loadCFGServerC(path)
	if err != nil {
		return nil, fmt.Errorf("GetAttrs->LCS: %w", err)
	}

	return cfg, nil
}

func loadCFGServerC(
	pth string,
) (*apim.CfgClient, error) {
	file, err := os.OpenFile(pth, os.O_RDONLY|os.O_EXCL, fmd)
	if err != nil {
		return nil, fmt.Errorf("loadCFGServerC->OF: %w", err)
	}

	defer file.Close()

	attr := &apim.CfgClient{}

	byteValue, _ := io.ReadAll(file)

	err = json.Unmarshal(byteValue, &attr)
	if err != nil {
		return nil, fmt.Errorf("loadCFGServerC->Unmarsh: %w", err)
	}

	return attr, nil
}
