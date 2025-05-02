package interactionpa

import (
	"fmt"
	"os"
	"sync"

	"github.com/dmitrovia/passkeeper/internal/client/metamanager"
	"github.com/dmitrovia/passkeeper/internal/client/proc/chunkerproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/chunkerproc/chunkerpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/getsecretbyidproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/getsecretbyidproc/getsecretbyidprocattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/getsecretsproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/getsecretsproc/getsecretsprocattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/initloadproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/initloadproc/initloadprocattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/initsingleloadproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/initsingleloadproc/initsingleloadprocattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/inituploadproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/inituploadproc/inituploadprocattr.go"
	"github.com/dmitrovia/passkeeper/internal/client/proc/loadproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/loadproc/loadprocattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/loginproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/loginproc/loginprocattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/logoutproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/logoutproc/logoutprocattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/registerproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/registerproc/registerprocpattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadproc/uploadpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadsecretproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/uploadsecretproc/uploadsecretprocattr"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type LoadChunkInfo struct {
	Chunks      map[string]*chunckmeta.ChunkMeta
	LoadIsAllow bool
}
type InteractionProcAttr struct {
	// upload and chunk
	Chproc             *chunkerproc.ChunkerProc
	Chunkerpa          *chunkerpa.ChunkerProcAttr
	Uploadproc         *uploadproc.UploadProc
	Uploadpa           *uploadpa.UploadProcAttr
	UploadChan         chan *chunckmeta.ChunkMeta
	ErrChan            chan error
	SpecificFileUpload bool
	// upload secret
	UploadsecretProc     *uploadsecretproc.UploadSecretProc
	UploadsecretProcAttr *uploadsecretprocattr.
				UploadSecretProcAttr
	// get secret by id
	GetSecretByIDProc     *getsecretbyidproc.GetSecretByID
	GetSecretByIDProcAttr *getsecretbyidprocattr.
				GetSecretByIDProcAttr
	// get all secrets
	GetSecretsProc     *getsecretsproc.GetSecrets
	GetSecretsProcAttr *getsecretsprocattr.GetSecretsProcAttr
	// singleinitload
	InitSingleLoadproc *initsingleloadproc.InitSingleProc
	InitSingleLoadpa   *initsingleloadprocattr.
				InitSingleLoadProcAttr
	// initload
	InitLoadproc *initloadproc.InitProc
	InitLoadpa   *initloadprocattr.InitLoadProcAttr
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
	AttrClintProc   *clientpa.ClientProcAttr
	WgSubProc       *sync.WaitGroup
	WorkerChunkWg   *sync.WaitGroup
	Metamanager     *metamanager.MetaManager
	CurrentMetadata map[string]chunckmeta.ChunkMeta
	// load and build
	LoadChan           chan *chunckmeta.ChunkMeta
	SpecificFileLoad   bool
	LoadMetadata       map[string]chunckmeta.ChunkMeta
	FileNamesAllowLoad map[string]LoadChunkInfo
	LoadProc           *loadproc.LoadProc
	LoadProcAttr       *loadprocattr.LoadProcAttr
}

func (ipa *InteractionProcAttr) SortLoadMetadata() error {
	ipa.FileNamesAllowLoad = make(map[string]LoadChunkInfo)
	for _, value := range ipa.LoadMetadata {
		val, ok := ipa.FileNamesAllowLoad[*value.OrigFileName]
		if !ok {
			lci := &LoadChunkInfo{}
			lci.Chunks = make(map[string]*chunckmeta.ChunkMeta)
			lci.Chunks[*value.FileName] = &value
			lci.LoadIsAllow = true
			ipa.FileNamesAllowLoad[*value.OrigFileName] = *lci
		} else {
			val.Chunks[*value.FileName] = &value
		}
	}

	clear(ipa.LoadMetadata)

	return nil
}

//nolint:nestif
func (ipa *InteractionProcAttr) SetRestrictions() error {
	names, err := ipa.GetExistingFilenames()
	if err != nil {
		return fmt.Errorf("loadAndBuild->GEF: %w", err)
	}

	var inValue string

	for key, value := range ipa.FileNamesAllowLoad {
		_, ok := names[key]
		if ok {
			str := "Do you want to overwrite the file?" +
				"Enter 1 if yes, 2 - all files, other - no"

			fmt.Println("FILE:" + key)
			fmt.Println(str)

			if ipa.AttrClintProc.TestData == nil {
				_, err1 := fmt.Fscan(os.Stdin, &inValue)
				if err1 != nil {
					continue
				}
			} else {
				inValue = ipa.AttrClintProc.TestData.
					TestSetRestrictionsInput
			}

			if inValue == "2" {
				break
			}

			if inValue == "1" {
				continue
			}

			value.LoadIsAllow = false
		}
	}

	return nil
}

func (ipa *InteractionProcAttr,
) InitGetSecrets() error {
	ipa.GetSecretsProcAttr = &getsecretsprocattr.
		GetSecretsProcAttr{}

	ipa.GetSecretsProcAttr.Init(ipa.AttrClintProc)

	ipa.GetSecretsProc = getsecretsproc.NewProc(
		ipa.GetSecretsProcAttr)

	return nil
}

func (ipa *InteractionProcAttr,
) InitGetSecretByID() error {
	ipa.GetSecretByIDProcAttr = &getsecretbyidprocattr.
		GetSecretByIDProcAttr{}

	ipa.GetSecretByIDProcAttr.Init(ipa.AttrClintProc)

	ipa.GetSecretByIDProc = getsecretbyidproc.NewProc(
		ipa.GetSecretByIDProcAttr)

	return nil
}

func (ipa *InteractionProcAttr) InitUploadSecret() error {
	ipa.UploadsecretProcAttr = &uploadsecretprocattr.
		UploadSecretProcAttr{}

	ipa.UploadsecretProcAttr.Init(ipa.AttrClintProc)

	ipa.UploadsecretProc = uploadsecretproc.NewProc(
		ipa.UploadsecretProcAttr)

	return nil
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

	ipa.UploadChan = make(chan *chunckmeta.ChunkMeta,
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

	return nil
}

func (ipa *InteractionProcAttr) InitAfterLoad() error {
	ipa.InitSingleLoadpa = &initsingleloadprocattr.
		InitSingleLoadProcAttr{}
	ipa.InitSingleLoadpa.Init(ipa.AttrClintProc)
	ipa.InitSingleLoadproc = initsingleloadproc.NewProc(
		ipa.InitSingleLoadpa)
	ipa.InitLoadpa = &initloadprocattr.
		InitLoadProcAttr{}
	ipa.InitLoadpa.Init(ipa.AttrClintProc)
	ipa.InitLoadproc = initloadproc.NewProc(
		ipa.InitLoadpa)

	return nil
}

func (ipa *InteractionProcAttr) InitLoad(
	countMeta int,
) error {
	ipa.LoadChan = make(chan *chunckmeta.ChunkMeta,
		countMeta)
	ipa.ErrChan = make(chan error, countMeta)

	ipa.WorkerChunkWg = &sync.WaitGroup{}
	ipa.LoadProcAttr = &loadprocattr.LoadProcAttr{}
	ipa.LoadProcAttr.Init(ipa.AttrClintProc)
	ipa.LoadProcAttr.WorkerChunkWg = ipa.WorkerChunkWg
	ipa.LoadProcAttr.LoadChan = ipa.LoadChan
	ipa.LoadProcAttr.ErrChan = ipa.ErrChan
	ipa.LoadProc = loadproc.NewProc(ipa.LoadProcAttr)

	ipa.Metamanager = metamanager.NewMetaManager(
		ipa.AttrClintProc.MetaPath)

	metadata, err := ipa.Metamanager.LoadMetadata()
	if err != nil {
		return fmt.Errorf("InitLoad->LM: %w", err)
	}

	ipa.CurrentMetadata = metadata

	return nil
}

func (ipa *InteractionProcAttr) GetExistingFilenames() (
	map[string]struct{}, error,
) {
	entries, err := os.ReadDir(
		ipa.AttrClintProc.FilesUploadPath)
	if err != nil {
		return nil, fmt.Errorf(
			"GetExistingFilenames->ReadDir: %w", err)
	}

	names := make(map[string]struct{})

	for _, e := range entries {
		names[e.Name()] = struct{}{}
	}

	return names, nil
}
