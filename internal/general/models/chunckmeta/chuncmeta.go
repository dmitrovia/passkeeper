package chunckmeta

import (
	"time"

	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
)

type ChunkMeta struct {
	FileName    *string     `json:"fileName"`
	ID          int32       `json:"id,omitempty"`
	Index       *int        `json:"index,omitempty"`
	Hash        *string     `json:"hash,omitempty"`
	FilePath    *string     `json:"filePath,omitempty"`
	Data        *[]byte     `json:"data,omitempty"`
	User        *userm.User `json:"user,omitempty"`
	Createddate *time.Time  `json:"createdDate,omitempty"`
}

func NewChunkMeta(
	fname string,
	hash string,
	index int,
	data *[]byte,
) *ChunkMeta {
	return &ChunkMeta{
		FileName: &fname,
		Hash:     &hash,
		Index:    &index,
		Data:     data,
	}
}

func (cm *ChunkMeta) ClearData() {
	cm.Data = nil
}
