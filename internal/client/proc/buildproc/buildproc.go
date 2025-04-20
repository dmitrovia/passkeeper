package buildproc

import (
	"cmp"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/dmitrovia/passkeeper/internal/client/proc/buildproc/buildprocattr"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type BuildProc struct {
	attr *buildprocattr.BuildProcAttr
}

func NewProc(
	attr *buildprocattr.BuildProcAttr,
) *BuildProc {
	return &BuildProc{
		attr: attr,
	}
}

func (proc *BuildProc) RunProcess() error {
	fmt.Println("BuildProc run")

	defer fmt.Println("BuildProc end")

	err := proc.build()
	if err != nil {
		return fmt.Errorf("RP->Build: %w", err)
	}

	fmt.Println("Successfully build")

	return nil
}

func (proc *BuildProc) build() error {
	oFile, err := os.Create(proc.attr.OutFilePath)
	if err != nil {
		return fmt.Errorf("RP->Create: %w", err)
	}
	defer oFile.Close()

	chunks := make([]chunckmeta.ChunkMeta, 0,
		len(proc.attr.BuildMetadata))

	for _, chunk := range proc.attr.BuildMetadata {
		chunks = append(chunks, *chunk)
	}

	slices.SortFunc(chunks,
		func(a, b chunckmeta.ChunkMeta) int {
			return cmp.Compare(*a.Index, *b.Index)
		})

	for _, chunk := range chunks {
		chunkFile, err := os.Open(*chunk.FilePath)
		if err != nil {
			return fmt.Errorf("RP->Open: %w", err)
		}

		_, err = io.Copy(oFile, chunkFile)
		if err != nil {
			return fmt.Errorf("RP->Copy: %w", err)
		}

		chunkFile.Close()

		err = os.Remove(*chunk.FilePath)
		if err != nil {
			return fmt.Errorf("RP->Remove: %w", err)
		}

		chunk.ClearAllExceptMeta()
		proc.attr.CurrentMetadata[*chunk.FileName] = chunk
	}

	return nil
}
