package interactionpa

import (
	"fmt"
	"sync"

	"github.com/dmitrovia/passkeeper/internal/client/proc/chunkerproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/chunkerproc/chunkerpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadproc/uploadpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type InteractionProcAttr struct {
	Uploadpa      *uploadpa.UploadProcAttr
	Chunkerpa     *chunkerpa.ChunkerProcAttr
	Chproc        *chunkerproc.ChunkerProc
	Uploadproc    *uploadproc.UploadProc
	Wgroup        *sync.WaitGroup
	AttrClintProc *clientpa.ClientProcAttr
	UploadChan    chan chunckmeta.ChunkMeta
	ErrChan       chan error
}

func (ipa *InteractionProcAttr) InitChunkAndUpload() error {
	ipa.Chunkerpa = &chunkerpa.ChunkerProcAttr{}

	err := ipa.Chunkerpa.Init(ipa.AttrClintProc)
	if err != nil {
		return fmt.Errorf("InitChunkAndUpload->Init: %w", err)
	}

	ipa.Uploadpa = &uploadpa.UploadProcAttr{}

	err = ipa.Uploadpa.Init(ipa.AttrClintProc)
	if err != nil {
		return fmt.Errorf("InitChunkAndUpload->Init: %w", err)
	}

	ipa.Wgroup = &sync.WaitGroup{}
	ipa.UploadChan = make(chan chunckmeta.ChunkMeta,
		ipa.Chunkerpa.CntChunks)
	ipa.ErrChan = make(chan error, ipa.Chunkerpa.CntChunks)

	ipa.Chunkerpa.Wgroup = ipa.Wgroup
	ipa.Chunkerpa.UploadChan = ipa.UploadChan
	ipa.Chunkerpa.ErrChan = ipa.ErrChan
	ipa.Chproc = chunkerproc.NewProc(ipa.Chunkerpa)
	ipa.Uploadpa.Wgroup = ipa.Wgroup
	ipa.Uploadpa.UploadChan = ipa.UploadChan
	ipa.Uploadpa.ErrChan = ipa.ErrChan
	ipa.Uploadproc = uploadproc.NewProc(ipa.Uploadpa)

	return nil
}
