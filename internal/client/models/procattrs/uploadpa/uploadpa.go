package uploadpa

import (
	"sync"
	"time"

	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type UploadProcAttr struct {
	CountWorkers int
	ReqTimeout   time.Duration
	ServerURL    string
	Chunks       []chunckmeta.ChunkMeta
	Metadata     map[string]chunckmeta.ChunkMeta
	Wgroup       *sync.WaitGroup
	Mutex        *sync.Mutex
}
