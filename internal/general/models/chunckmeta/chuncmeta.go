package chunckmeta

import "time"

type ChunkMeta struct {
	Createddate *time.Time
	ID          int32
	FileName    *string `json:"fileName"`
	Hash        *string `json:"hash"`
	Index       *int    `json:"index"`
}
