package signals

import (
	"os"
	"os/signal"
	"syscall"
)

func SetupSignalHandler() (stopCh <-chan struct{}) {
	var shutdownSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c // If a second signal is caught, exit immediately.
		os.Exit(1)
	}()

	return stop
}
