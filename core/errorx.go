package core

import (
	"context"
)

// ErrMonitor ErrMonitor
type ErrMonitor struct {
	name string
	err  error
	errc chan error
}

// NewErrMonitor NewErrMonitor
func NewErrMonitor(name string) *ErrMonitor {
	return &ErrMonitor{
		name: name,
		errc: make(chan error),
	}
}

// Monitor Monitor
func (m *ErrMonitor) Monitor(ctx context.Context, ctxCancel context.CancelFunc) {
	select {
	case err, ok := <-m.errc:
		if !ok {
			break
		}
		m.err = err
		ctxCancel()
	case <-ctx.Done():
		ctxCancel()
	}
}

// Receive Receive
func (m *ErrMonitor) Receive(err error) {
	if m.err == nil {
		m.errc <- err
	}
	m.err = err
}

// Close Close
func (m *ErrMonitor) Close() {
	close(m.errc)
}

// Error Error
func (m *ErrMonitor) Error() error {
	return m.err
}
