package buildprocattr

import (
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type BuildProcAttr struct {
	OutFilePath   string
	BuildMetadata map[string]chunckmeta.ChunkMeta
}
