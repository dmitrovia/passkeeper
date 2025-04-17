package interactionpa

import (
	"fmt"
	"sync"

	"github.com/dmitrovia/passkeeper/internal/client/metamanager"
	"github.com/dmitrovia/passkeeper/internal/client/proc/chunkerproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/chunkerproc/chunkerpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/initsingleloadproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/initsingleloadproc/initsingleloadprocattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/inituploadproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/inituploadproc/inituploadprocattr.go"
	"github.com/dmitrovia/passkeeper/internal/client/proc/loginproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/loginproc/loginprocattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/logoutproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/logoutproc/logoutprocattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/registerproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/registerproc/registerprocpattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadproc/uploadpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type InteractionProcAttr struct {
	// upload and chunk
	Chproc     *chunkerproc.ChunkerProc
	Chunkerpa  *chunkerpa.ChunkerProcAttr
	Uploadproc *uploadproc.UploadProc
	Uploadpa   *uploadpa.UploadProcAttr
	UploadChan chan chunckmeta.ChunkMeta
	ErrChan    chan error
	// init singleinitload
	InitSingleLoadproc *initsingleloadproc.InitSingleProc
	InitSingleLoadpa   *initsingleloadprocattr.
				InitUploadProcAttr
	// init upload
	InitUploadproc *inituploadproc.InitUploadProc
	InitUploadpa   *inituploadprocattr.InitUploadProcAttr
	// register
	Registerproc *registerproc.RegisterProc
	Registerpa   *registerprocpattr.RegisterProcAttr
	// login
	Loginproc *loginproc.LoginProc
	Loginpa   *loginprocattr.LoginProcAttr
	// logout
	Logoutproc *logoutproc.LogoutProc
	Logoutpa   *logoutprocattr.LogoutProcAttr
	// general
	AttrClintProc      *clientpa.ClientProcAttr
	WGsubprocess       *sync.WaitGroup
	WorkerChunkWg      *sync.WaitGroup
	Metamanager        *metamanager.MetaManager
	CurrentMetadata    map[string]chunckmeta.ChunkMeta
	LoadMetadata       map[string]chunckmeta.ChunkMeta
	SpecificFileUpload bool
	SpecificFileLoad   bool
}

func (ipa *InteractionProcAttr) InitRegister() error {
	ipa.Registerpa = &registerprocpattr.RegisterProcAttr{}

	err := ipa.Registerpa.Init(ipa.AttrClintProc)
	if err != nil {
		return fmt.Errorf("InitRegister->Init: %w", err)
	}

	ipa.Registerproc = registerproc.NewProc(ipa.Registerpa)

	return nil
}

func (ipa *InteractionProcAttr) InitLogin() error {
	ipa.Loginpa = &loginprocattr.LoginProcAttr{}

	err := ipa.Loginpa.Init(ipa.AttrClintProc)
	if err != nil {
		return fmt.Errorf("InitLogin->Init: %w", err)
	}

	ipa.Loginproc = loginproc.NewProc(ipa.Loginpa)

	return nil
}

func (ipa *InteractionProcAttr) InitLogout() error {
	ipa.Logoutpa = &logoutprocattr.LogoutProcAttr{}

	err := ipa.Logoutpa.Init(ipa.AttrClintProc)
	if err != nil {
		return fmt.Errorf("InitLogout->Init: %w", err)
	}

	ipa.Logoutproc = logoutproc.NewProc(ipa.Logoutpa)

	return nil
}

func (ipa *InteractionProcAttr) InitChunkAndUpload() error {
	ipa.Chunkerpa = &chunkerpa.ChunkerProcAttr{}
	ipa.Metamanager = metamanager.NewMetaManager(
		ipa.AttrClintProc.MetaPath)
	ipa.Chunkerpa.Metamanager = ipa.Metamanager

	metadata, err := ipa.Metamanager.LoadMetadata()
	if err != nil {
		return fmt.Errorf("InitChunkAndUpload->LM: %w", err)
	}

	ipa.CurrentMetadata = metadata

	err = ipa.Chunkerpa.Init(ipa.AttrClintProc)
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

	ipa.WorkerChunkWg = &sync.WaitGroup{}

	ipa.Chunkerpa.UploadChan = ipa.UploadChan
	ipa.Chunkerpa.ErrChan = ipa.ErrChan
	ipa.Chunkerpa.WorkerChunkWg = ipa.WorkerChunkWg
	ipa.Chunkerpa.CurrentMetadata = metadata
	ipa.Chproc = chunkerproc.NewProc(ipa.Chunkerpa)
	ipa.Uploadpa.UploadChan = ipa.UploadChan
	ipa.Uploadpa.CurrentMetadata = metadata
	ipa.Uploadpa.ErrChan = ipa.ErrChan
	ipa.Uploadpa.CountChunk = ipa.Chunkerpa.CntChunks
	ipa.Uploadpa.WorkerChunkWg = ipa.WorkerChunkWg
	ipa.Uploadproc = uploadproc.NewProc(ipa.Uploadpa)

	ipa.InitUploadpa = &inituploadprocattr.InitUploadProcAttr{}
	ipa.InitUploadpa.Init(ipa.AttrClintProc)
	ipa.InitUploadproc = inituploadproc.NewProc(
		ipa.InitUploadpa)

	ipa.InitSingleLoadpa = &initsingleloadprocattr.
		InitUploadProcAttr{}
	ipa.InitSingleLoadpa.LoadMetadata = ipa.LoadMetadata
	ipa.InitSingleLoadpa.Init(ipa.AttrClintProc)
	ipa.InitSingleLoadproc = initsingleloadproc.NewProc(
		ipa.InitSingleLoadpa)

	return nil
}
