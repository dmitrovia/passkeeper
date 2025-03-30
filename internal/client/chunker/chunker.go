package chunker

import "github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"

type FileChunker struct {
	chunkSize int
	filesPath string
}

func NewFileChunker(
	chunkSize int,
	filesPath string,
) *FileChunker {
	return &FileChunker{
		chunkSize: chunkSize,
		filesPath: filesPath,
	}
}

func (c *FileChunker) ChunkFiles() (
	[]chunckmeta.ChunkMeta, error,
) {
	return nil, nil
}

func (c *FileChunker) ChunklargeFiles() (
	[]chunckmeta.ChunkMeta, error,
) {
	return nil, nil
}
