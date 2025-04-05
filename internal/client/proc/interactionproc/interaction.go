package interactionproc

import (
	"fmt"

	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/interactionpa"
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

	if ip.attr == nil {
		ip.attr = &interactionpa.InteractionProcAttr{}
	}

	err := ip.attr.InitChunkAndUpload()
	if err != nil {
		return fmt.Errorf("RP->initChunkAndUpload: %w", err)
	}

	go ip.RunChunker()
	go ip.RunUploader()

	ip.attr.Wgroup.Wait()
	close(ip.attr.Uploadpa.UploadChan)
	close(ip.attr.Uploadpa.ErrChan)

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
