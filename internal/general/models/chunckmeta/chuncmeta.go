package chunckmeta

import (
	"time"

	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
)

type ChunkMeta struct {
	FileName    *string     `json:"fileName"`
	Createddate *time.Time  `json:"createdDate,omitempty"`
	ID          int32       `json:"id,omitempty"`
	Hash        *string     `json:"hash"`
	Index       *int        `json:"index"`
	Data        *[]byte     `json:"data,omitempty"`
	User        *userm.User `json:"user,omitempty"`
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
