package misto

import (
	"bufio"
	"context"
	"io"
	"sync"
)

// ConcurrentScanner works like io.Scanner, but with multiple io.Readers
type ConcurrentScanner struct {
	scans  chan []byte   // Scanned data from readers
	errors chan error    // Errors from readers
	done   chan struct{} // Signal that all readers have completed
	cancel func()        // Cancel all readers (stop on first error)

	data []byte // Last scanned value
	err  error
}

// NewConcurrentScanner starts scanning each reader in a separate goroutine
// and returns a *ConcurrentScanner.
func NewConcurrentScanner(readers []io.Reader) *ConcurrentScanner {
	ctx, cancel := context.WithCancel(context.Background())
	s := &ConcurrentScanner{
		scans:  make(chan []byte),
		errors: make(chan error),
		done:   make(chan struct{}),
		cancel: cancel,
	}

	var wg sync.WaitGroup
	wg.Add(len(readers))

	for _, reader := range readers {
		// Start a scanner for each reader in it's own goroutine.
		go func(reader io.Reader) {
			defer wg.Done()
			scanner := bufio.NewScanner(reader)

			for scanner.Scan() {
				select {
				case s.scans <- scanner.Bytes():
					// While there is data, send it to s.scans,
					// this will block until Scan() is called.
				case <-ctx.Done():
					// This fires when context is cancelled,
					// indicating that we should exit now.
					return
				}
			}
			if err := scanner.Err(); err != nil {
				select {
				case s.errors <- err:
					// Reprort we got an error
				case <-ctx.Done():
					// Exit now if context was cancelled, otherwise sending
					// the error and this goroutine will never exit.
					return
				}
			}
		}(reader)
	}

	go func() {
		// Signal that all scanners have completed
		wg.Wait()
		close(s.done)
	}()

	return s
}

// Scan ...
func (s *ConcurrentScanner) Scan() bool {
	select {
	case s.data = <-s.scans:
		// Got data from a scanner
		return true
	case <-s.done:
		// All scanners are done, nothing to do.
	case s.err = <-s.errors:
		// One of the scanners error'd, were done.
	}
	s.cancel() // Cancel context regardless of how we exited.
	return false
}

// Bytes ...
func (s *ConcurrentScanner) Bytes() []byte {
	return s.data
}

// Text ...
func (s *ConcurrentScanner) Text() string {
	return string(s.data)
}

// Err ...
func (s *ConcurrentScanner) Err() error {
	return s.err
}
