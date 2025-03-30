package interaction

import (
	"fmt"
	"sync"

	"github.com/dmitrovia/passkeeper/internal/client/chunker"
	"github.com/dmitrovia/passkeeper/internal/client/metamanager"
	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/uploadpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadproc"
)

func RunProcess(clientAttr *clientpa.ClientProcAttr) error {
	fmt.Println("InteractionProc run")
	defer fmt.Println("InteractionProc end")
	fmt.Println(clientAttr)

	attr := &uploadpa.UploadProcAttr{}
	chunker := chunker.NewFileChunker(clientAttr.DefChunkSize,
		clientAttr.FileSynchronizePath)

	metaManager := metamanager.NewMetaManager(
		clientAttr.MetaPath)

	metadata, err := metaManager.LoadMetadata()
	if err != nil {
		return fmt.Errorf("RP->LoadMetadata: %w", err)
	}

	fmt.Println(metadata)

	chunks, err := chunker.ChunkFiles()
	if err != nil {
		return fmt.Errorf("RP->ChunkFiles: %w", err)
	}

	fmt.Println(chunks)

	err = uploadproc.RunProcess(attr)
	if err != nil {
		return fmt.Errorf("RP->uploadproc.RP: %w", err)
	}

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	fmt.Println(wg)
	fmt.Println(mu)

	return nil
}
