// Package sigcontext provides a Context that exits when the program receives SIGINT or SIGTERM.
package sigcontext

import(
	"context"
	"os"
	"os/signal"
	"syscall"
	"log"
)

// NewForSignals returns a Context that exits when the process receives one of the provided signals.
func NewForSignals(sigs ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, sigs...)

	go func() {
		sig := <-signals
		log.Printf("Received signal %s, starting shutdown", sig.String())
		cancel()
	}()

	return ctx
}

// New returns a Context that exits when the process receives SIGINT or SIGTERM.
func New() context.Context {
	return NewForSignals(syscall.SIGINT, syscall.SIGTERM)
}
