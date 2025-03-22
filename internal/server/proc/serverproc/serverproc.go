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

	"github.com/dmitrovia/passkeeper/internal/general/flags"
	"github.com/dmitrovia/passkeeper/internal/general/functions/loggerf"
	"github.com/dmitrovia/passkeeper/internal/server/config"
	"github.com/dmitrovia/passkeeper/internal/server/migrator"
	"github.com/dmitrovia/passkeeper/internal/server/models/procattrs/serverpa"
)

func RunProcess() error {
	attr := &serverpa.ServerProcAttr{}

	err := attr.Init()
	if err != nil {
		return fmt.Errorf("RP->Init: %w", err)
	}

	flags.InitFlags(attr)

	err = config.GetAttrs(attr)
	if err != nil {
		return fmt.Errorf("RP->GetAttrs: %w", err)
	}

	waitGroup := new(sync.WaitGroup)
	ctxDB, cancel := context.WithTimeout(
		context.Background(), attr.DBtimeout)

	defer cancel()

	err = attr.SetPgxConn(ctxDB)
	if err != nil {
		return fmt.Errorf("RP->SetPgxConn: %w", err)
	}

	err = migrator.UseMigrations(attr)
	if err != nil {
		return fmt.Errorf("RP->UseMigrations: %w", err)
	}

	waitGroup.Add(1)

	go runServer(attr)
	go waitClose(attr, waitGroup)

	waitGroup.Wait()

	return nil
}

func waitClose(
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
			err := attr.GetServer().Shutdown(context.TODO())
			if err != nil {
				loggerf.Log("waitClose->Shutdown", err)
			}

			waitG.Done()

			return
		}
	}
}

func runServer(attr *serverpa.ServerProcAttr) {
	err := attr.GetServer().ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		loggerf.Log("runServer->GetServer.LAS", err)
	}
}
