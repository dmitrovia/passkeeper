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
	registerOption      string = "1"
	loginOption         string = "2"
	uploadOption        string = "3"
	logoutOption        string = "4"
	loadOption          string = "5"
	uploadSecretOption  string = "6"
	getSecretsOption    string = "7"
	GetSecretByIDOption string = "8"
	exitOption          string = "99"
	nonExistentOption   string = "999"
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

	err := ip.ChooseProc()
	if err != nil {
		return fmt.Errorf("RP->chooseProc: %w", err)
	}

	return nil
}

func (ip *InteractionProc) ChooseProc() error {
	for {
		ip.printOptions()

		if ip.attr.AttrClintProc.SelectedProc != nil &&
			*ip.attr.AttrClintProc.SelectedProc == exitOption {
			fmt.Println("Press ctrl+c to exit")

			return nil
		}

		var inValue string

		if ip.attr.AttrClintProc.TestData == nil {
			_, err1 := fmt.Fscan(os.Stdin, &inValue)
			if err1 != nil {
				return fmt.Errorf("chooseProc->Fscan: %w", err1)
			}
		} else {
			inValue = ip.attr.AttrClintProc.TestData.
				TestChooseProcInput
		}

		ip.checkIncorrectOptionAuth(inValue)
		ip.checkIncorrectOptionNotAuth(inValue)

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
	fmt.Println("6.Upload secret")
	fmt.Println("7.Get all secrets")
	fmt.Println("8.Get secret by identifier")
	fmt.Println("4.Logout")
	fmt.Println("99.Exit")
	fmt.Println("")
	fmt.Println("--------------------------------------------")
}

func (ip *InteractionProc) checkIncorrectOptionAuth(
	option string,
) {
	checkBan := ip.attr.AttrClintProc.IsAuth &&
		(option == registerOption || option == loginOption)

	if checkBan {
		notOption := nonExistentOption
		ip.attr.AttrClintProc.SelectedProc = &notOption
	} else {
		ip.attr.AttrClintProc.SelectedProc = &option
	}
}

func (ip *InteractionProc) checkIncorrectOptionNotAuth(
	option string,
) {
	checkBan := !ip.attr.AttrClintProc.IsAuth &&
		(option == uploadOption ||
			option == logoutOption ||
			option == loadOption ||
			option == uploadSecretOption ||
			option == getSecretsOption ||
			option == GetSecretByIDOption)

	if checkBan {
		notOption := nonExistentOption
		ip.attr.AttrClintProc.SelectedProc = &notOption
	} else {
		ip.attr.AttrClintProc.SelectedProc = &option
	}
}

//nolint:cyclop
func (ip *InteractionProc) selectOption() bool {
	var err error

	switch *ip.attr.AttrClintProc.SelectedProc {
	case registerOption:
		err = ip.RunRegister()
	case loginOption:
		err = ip.RunLogin()
	case logoutOption:
		err = ip.RunLogout()
	case uploadSecretOption:
		err = ip.RunUploadSecret()
	case getSecretsOption:
		err = ip.RunGetSecrets()
	case GetSecretByIDOption:
		err = ip.RunGetSecretByID()
	case exitOption:
		fmt.Println("Press ctrl+c to exit")

		return true
	case loadOption:
		err = ip.LoadAndChunksSelectMode()
	case uploadOption:
		err = ip.UploadAndChunksSelectMode()
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

	var inValue string

	if ip.attr.AttrClintProc.TestData == nil {
		_, err1 := fmt.Fscan(os.Stdin, &inValue)
		if err1 != nil {
			return fmt.Errorf("chooseLoadType->Fscan1: %w", err1)
		}
	} else {
		inValue = ip.attr.AttrClintProc.TestData.
			TestChooseLoadTypeInput
	}

	if inValue == "1" {
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
	ip *InteractionProc) LoadAndChunksSelectMode() error {
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
	ip *InteractionProc) UploadAndChunksSelectMode() error {
	str := "Enter 1 if you want to download" +
		"a specific file, any other value if all"
	fmt.Println(str)

	var inValue string

	if ip.attr.AttrClintProc.TestData == nil {
		_, err1 := fmt.Fscan(os.Stdin, &inValue)
		if err1 != nil {
			return fmt.Errorf("UACSM->Fscan1: %w", err1)
		}
	} else {
		inValue = ip.attr.AttrClintProc.TestData.
			TestUploadAndChunksSelectModeInput
	}

	if inValue == "1" {
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

	if ip.attr.AttrClintProc.TestData == nil {
		_, err1 := fmt.Fscan(os.Stdin, &fileName)
		if err1 != nil {
			return fmt.Errorf("USF->Fscan: %w", err1)
		}
	} else {
		fileName = ip.attr.AttrClintProc.TestData.
			TestUploadSingleFileInput
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

func (ip *InteractionProc) RunRegister() error {
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

func (ip *InteractionProc) RunLogin() error {
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

func (ip *InteractionProc) RunLogout() error {
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

func (ip *InteractionProc) RunGetSecrets() error {
	ip.attr.WgSubProc.Add(1)
	defer ip.attr.WgSubProc.Done()

	err := ip.attr.InitGetSecrets()
	if err != nil {
		return fmt.Errorf("runGetSecrets->IGS: %w", err)
	}

	err = ip.attr.GetSecretsProc.RunProcess()
	if err != nil {
		return fmt.Errorf("runGetSecrets->RP: %w", err)
	}

	return nil
}

func (ip *InteractionProc) RunGetSecretByID() error {
	ip.attr.WgSubProc.Add(1)
	defer ip.attr.WgSubProc.Done()

	err := ip.attr.InitGetSecretByID()
	if err != nil {
		return fmt.Errorf("runGetSecretByID->ISBI: %w", err)
	}

	err = ip.attr.GetSecretByIDProc.RunProcess()
	if err != nil {
		return fmt.Errorf("runGetSecretByID->RP: %w", err)
	}

	return nil
}

func (ip *InteractionProc) RunUploadSecret() error {
	ip.attr.WgSubProc.Add(1)
	defer ip.attr.WgSubProc.Done()

	err := ip.attr.InitUploadSecret()
	if err != nil {
		return fmt.Errorf("runUploadSecret->IL: %w", err)
	}

	err = ip.attr.UploadsecretProc.RunProcess()
	if err != nil {
		return fmt.Errorf("runUploadSecret->RP: %w", err)
	}

	return nil
}
