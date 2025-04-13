package interactionproc

import (
	"fmt"
	"maps"
	"os"

	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc/interactionpa"
	"github.com/dmitrovia/passkeeper/internal/general/functions/loggerf"
)

const (
	registerOption    int = 1
	loginOption       int = 2
	uploadOption      int = 3
	logoutOption      int = 4
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

func (ip *InteractionProc) printOptions() {
	fmt.Println("--------------------------------------------")
	fmt.Println("")

	if !ip.attr.AttrClintProc.IsAuth {
		fmt.Println("1.Register")
		fmt.Println("2.Login")

		return
	}

	fmt.Println("3.Send data to server")
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
	ip *InteractionProc) uploadAndChunksSelectMode() error {
	str := "Enter 1 if you want to download" +
		"a specific file, any other value if all"
	fmt.Println(str)

	var inValue int

	var fileName string

	_, err1 := fmt.Fscan(os.Stdin, &inValue)
	if err1 != nil {
		return fmt.Errorf("UACSM->Fscan1: %w", err1)
	}

	if inValue == 1 {
		ip.attr.SpecificFileUpload = true
	}

	if ip.attr.SpecificFileUpload {
		fmt.Println("Enter file name")

		_, err1 := fmt.Fscan(os.Stdin, &fileName)
		if err1 != nil {
			return fmt.Errorf("UACSM->Fscan2: %w", err1)
		}

		fsp := ip.attr.AttrClintProc.FileSynchronizePath
		ip.attr.AttrClintProc.SelectFilePath = fsp +
			fileName

		err := ip.uploadAndChunk()
		if err != nil {
			return fmt.Errorf("UACSM->UAC: %w", err)
		}

		return nil
	}

	return nil
}

func (ip *InteractionProc) uploadAndChunk() error {
	err := ip.attr.InitChunkAndUpload()
	if err != nil {
		return fmt.Errorf("uploadAndChunk->ICAU: %w", err)
	}

	err = ip.attr.InitUploadproc.RunProcess()
	if err != nil {
		return fmt.Errorf("uploadAndChunk->IUPRP: %w", err)
	}

	ip.attr.WGsubprocess.Add(
		ip.attr.Chunkerpa.CountWorkersChunker)
	ip.attr.WGsubprocess.Add(
		ip.attr.Uploadpa.CountWorkersUpload)
	ip.attr.WorkerChunkWg.Add(ip.attr.Chunkerpa.CntChunks)

	go ip.RunChunker()
	go ip.RunUploader()

	ip.attr.WGsubprocess.Wait()

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
	ip.attr.WGsubprocess.Add(1)

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
	ip.attr.WGsubprocess.Add(1)

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
	ip.attr.WGsubprocess.Add(1)

	err := ip.attr.InitLogout()
	if err != nil {
		return fmt.Errorf("RunLogout->IL: %w", err)
	}

	err = ip.attr.Logoutproc.RunProcess()
	if err != nil {
		return fmt.Errorf("RunLogout->RP: %w", err)
	}

	return nil
}
