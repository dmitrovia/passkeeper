package chunkerpa

import "github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"

type ChunkerProcAttr struct {
	CountWorkers    int
	ChunkSize       int
	FilePath        string
	CurrentMetadata map[string]chunckmeta.ChunkMeta
}
