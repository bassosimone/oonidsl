package dslx

//
// Functional extensions (core)
//

import (
	"context"
	"sync"

	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/netxlite"
	"github.com/bassosimone/oonidsl/internal/runtimex"
)

// Func is a function f: (context.Context, A) -> B.
type Func[A, B any] interface {
	Apply(ctx context.Context, a A) B
}

// Result is the result of an operation implemented by this package
// such as [TCPConnect], [TLSHandshake], etc.
type Result[State any] struct {
	// Error is either the error that occurred or nil.
	Error error

	// Observations contains the collected observations.
	Observations []*Observations

	// Skipped indicates whether an operation decided
	// that subsequent steps should be skipped.
	Skipped bool

	// State contains state passed between function calls.
	State State
}

// Compose2 composes two operations such as [TCPConnect] and [TLSHandshake].
func Compose2[A, B, C any](f Func[A, *Result[B]], g Func[B, *Result[C]]) Func[A, *Result[C]] {
	return &compose2[A, B, C]{
		f: f,
		g: g,
	}
}

// compose2 is the type returned by [Compose2].
type compose2[A, B, C any] struct {
	f Func[A, *Result[B]]
	g Func[B, *Result[C]]
}

// Apply implements Func
func (h *compose2[A, B, C]) Apply(ctx context.Context, a A) *Result[C] {
	mb := h.f.Apply(ctx, a)
	runtimex.Assert(mb != nil, "h.f.Apply returned a nil pointer")
	if mb.Skipped || mb.Error != nil {
		return &Result[C]{
			Error:        mb.Error,
			Observations: mb.Observations,
			Skipped:      mb.Skipped,
			State:        *new(C), // zero value
		}
	}
	mc := h.g.Apply(ctx, mb.State)
	runtimex.Assert(mc != nil, "h.g.Apply returned a nil pointer")
	return &Result[C]{
		Error:        mc.Error,
		Observations: append(mb.Observations, mc.Observations...), // merge observations
		Skipped:      mc.Skipped,
		State:        mc.State,
	}
}

// Counter generates an instance of *CounterState.
func Counter[T any]() *CounterState[T] {
	return &CounterState[T]{}
}

// CounterState allows to count how many times
// a fx.Func[T, fx.Result[T]] is invoked.
type CounterState[T any] struct {
	n atomicx.Int64
}

// Value returns the counter's value.
func (c *CounterState[T]) Value() int64 {
	return c.n.Load()
}

// Func returns a fx.Func[T, fx.Result[T]] that updates the counter.
func (c *CounterState[T]) Func() Func[T, *Result[T]] {
	return &counterFunc[T]{c}
}

// counterFunc is the Func returned by CounterFunc.Func.
type counterFunc[T any] struct {
	c *CounterState[T]
}

// Apply implements Func.
func (c *counterFunc[T]) Apply(ctx context.Context, value T) *Result[T] {
	c.c.n.Add(1)
	return &Result[T]{
		Error:        nil,
		Observations: nil,
		Skipped:      false,
		State:        value,
	}
}

// ErrorLogger logs errors emitted by Func[A, B].
type ErrorLogger struct {
	errors []error
	mu     sync.Mutex
}

// Errors returns the a copy of the internal array of errors and clears
// the internal array of errors as a side effect.
func (e *ErrorLogger) Errors() []error {
	defer e.mu.Unlock()
	e.mu.Lock()
	v := []error{}
	v = append(v, e.errors...)
	e.errors = nil // as documented
	return v
}

// Record records that an error occurred.
func (e *ErrorLogger) Record(err error) {
	defer e.mu.Unlock()
	e.mu.Lock()
	e.errors = append(e.errors, err)
}

// RecordErrors records errors returned by fx.
func RecordErrors[A, B any](logger *ErrorLogger, fx Func[A, *Result[B]]) Func[A, *Result[B]] {
	return &recordErrorsFunc[A, B]{
		fx: fx,
		p:  logger,
	}
}

// recordErrorsFunc is the type returned by ErrorLogger.Wrap.
type recordErrorsFunc[A, B any] struct {
	fx Func[A, *Result[B]]
	p  *ErrorLogger
}

// Apply implements Func.
func (elw *recordErrorsFunc[A, B]) Apply(ctx context.Context, a A) *Result[B] {
	r := elw.fx.Apply(ctx, a)
	if r.Error != nil {
		elw.p.Record(r.Error)
	}
	return r
}

// FirstErrorExcludingBrokenIPv6Errors returns the first error in a list of
// fx.Result[T] excluding errors known to be linked with IPv6 issues.
func FirstErrorExcludingBrokenIPv6Errors[T any](entries ...*Result[T]) error {
	for _, entry := range entries {
		if entry.Error != nil {
			continue
		}
		err := entry.Error
		switch err.Error() {
		case netxlite.FailureNetworkUnreachable, netxlite.FailureHostUnreachable:
			// This class of errors is often times linked with wrongly
			// configured IPv6, therefore we skip them.
		default:
			return err
		}
	}
	return nil
}

// FirstError returns the first error in a list of fx.Result[T].
func FirstError[T any](entries ...*Result[T]) error {
	for _, entry := range entries {
		if entry.Error != nil {
			continue
		}
		return entry.Error
	}
	return nil
}
