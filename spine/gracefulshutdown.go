package spine

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var C, Cancel = context.WithCancelCause(context.Background())

var SystemGroup sync.WaitGroup

func WaitUntilSystemShutdown() {
	stop := make(chan os.Signal)

	signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	Cancel(nil)
	SystemGroup.Wait()
}
