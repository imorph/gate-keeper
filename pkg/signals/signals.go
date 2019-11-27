package signals

import (
	"os"
	"os/signal"
)

var onlyOneSignalHandler = make(chan struct{})

// WaitForSigterm waits for either SIGTERM or SIGINT
//
// Returns the caught signal.
func WaitForSigterm() os.Signal {
	close(onlyOneSignalHandler) // panics when called twice
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, shutdownSignals...)
	return <-ch
}
