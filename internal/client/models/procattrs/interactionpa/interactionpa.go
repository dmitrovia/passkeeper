package interactionpa

import (
	"fmt"
	"sync"

	"github.com/dmitrovia/passkeeper/internal/client/metamanager"
	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/chunkerpa"
	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/uploadpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/chunkerproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadproc"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

const wgCount int = 2

type InteractionProcAttr struct {
	Uploadpa      *uploadpa.UploadProcAttr
	Chunkerpa     *chunkerpa.ChunkerProcAttr
	Chproc        *chunkerproc.ChunkerProc
	Uploadproc    *uploadproc.UploadProc
	Wgroup        *sync.WaitGroup
	AttrClintProc *clientpa.ClientProcAttr
}

func (ipa *InteractionProcAttr) InitChunkAndUpload() error {
	ipa.Wgroup = &sync.WaitGroup{}

	chpa := &chunkerpa.ChunkerProcAttr{}
	chpa.ChunkSize = ipa.AttrClintProc.DefChunkSize
	chpa.FilePath = ipa.AttrClintProc.FileSynchronizePath

	err := chpa.Init()
	if err != nil {
		return fmt.Errorf("InitChunkAndUpload->Init: %w", err)
	}

	uploadChan := make(chan chunckmeta.ChunkMeta,
		chpa.CntChunks)
	errChan := make(chan error, chpa.CntChunks)

	ipa.Wgroup.Add(chpa.CntChunks * wgCount)

	chpa.CountWorkersChunker = ipa.
		AttrClintProc.CountWorkersChunker
	chpa.Wgroup = ipa.Wgroup
	chpa.UploadChan = uploadChan
	chpa.ErrChan = errChan

	ipa.Chproc = chunkerproc.NewProc(chpa)

	upa := &uploadpa.UploadProcAttr{}

	upa.Init()

	metaManager := metamanager.NewMetaManager(
		ipa.AttrClintProc.MetaPath)

	metadata, err := metaManager.LoadMetadata()
	if err != nil {
		return fmt.Errorf("InitChunkAndUpload->LM: %w", err)
	}

	upa.CurrentMetadata = metadata
	upa.CountWorkersUpload = ipa.
		AttrClintProc.CountWorkersUpload
	upa.Wgroup = ipa.Wgroup
	upa.UploadChan = uploadChan
	upa.ErrChan = errChan
	upa.ServerURL = ipa.AttrClintProc.ServerAddr
	upa.Mutex = &sync.Mutex{}

	ipa.Uploadproc = uploadproc.NewProc(upa)

	return nil
}
