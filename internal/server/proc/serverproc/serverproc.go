package serverproc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dmitrovia/passkeeper/internal/general/functions/loggerf"
	"github.com/dmitrovia/passkeeper/internal/server/migrator"
	"github.com/dmitrovia/passkeeper/internal/server/proc/serverproc/serverpa"
)

type ServerProc struct {
	attr *serverpa.ServerProcAttr
}

func NewProc(attr *serverpa.ServerProcAttr) *ServerProc {
	return &ServerProc{attr: attr}
}

func (sp *ServerProc) RunProcess() error {
	fmt.Println("ServerProc run")
	defer fmt.Println("ServerProc end")

	if sp.attr == nil {
		sp.attr = &serverpa.ServerProcAttr{}
	}

	err := sp.attr.Init()
	if err != nil {
		return fmt.Errorf("RP->Init: %w", err)
	}

	waitGroup := new(sync.WaitGroup)
	ctxDB, cancel := context.WithTimeout(
		context.Background(), sp.attr.Dbtimeout)

	defer cancel()

	err = sp.attr.SetPgxPool(ctxDB)
	if err != nil {
		return fmt.Errorf("RP->SetPgxPool: %w", err)
	}

	err = migrator.UseMigrations(sp.attr)
	if err != nil {
		return fmt.Errorf("RP->UseMigrations: %w", err)
	}

	waitGroup.Add(1)

	go sp.runServer(sp.attr)
	go sp.waitClose(sp.attr, waitGroup)

	waitGroup.Wait()

	return nil
}

func (sp *ServerProc) waitClose(
	attr *serverpa.ServerProcAttr,
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
			err := attr.Server.Shutdown(context.TODO())
			if err != nil {
				loggerf.Log("waitClose->Shutdown", err)
			}

			waitG.Done()

			return
		}
	}
}

func (sp *ServerProc) runServer(
	attr *serverpa.ServerProcAttr,
) {
	err := attr.Server.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		loggerf.Log("runServer->GetServer.LAS", err)
	}
}
