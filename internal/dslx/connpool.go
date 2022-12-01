package dslx

import (
	"io"
	"sync"
)

// ConnPool contains established connections.
type ConnPool struct {
	mu sync.Mutex
	v  []io.Closer
}

// maybeRegister registers a conn for late close if not nil.
func (p *ConnPool) maybeRegister(c io.Closer) {
	if c != nil {
		defer p.mu.Unlock()
		p.mu.Lock()
		p.v = append(p.v, c)
	}
}

// Close closes all the registered connections in reverse order
// such that we TLS-close before we try TCP-close
func (p *ConnPool) Close() error {
	defer p.mu.Unlock()
	p.mu.Lock()
	for idx := len(p.v) - 1; idx >= 0; idx-- {
		_ = p.v[idx].Close()
	}
	p.v = nil // reset
	return nil
}
