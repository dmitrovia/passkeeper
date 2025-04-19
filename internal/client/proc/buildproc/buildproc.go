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

	err := proc.Build()
	if err != nil {
		return fmt.Errorf("RP->Build: %w", err)
	}

	fmt.Println("Successfully build")

	return nil
}

func (proc *BuildProc) Build() error {
	oFile, err := os.Create(proc.attr.OutFilePath)
	if err != nil {
		return fmt.Errorf("RP->Create: %w", err)
	}
	defer oFile.Close()

	chunks := make([]chunckmeta.ChunkMeta, 0,
		len(proc.attr.BuildMetadata))

	for _, chunk := range proc.attr.BuildMetadata {
		chunks = append(chunks, chunk)
	}

	slices.SortFunc(chunks,
		func(a, b chunckmeta.ChunkMeta) int {
			return cmp.Compare(*a.Index, *b.Index)
		})

	for _, chunk := range chunks {
		chunkFile, err := os.Open(*chunk.FileName)
		if err != nil {
			return fmt.Errorf("RP->Open: %w", err)
		}

		chunkFile.Close()

		_, err = io.Copy(oFile, chunkFile)
		if err != nil {
			return fmt.Errorf("RP->Copy: %w", err)
		}

		err = os.Remove(*chunk.FileName)
		if err != nil {
			return fmt.Errorf("RP->Remove: %w", err)
		}
	}

	return nil
}
