package interactionpa

import (
	"fmt"
	"sync"

	"github.com/dmitrovia/passkeeper/internal/client/proc/chunkerproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/chunkerproc/chunkerpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/registerproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/registerproc/registerprocpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadproc/uploadpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type InteractionProcAttr struct {
	Chproc        *chunkerproc.ChunkerProc
	Chunkerpa     *chunkerpa.ChunkerProcAttr
	Uploadproc    *uploadproc.UploadProc
	Uploadpa      *uploadpa.UploadProcAttr
	AttrClintProc *clientpa.ClientProcAttr
	Registerproc  *registerproc.RegisterProc
	Registerpa    *registerprocpa.RegisterProcAttr
	UploadChan    chan chunckmeta.ChunkMeta
	ErrChan       chan error
	WGsubprocess  *sync.WaitGroup
}

func (ipa *InteractionProcAttr) InitRegister() error {
	ipa.Registerpa = &registerprocpa.RegisterProcAttr{}

	err := ipa.Registerpa.Init(ipa.AttrClintProc)
	if err != nil {
		return fmt.Errorf("InitRegister->Init: %w", err)
	}

	ipa.Registerpa.Wgroup = ipa.WGsubprocess

	ipa.Registerproc = registerproc.NewProc(ipa.Registerpa)

	return nil
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

	ipa.UploadChan = make(chan chunckmeta.ChunkMeta,
		ipa.Chunkerpa.CntChunks)
	ipa.ErrChan = make(chan error, ipa.Chunkerpa.CntChunks)

	ipa.Chunkerpa.Wgroup = ipa.WGsubprocess
	ipa.Chunkerpa.UploadChan = ipa.UploadChan
	ipa.Chunkerpa.ErrChan = ipa.ErrChan
	ipa.Chproc = chunkerproc.NewProc(ipa.Chunkerpa)
	ipa.Uploadpa.Wgroup = ipa.WGsubprocess
	ipa.Uploadpa.UploadChan = ipa.UploadChan
	ipa.Uploadpa.ErrChan = ipa.ErrChan
	ipa.Uploadproc = uploadproc.NewProc(ipa.Uploadpa)

	return nil
}
