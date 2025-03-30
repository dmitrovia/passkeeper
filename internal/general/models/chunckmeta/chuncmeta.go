package chunckmeta

import "time"

type ChunkMeta struct {
	FileName    *string    `json:"fileName"`
	Createddate *time.Time `json:"createdDate"`
	ID          int32      `json:"id"`
	Hash        *string    `json:"hash"`
	Index       *int       `json:"index"`
}
