package interactionproc

import (
	"fmt"
	"maps"
	"os"

	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc/interactionpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/logoutproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/logoutproc/logoutprocattr"
	"github.com/dmitrovia/passkeeper/internal/general/functions/loggerf"
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

		err = ip.RunRegister()
	case loginOption:
		fmt.Println("Login")

		err = ip.RunLogin()
	case logoutOption:
		fmt.Println("Logout")

		err = ip.RunLogout()
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
	ip *InteractionProc) loadAndChunksSelectMode() error {
	str := "Enter 1 if you want to load" +
		"a specific file, any other value if all"
	fmt.Println(str)

	var inValue int

	_, err1 := fmt.Fscan(os.Stdin, &inValue)
	if err1 != nil {
		return fmt.Errorf("LACSM->Fscan1: %w", err1)
	}

	if inValue == 1 {
		ip.attr.SpecificFileLoad = true
	} else {
		ip.attr.SpecificFileLoad = false
	}

	ip.attr.WgSubProc.Add(1)

	err := ip.attr.InitLoad()
	if err != nil {
		return fmt.Errorf("LACSM->IL: %w", err)
	}

	if ip.attr.SpecificFileLoad {
		err := ip.attr.InitSingleLoadproc.RunProcess()
		if err != nil {
			return fmt.Errorf("LACSM->ISLPRP: %w", err)
		}

		ip.attr.LoadMetadata = ip.
			attr.InitSingleLoadpa.LoadMetadata
	} else {
		err := ip.attr.InitLoadproc.RunProcess()
		if err != nil {
			return fmt.Errorf("LACSM->ILPRP: %w", err)
		}

		ip.attr.LoadMetadata = ip.
			attr.InitLoadpa.LoadMetadata
	}

	err = ip.loadAndBuild()
	if err != nil {
		return fmt.Errorf("LACSM->UAC: %w", err)
	}

	return nil
}

func (ip *InteractionProc) loadAndBuild() error {
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
	entries, err := os.ReadDir(
		ip.attr.AttrClintProc.FileSynchronizePath)
	if err != nil {
		return fmt.Errorf("UACSM->ReadDir: %w", err)
	}

	for _, e := range entries {
		fsp := ip.attr.AttrClintProc.FileSynchronizePath
		ip.attr.AttrClintProc.SelectFilePath = fsp +
			e.Name()

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

	fsp := ip.attr.AttrClintProc.FileSynchronizePath
	ip.attr.AttrClintProc.SelectFilePath = fsp +
		fileName

	err := ip.uploadAndChunk()
	if err != nil {
		return fmt.Errorf("USF->UAC: %w", err)
	}

	return nil
}

func (ip *InteractionProc) uploadAndChunk() error {
	err := ip.attr.InitChunkAndUpload()
	if err != nil {
		return fmt.Errorf("uploadAndChunk->ICAU: %w", err)
	}

	ip.attr.WgSubProc.Add(1)

	err = ip.attr.InitUploadproc.RunProcess()
	if err != nil {
		return fmt.Errorf("uploadAndChunk->IUPRP: %w", err)
	}

	ip.attr.WgSubProc.Add(
		ip.attr.Chunkerpa.CountWorkersChunker)
	ip.attr.WgSubProc.Add(
		ip.attr.Uploadpa.CountWorkersUpload)
	ip.attr.WorkerChunkWg.Add(ip.attr.Chunkerpa.CntChunks)

	go ip.RunChunker()
	go ip.RunUploader()

	ip.attr.WgSubProc.Wait()

	ip.attr.Chunkerpa.ChFile.Close()
	close(ip.attr.ErrChan)

	if len(ip.attr.Uploadpa.UploadedMetadata) > 0 {
		maps.Copy(ip.attr.CurrentMetadata,
			ip.attr.Uploadpa.UploadedMetadata)

		err = ip.attr.Metamanager.SaveMetadata(
			ip.attr.CurrentMetadata)
		if err != nil {
			return fmt.Errorf("uploadAndChunk->SM: %w", err)
		}
	}

	for err := range ip.attr.Uploadpa.ErrChan {
		if err != nil {
			return err
		}
	}

	fmt.Println("Successfully uploaded")

	return nil
}

func (ip *InteractionProc) RunChunker() {
	err := ip.attr.Chproc.RunProcess()
	if err != nil {
		loggerf.Log("RunChunker->RP", err)

		return
	}
}

func (ip *InteractionProc) RunUploader() {
	err := ip.attr.Uploadproc.RunProcess()
	if err != nil {
		loggerf.Log("RunUploader->RP", err)

		return
	}
}

func (ip *InteractionProc) RunRegister() error {
	ip.attr.WgSubProc.Add(1)

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
