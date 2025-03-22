package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/server/models/procattrs/serverpa"
)

const fmd os.FileMode = 0o666

func GetAttrs(
	attr *serverpa.ServerProcAttr,
) error {
	cfg, err := loadCFGServer(*attr.GetConfigPath())
	if err != nil {
		return fmt.Errorf("GetAttrs->loadConfigServer: %w", err)
	}

	if *attr.GetDBDSN() == "" {
		attr.SetdBDSN(cfg.DBDSN)
	}

	if *attr.GetServerAddr() == "" {
		attr.SetServerAddr(cfg.ServerAddr)
	}

	return nil
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
