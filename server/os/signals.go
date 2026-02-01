package os

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
)

// Intercept shutdown signals to log shutdown events
func InterceptShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(
		sigs,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)

	go func() {
		<-sigs // Wait for a shutdown signal
		systray.Quit()
	}()
}
