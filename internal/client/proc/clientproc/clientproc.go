package clientproc

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interactionproc"
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

	waitGroup := new(sync.WaitGroup)
	_, cancel := context.WithTimeout(
		context.Background(), sp.attr.ReqTimeout)

	defer cancel()

	waitGroup.Add(1)

	go sp.runInteraction(sp.attr)
	go sp.waitClose(waitGroup)

	waitGroup.Wait()

	return nil
}

func (sp *ClientProc) runInteraction(
	attr *clientpa.ClientProcAttr,
) {
	ip := interactionproc.NewProc(attr, nil)

	err := ip.RunProcess()
	if err != nil {
		loggerf.Log("runServer->interaction.RunProcess", err)
	}
}

func (sp *ClientProc) waitClose(
	waitG *sync.WaitGroup,
) {
	channelCancel := make(chan os.Signal, 1)
	signal.Notify(channelCancel,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT)

	for {
		_, ok := <-channelCancel
		if ok {
			waitG.Done()

			return
		}
	}
}
