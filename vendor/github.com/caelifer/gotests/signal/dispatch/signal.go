package dispatch

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

// Dispatcher shared global
var dispatcher = struct {
	*sync.Mutex
	signals map[os.Signal]chan os.Signal
}{
	new(sync.Mutex),
	make(map[os.Signal]chan os.Signal),
}

// Package initializer
func init() {
	// Set logger
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
}

// SignalHandler is a custom function that handles os.Signal
type SignalHandler func(os.Signal)

// HandleSignal installs custom handler for a particular os.Signal provided by signal.
func HandleSignal(sig os.Signal, handler SignalHandler) {
	// Unregister handler if it exists
	StopHandleSignal(sig)

	log.Printf("registering new [%s] handler", sig)

	// Create buffered channel of os.Signal values
	ch := make(chan os.Signal, 1)

	/////////////////////// protected section ///////////////////
	// Take exclusive lock
	dispatcher.Lock()

	// Install our new channel
	dispatcher.signals[sig] = ch

	// Fast unlock
	dispatcher.Unlock()
	/////////////////////// protected section ///////////////////

	// Install custom handler in the separate gorutine
	go func(c <-chan os.Signal, sig os.Signal) {
		for s := range c {
			handler(s)
		}
		log.Printf("exiting [%s] handler", sig)
	}(ch, sig)

	// Set notification
	signal.Notify(ch, sig)
}

// StopHandleSignal safely stops signal handling for signal specified by signal.
// If no handler exists, this function is noop.
func StopHandleSignal(sig os.Signal) {
	// Take exclusive lock
	dispatcher.Lock()
	defer dispatcher.Unlock()

	// Check if we already have registered handler
	if ch, ok := dispatcher.signals[sig]; ok {
		// Signal handler already exists - do clean-up
		log.Printf("unregistering existing [%s] handler", sig)

		// Stop receiving signlas
		signal.Stop(ch)
		// Close signal channel so gorutine can safely exit
		close(ch)

		// Clear our signal table
		delete(dispatcher.signals, sig)
	}
}
