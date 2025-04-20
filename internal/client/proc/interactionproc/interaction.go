package interactionproc

import (
	"fmt"
	"os"

	"github.com/dmitrovia/passkeeper/internal/client/proc/buildproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/buildproc/buildprocattr"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc/interactionpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/logoutproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/logoutproc/logoutprocattr"
	"github.com/dmitrovia/passkeeper/internal/general/functions/loggerf"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

const (
	registerOption    int = 1
	loginOption       int = 2
	uploadOption      int = 3
	logoutOption      int = 4
	loadOption        int = 5
	exitOption        int = 99
	nonExistentOption     = 999
)

type InteractionProc struct {
	attr *interactionpa.InteractionProcAttr
}

func NewProc(
	attr *interactionpa.InteractionProcAttr,
) *InteractionProc {
	return &InteractionProc{
		attr: attr,
	}
}

func (ip *InteractionProc) RunProcess() error {
	fmt.Println("InteractionProc run")
	defer fmt.Println("InteractionProc end")

	err := ip.chooseProc()
	if err != nil {
		return fmt.Errorf("RP->chooseProc: %w", err)
	}

	return nil
}

func (ip *InteractionProc) chooseProc() error {
	for {
		ip.printOptions()

		if ip.attr.AttrClintProc.SelectedProc != nil &&
			*ip.attr.AttrClintProc.SelectedProc == exitOption {
			fmt.Println("Press ctrl+c to exit")

			return nil
		}

		var inValue int

		_, err1 := fmt.Fscan(os.Stdin, &inValue)
		if err1 != nil {
			return fmt.Errorf("chooseProc->Fscan: %w", err1)
		}

		ip.checkIncorrectOption(inValue)

		isExit := ip.selectOption()
		if isExit {
			return nil
		}
	}
}

func (ip *InteractionProc) printOptions() {
	fmt.Println("--------------------------------------------")
	fmt.Println("")

	if !ip.attr.AttrClintProc.IsAuth {
		fmt.Println("1.Register")
		fmt.Println("2.Login")

		return
	}

	fmt.Println("3.Send data to server")
	fmt.Println("5.Get data from server")
	fmt.Println("4.Logout")
	fmt.Println("99.Exit")
	fmt.Println("")
	fmt.Println("--------------------------------------------")
}

func (ip *InteractionProc) checkIncorrectOption(
	option int,
) {
	checkBan := ip.attr.AttrClintProc.IsAuth &&
		(option == registerOption || option == loginOption)
	checkBan1 := !ip.attr.AttrClintProc.IsAuth &&
		(option == uploadOption)

	if checkBan || checkBan1 {
		notOption := nonExistentOption
		ip.attr.AttrClintProc.SelectedProc = &notOption
	} else {
		ip.attr.AttrClintProc.SelectedProc = &option
	}
}

func (ip *InteractionProc) selectOption() bool {
	var err error

	switch *ip.attr.AttrClintProc.SelectedProc {
	case registerOption:
		fmt.Println("Register")

		err = ip.runRegister()
	case loginOption:
		fmt.Println("Login")

		err = ip.runLogin()
	case logoutOption:
		fmt.Println("Logout")

		err = ip.runLogout()
	case exitOption:
		fmt.Println("Press ctrl+c to exit")

		return true
	case loadOption:
		err = ip.loadAndChunksSelectMode()
	case uploadOption:
		err = ip.uploadAndChunksSelectMode()
	default:
		fmt.Println("No such option")
	}

	if err != nil {
		loggerf.Log("selectOption", err)
	}

	return false
}

func (
	ip *InteractionProc) chooseLoadType() error {
	str := "Enter 1 if you want to load" +
		"a specific file, any other value if all"
	fmt.Println(str)

	var inValue int

	_, err1 := fmt.Fscan(os.Stdin, &inValue)
	if err1 != nil {
		return fmt.Errorf("chooseLoadType->Fscan1: %w", err1)
	}

	if inValue == 1 {
		ip.attr.SpecificFileLoad = true
	} else {
		ip.attr.SpecificFileLoad = false
	}

	return nil
}

func (
	ip *InteractionProc) getLoadMetadata() error {
	if ip.attr.SpecificFileLoad {
		err := ip.attr.InitSingleLoadproc.RunProcess()
		if err != nil {
			return fmt.Errorf("getLoadMetadata->ISLPRP: %w", err)
		}

		ip.attr.LoadMetadata = ip.
			attr.InitSingleLoadpa.LoadMetadata
	} else {
		err := ip.attr.InitLoadproc.RunProcess()
		if err != nil {
			return fmt.Errorf("getLoadMetadata->ILPRP: %w", err)
		}

		ip.attr.LoadMetadata = ip.
			attr.InitLoadpa.LoadMetadata
	}

	return nil
}

func (
	ip *InteractionProc) loadAndChunksSelectMode() error {
	err := ip.chooseLoadType()
	if err != nil {
		return fmt.Errorf("LACSM->chooseLoadType: %w", err)
	}

	ip.attr.WgSubProc.Add(1)
	defer ip.attr.WgSubProc.Done()

	err = ip.attr.InitAfterLoad()
	if err != nil {
		return fmt.Errorf("LACSM->IAL: %w", err)
	}

	err = ip.getLoadMetadata()
	if err != nil {
		return fmt.Errorf("LACSM->SLM: %w", err)
	}

	if len(ip.attr.LoadMetadata) == 0 {
		return nil
	}

	err = ip.attr.SortLoadMetadata()
	if err != nil {
		return fmt.Errorf("LACSM->SLM: %w", err)
	}

	err = ip.attr.SetRestrictions()
	if err != nil {
		return fmt.Errorf("LACSM->SetRestrictions: %w", err)
	}

	for fileName, value := range ip.attr.FileNamesAllowLoad {
		if !value.LoadIsAllow {
			continue
		}

		err = ip.loadAndBuild(fileName, value.Chunks)
		if err != nil {
			return fmt.Errorf("LACSM->LAB: %w", err)
		}
	}

	return nil
}

func (ip *InteractionProc) loadAndBuild(
	fileName string,
	metas map[string]*chunckmeta.ChunkMeta,
) error {
	err := ip.attr.InitLoad(len(metas))
	if err != nil {
		return fmt.Errorf("LACSM->IL: %w", err)
	}

	for _, val := range metas {
		ip.attr.LoadChan <- val
	}

	ip.attr.WorkerChunkWg.Add(len(metas))

	go ip.runLoader()
	ip.attr.WorkerChunkWg.Wait()

	close(ip.attr.ErrChan)
	close(ip.attr.LoadChan)

	for err := range ip.attr.ErrChan {
		if err != nil {
			return err
		}
	}

	builderAttr := &buildprocattr.BuildProcAttr{}
	builderAttr.BuildMetadata = metas
	builderAttr.CurrentMetadata = ip.attr.CurrentMetadata
	newPath := fmt.Sprintf("%s%s",
		ip.attr.AttrClintProc.FilesUploadPath,
		fileName)
	builderAttr.OutFilePath = newPath

	builder := buildproc.NewProc(builderAttr)

	err = builder.RunProcess()
	if err != nil {
		return fmt.Errorf("LAB->RPBuilder: %w", err)
	}

	err = ip.attr.Metamanager.SaveMetadata(
		ip.attr.CurrentMetadata)
	if err != nil {
		return fmt.Errorf("LAB->SM: %w", err)
	}

	fmt.Println("Successfully loaded and builded")

	return nil
}

func (
	ip *InteractionProc) uploadAndChunksSelectMode() error {
	str := "Enter 1 if you want to download" +
		"a specific file, any other value if all"
	fmt.Println(str)

	var inValue int

	_, err1 := fmt.Fscan(os.Stdin, &inValue)
	if err1 != nil {
		return fmt.Errorf("UACSM->Fscan1: %w", err1)
	}

	if inValue == 1 {
		ip.attr.SpecificFileUpload = true
	} else {
		ip.attr.SpecificFileUpload = false
	}

	if ip.attr.SpecificFileUpload {
		err := ip.uploadSingleFile()
		if err != nil {
			return fmt.Errorf("UACSM->uploadSingleFile: %w", err)
		}

		return nil
	}

	err := ip.uploadMultiFiles()
	if err != nil {
		return fmt.Errorf("UACSM->uploadMultiFiles: %w", err)
	}

	return nil
}

func (ip *InteractionProc) uploadMultiFiles() error {
	names, err := ip.attr.GetExistingFilenames()
	if err != nil {
		return fmt.Errorf("uploadMultiFiles->GEFN: %w", err)
	}

	for key := range names {
		fsp := ip.attr.AttrClintProc.FilesUploadPath
		ip.attr.AttrClintProc.SelectFilePath = fsp +
			key

		err := ip.uploadAndChunk()
		if err != nil {
			return fmt.Errorf("UACSM->UAC: %w", err)
		}
	}

	return nil
}

func (ip *InteractionProc) uploadSingleFile() error {
	fmt.Println("Enter file name")

	var fileName string

	_, err1 := fmt.Fscan(os.Stdin, &fileName)
	if err1 != nil {
		return fmt.Errorf("USF->Fscan: %w", err1)
	}

	fsp := ip.attr.AttrClintProc.FilesUploadPath
	ip.attr.AttrClintProc.SelectFilePath = fsp +
		fileName

	err := ip.uploadAndChunk()
	if err != nil {
		return fmt.Errorf("USF->UAC: %w", err)
	}

	return nil
}

func (ip *InteractionProc) uploadAndChunk() error {
	ip.attr.WgSubProc.Add(1)
	defer ip.attr.WgSubProc.Done()

	err := ip.attr.InitChunkAndUpload()
	if err != nil {
		return fmt.Errorf("uploadAndChunk->ICAU: %w", err)
	}

	err = ip.attr.InitUploadproc.RunProcess()
	if err != nil {
		return fmt.Errorf("uploadAndChunk->IUPRP: %w", err)
	}

	ip.attr.WorkerChunkWg.Add(ip.attr.Chunkerpa.CntChunks)

	go ip.runChunker()
	go ip.runUploader()

	ip.attr.WorkerChunkWg.Wait()
	ip.attr.Chunkerpa.ChFile.Close()
	close(ip.attr.ErrChan)
	close(ip.attr.UploadChan)

	for err := range ip.attr.ErrChan {
		if err != nil {
			return err
		}
	}

	err = ip.attr.Metamanager.SaveMetadata(
		ip.attr.CurrentMetadata)
	if err != nil {
		return fmt.Errorf("uploadAndChunk->SM: %w", err)
	}

	fmt.Println("Successfully uploaded")

	return nil
}

func (ip *InteractionProc) runChunker() {
	err := ip.attr.Chproc.RunProcess()
	if err != nil {
		loggerf.Log("RunChunker->RP", err)

		return
	}
}

func (ip *InteractionProc) runUploader() {
	err := ip.attr.Uploadproc.RunProcess()
	if err != nil {
		loggerf.Log("RunUploader->RP", err)

		return
	}
}

func (ip *InteractionProc) runLoader() {
	err := ip.attr.LoadProc.RunProcess()
	if err != nil {
		loggerf.Log("RunLoader->RP", err)

		return
	}
}

func (ip *InteractionProc) runRegister() error {
	ip.attr.WgSubProc.Add(1)
	defer ip.attr.WgSubProc.Done()

	err := ip.attr.InitRegister()
	if err != nil {
		return fmt.Errorf("RunRegister->IR: %w", err)
	}

	err = ip.attr.Registerproc.RunProcess()
	if err != nil {
		return fmt.Errorf("RunRegister->RP: %w", err)
	}

	return nil
}

func (ip *InteractionProc) runLogin() error {
	ip.attr.WgSubProc.Add(1)
	defer ip.attr.WgSubProc.Done()

	err := ip.attr.InitLogin()
	if err != nil {
		return fmt.Errorf("RunLogin->IL: %w", err)
	}

	err = ip.attr.Loginproc.RunProcess()
	if err != nil {
		return fmt.Errorf("RunLogin->RP: %w", err)
	}

	return nil
}

func (ip *InteractionProc) runLogout() error {
	ip.attr.WgSubProc.Add(1)
	defer ip.attr.WgSubProc.Done()

	ip.attr.Logoutpa = &logoutprocattr.LogoutProcAttr{}

	err := ip.attr.Logoutpa.Init(ip.attr.AttrClintProc)
	if err != nil {
		return fmt.Errorf("RunLogout->Init: %w", err)
	}

	ip.attr.Logoutproc = logoutproc.NewProc(ip.attr.Logoutpa)

	err = ip.attr.Logoutproc.RunProcess()
	if err != nil {
		return fmt.Errorf("RunLogout->RP: %w", err)
	}

	return nil
}
