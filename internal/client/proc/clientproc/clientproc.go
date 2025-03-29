package clientproc

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/clientpa"
	"github.com/dmitrovia/passkeeper/internal/client/proc/interaction"
	"github.com/dmitrovia/passkeeper/internal/general/functions/loggerf"
)

func RunProcess() error {
	fmt.Println("ClientProc run")
	defer fmt.Println("ClientProc end")

	attr := &clientpa.ClientProcAttr{}

	err := attr.Init()
	if err != nil {
		return fmt.Errorf("RP->Init: %w", err)
	}

	waitGroup := new(sync.WaitGroup)
	_, cancel := context.WithTimeout(
		context.Background(), attr.GetReqtimeout())

	defer cancel()

	waitGroup.Add(1)

	go runInteraction(attr)
	go waitClose(waitGroup)

	waitGroup.Wait()

	return nil
}

func runInteraction(attr *clientpa.ClientProcAttr) {
	err := interaction.RunProcess(attr)
	if err != nil {
		loggerf.Log("runServer->interaction.RunProcess", err)
	}
}

func waitClose(
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
