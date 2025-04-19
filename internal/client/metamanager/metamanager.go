package metamanager

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type MetaManager struct {
	metaPath string
}

func NewMetaManager(
	metaPath string,
) *MetaManager {
	return &MetaManager{
		metaPath: metaPath,
	}
}

const fmd os.FileMode = 0o666

func (m *MetaManager) LoadMetadata() (
	map[string]chunckmeta.ChunkMeta, error,
) {
	metadata := make(map[string]chunckmeta.ChunkMeta)

	data, err := os.ReadFile(m.metaPath)
	if err != nil {
		return metadata, fmt.Errorf("LoadMetadata->RF: %w", err)
	}

	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return metadata, err
	}

	return metadata, nil
}

func (m *MetaManager) SaveMetadata(
	metadata map[string]chunckmeta.ChunkMeta,
) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(m.metaPath, data, fmd)
	if err != nil {
		return fmt.Errorf("SaveMetadata->WF: %w", err)
	}

	return nil
}
