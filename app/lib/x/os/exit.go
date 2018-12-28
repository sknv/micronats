package os

import (
	"os"
	"os/signal"
	"syscall"
)

func NotifyAboutExit() <-chan os.Signal {
	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	return exit
}

func WaitForExit() {
	<-NotifyAboutExit()
}
