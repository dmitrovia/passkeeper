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

func (proc *ClientProc) RunProcess() error {
	fmt.Println("ClientProc run")
	defer fmt.Println("ClientProc end")

	if proc.attr == nil {
		proc.attr = &clientpa.ClientProcAttr{}
	}

	err := proc.attr.Init()
	if err != nil {
		return fmt.Errorf("RP->Init: %w", err)
	}

	proc.attr.WGMainProc.Add(1)

	go proc.waitClose()
	go proc.runInteraction()

	proc.attr.WGMainProc.Wait()
	fmt.Println("Wait for sub processes to complete")
	proc.attr.WgSubProc.Wait()

	return nil
}

func (proc *ClientProc) runInteraction() {
	newAttr := &interactionpa.InteractionProcAttr{}
	newAttr.AttrClintProc = proc.attr
	newAttr.WgSubProc = proc.attr.WgSubProc

	ip := interactionproc.NewProc(newAttr)

	err := ip.RunProcess()
	if err != nil {
		loggerf.Log("runInteraction->interaction.RP", err)
	}
}

func (proc *ClientProc) waitClose() {
	channelCancel := make(chan os.Signal, 1)
	signal.Notify(channelCancel,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT)

	for {
		_, ok := <-channelCancel
		if ok {
			exitV := "99"
			proc.attr.SelectedProc = &exitV
			proc.attr.WGMainProc.Done()

			return
		}
	}
}
