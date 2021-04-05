package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// signalHandler for storing a signal
type signalHandler struct {
	signal chan interface{}
}

// NewSignal for creating a new signal
func NewSignal() signalHandler {
	sig := signalHandler{}

	sig.signal = make(chan interface{}, 1)

	return sig
}

// Wait for waiting for a signal
func (s *signalHandler) Wait() {
	<-s.signal
}

// Stop for ending a singal
func (s *signalHandler) Stop() {
	s.signal <- true
}

// osSignalHandler for storing a signal
type osSignalHandler struct {
	signal chan os.Signal
}

// NewOSSignal for creating a new signal
func NewOSSignal() osSignalHandler {
	osSig := osSignalHandler{}

	osSig.signal = make(chan os.Signal, 2)
	signal.Notify(
		osSig.signal,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)

	return osSig
}

// Wait for waiting for an OS signal
func (s *osSignalHandler) Wait() {
	<-s.signal
}
