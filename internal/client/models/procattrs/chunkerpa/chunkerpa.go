package chunkerpa

import (
	"sync"
	"time"

	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type ChunkerProcAttr struct {
	CountWorkersUpload  int
	CountWorkersChunker int
	ChunkSize           int
	FilePath            string
	ServerURL           string
	CurrentMetadata     map[string]chunckmeta.ChunkMeta
	ReqTimeout          time.Duration
	Wgroup              *sync.WaitGroup
	Mutex               *sync.Mutex
}
