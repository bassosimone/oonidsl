package dslx

import (
	"io"
	"sync"

	"github.com/bassosimone/oonidsl/internal/measurexlite"
	"github.com/lucas-clemente/quic-go"
)

// ConnPool tracks established connections.
type ConnPool struct {
	mu sync.Mutex
	v  []io.Closer
}

// MaybeTrack tracks the given connection if not nil. This
// method is safe for use by multiple goroutines.
func (p *ConnPool) MaybeTrack(c io.Closer) {
	if c != nil {
		defer p.mu.Unlock()
		p.mu.Lock()
		p.v = append(p.v, c)
	}
}

// Close closes all the tracked connections in reverse order. This
// method is safe for use by multiple goroutines.
func (p *ConnPool) Close() error {
	// Implementation note: reverse order is such that we close TLS
	// connections before we close the TCP connections they use. Hence
	// we'll _gracefully_ close TLS connections.
	defer p.mu.Unlock()
	p.mu.Lock()
	for idx := len(p.v) - 1; idx >= 0; idx-- {
		_ = p.v[idx].Close()
	}
	p.v = nil // reset
	return nil
}

// TODO(kelmenhorst) It would be nice to use the same ConnPool for any connection.
// The problem is that quic.EarlyConnection does not implement the io.Closer interface
// because there is no Close (instead: CloseWithError()).
// For now, this is a ConnPool specific to quic.EarlyConnection.

// ConnPool contains established connections.
type QUICConnPool struct {
	mu sync.Mutex
	v  []quic.EarlyConnection
}

// maybeRegister registers a conn for late close if not nil.
func (p *QUICConnPool) maybeRegister(c quic.EarlyConnection) {
	if c != nil {
		defer p.mu.Unlock()
		p.mu.Lock()
		p.v = append(p.v, c)
	}
}

// Close closes all the registered QUIC connections
func (p *QUICConnPool) Close() error {
	defer p.mu.Unlock()
	p.mu.Lock()
	for _, c := range p.v {
		_ = measurexlite.MaybeCloseQUICConn(c)
	}
	p.v = nil // reset
	return nil
}
