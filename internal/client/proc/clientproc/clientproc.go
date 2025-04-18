package clientproc

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmitrovia/passkeeper/internal/client/proc/clientproc/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc/interactionpa"
	"github.com/dmitrovia/passkeeper/internal/general/functions/loggerf"
)

type ClientProc struct {
	attr *clientpa.ClientProcAttr
}

func NewProc(attr *clientpa.ClientProcAttr) *ClientProc {
	return &ClientProc{attr: attr}
}

func (sp *ClientProc) RunProcess() error {
	fmt.Println("ClientProc run")
	defer fmt.Println("ClientProc end")

	if sp.attr == nil {
		sp.attr = &clientpa.ClientProcAttr{}
	}

	err := sp.attr.Init()
	if err != nil {
		return fmt.Errorf("RP->Init: %w", err)
	}

	sp.attr.WGMainProc.Add(1)

	go sp.waitClose()
	go sp.runInteraction(sp.attr)

	sp.attr.WGMainProc.Wait()
	fmt.Println("Wait for sub processes to complete")
	sp.attr.WgSubProc.Wait()

	return nil
}

func (sp *ClientProc) runInteraction(
	attr *clientpa.ClientProcAttr,
) {
	newAttr := &interactionpa.InteractionProcAttr{}
	newAttr.AttrClintProc = attr
	newAttr.WgSubProc = sp.attr.WgSubProc

	ip := interactionproc.NewProc(newAttr)

	err := ip.RunProcess()
	if err != nil {
		loggerf.Log("runServer->interaction.RunProcess", err)
	}
}

func (sp *ClientProc) waitClose() {
	channelCancel := make(chan os.Signal, 1)
	signal.Notify(channelCancel,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT)

	for {
		_, ok := <-channelCancel
		if ok {
			exitV := 99
			sp.attr.SelectedProc = &exitV

			sp.attr.WGMainProc.Done()

			return
		}
	}
}
