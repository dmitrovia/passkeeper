package interactionproc

import (
	"fmt"

	"github.com/dmitrovia/passkeeper/internal/client/metamanager"
	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/chunkerpa"
	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/interactionpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/chunkerproc"
)

type InteractionProc struct {
	attr          *interactionpa.InteractionProcAttr
	attrClintProc *clientpa.ClientProcAttr
}

func NewProc(attrClintProc *clientpa.ClientProcAttr,
	attr *interactionpa.InteractionProcAttr,
) *InteractionProc {
	return &InteractionProc{
		attrClintProc: attrClintProc,
		attr:          attr,
	}
}

func (ip *InteractionProc) RunProcess() error {
	fmt.Println("InteractionProc run")
	defer fmt.Println("InteractionProc end")

	if ip.attr == nil {
		ip.attr = &interactionpa.InteractionProcAttr{}
	}

	metaManager := metamanager.NewMetaManager(
		ip.attrClintProc.MetaPath)

	metadata, err := metaManager.LoadMetadata()
	if err != nil {
		return fmt.Errorf("RP->LoadMetadata: %w", err)
	}

	chpa := &chunkerpa.ChunkerProcAttr{}
	chpa.ChunkSize = ip.attrClintProc.DefChunkSize
	chpa.FilePath = ip.attrClintProc.FileSynchronizePath
	chpa.CurrentMetadata = metadata

	chproc := chunkerproc.NewProc(chpa)

	err = chproc.RunProcess()
	if err != nil {
		return fmt.Errorf("RP->chproc.RP: %w", err)
	}

	return nil
}
