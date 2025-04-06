package interactionproc

import (
	"fmt"
	"os"

	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc/interactionpa"
	"github.com/dmitrovia/passkeeper/internal/general/functions/loggerf"
)

const (
	registerOption int = 1
	loginOption    int = 2
	uploadOption   int = 3
	exitOption     int = 99
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
	fmt.Println("1.Register")

	if !ip.attr.AttrClintProc.IsAuth {
		fmt.Println("2.Login")
	}

	fmt.Println("3.Send data to server")
	fmt.Println("99.Exit")
	fmt.Println("")
	fmt.Println("--------------------------------------------")
}

func (ip *InteractionProc) chooseProc() error {
	for {
		ip.printOptions()

		if ip.attr.AttrClintProc.SelectedProc != nil &&
			*ip.attr.AttrClintProc.SelectedProc == exitOption {
			fmt.Println("Press ctrl+f to exit")

			return nil
		}

		var inValue int

		_, err1 := fmt.Fscan(os.Stdin, &inValue)
		if err1 != nil {
			return fmt.Errorf("chooseProc->Fscan: %w", err1)
		}

		ip.attr.AttrClintProc.SelectedProc = &inValue

		switch *ip.attr.AttrClintProc.SelectedProc {
		case registerOption:
			fmt.Println("Register")
		case loginOption:
			fmt.Println("Login")
		case exitOption:
			fmt.Println("Press ctrl+f to exit")

			return nil
		case uploadOption:
			fmt.Println("Send data to server")

			err := ip.uploadAndChunk()
			if err != nil {
				loggerf.Log("chooseProc->uploadAndChunk", err)
			}
		default:
			fmt.Println("No such option")
		}
	}
}

func (ip *InteractionProc) uploadAndChunk() error {
	err := ip.attr.InitChunkAndUpload()
	if err != nil {
		return fmt.Errorf("uploadAndChunk->ICAU: %w", err)
	}

	ip.attr.WGsubprocess.Add(
		ip.attr.Chunkerpa.CountWorkersChunker)
	ip.attr.WGsubprocess.Add(
		ip.attr.Uploadpa.CountWorkersUpload)

	go ip.RunChunker()
	go ip.RunUploader()

	ip.attr.WGsubprocess.Wait()

	ip.attr.Chunkerpa.ChFile.Close()
	close(ip.attr.UploadChan)
	close(ip.attr.ErrChan)

	for err := range ip.attr.Uploadpa.ErrChan {
		if err != nil {
			return err
		}
	}

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
