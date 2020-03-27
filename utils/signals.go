package utils

import (
	"os"
	"os/signal"
)

func InterruptSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for {
		<-c
		os.Exit(0)
	}
}
