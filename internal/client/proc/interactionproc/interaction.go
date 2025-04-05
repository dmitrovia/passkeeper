package interactionproc

import (
	"fmt"

	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc/interactionpa"
	"github.com/dmitrovia/passkeeper/internal/general/functions/loggerf"
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

	err := ip.attr.InitChunkAndUpload()
	if err != nil {
		return fmt.Errorf("RP->initChunkAndUpload: %w", err)
	}

	ip.attr.Wgroup.Add(ip.attr.Chunkerpa.CountWorkersChunker)
	ip.attr.Wgroup.Add(ip.attr.Uploadpa.CountWorkersUpload)

	go ip.RunChunker()
	go ip.RunUploader()

	ip.attr.Wgroup.Wait()

	err = ip.endProc()
	if err != nil {
		return fmt.Errorf("RP->endProc: %w", err)
	}

	return nil
}

func (ip *InteractionProc) endProc() error {
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
