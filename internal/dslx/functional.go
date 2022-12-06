package dslx

//
// Additional function algorithms
//

import (
	"context"
	"sync"

	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/fx"
)

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
func (c *CounterState[T]) Func() fx.Func[T, fx.Result[T]] {
	return &counterFunc[T]{c}
}

// counterFunc is the Func returned by CounterFunc.Func.
type counterFunc[T any] struct {
	c *CounterState[T]
}

// Apply implements Func.
func (c *counterFunc[T]) Apply(ctx context.Context, value T) fx.Result[T] {
	c.c.n.Add(1)
	return fx.Ok(value)
}

// ErrorLogger logs errors emitted by Func[A, B].
type ErrorLogger[A, B any] struct {
	errors []error
	mu     sync.Mutex
}

// Errors returns the a copy of the internal array of errors and clears
// the internal array of errors as a side effect.
func (e *ErrorLogger[A, B]) Errors() []error {
	defer e.mu.Unlock()
	e.mu.Lock()
	v := []error{}
	v = append(v, e.errors...)
	e.errors = nil // as documented
	return v
}

// record records that an error occurred.
func (e *ErrorLogger[A, B]) record(err error) {
	defer e.mu.Unlock()
	e.mu.Lock()
	e.errors = append(e.errors, err)
}

// Wrap wraps the given function to log errors.
func (e *ErrorLogger[A, B]) Wrap(fx fx.Func[A, fx.Result[B]]) fx.Func[A, fx.Result[B]] {
	return &errorLoggerWrapper[A, B]{
		fx: fx,
		p:  e,
	}
}

// errorLoggerWrapper is the type returned by ErrorLogger.Wrap.
type errorLoggerWrapper[A, B any] struct {
	fx fx.Func[A, fx.Result[B]]
	p  *ErrorLogger[A, B]
}

// Apply implements Func.
func (elw *errorLoggerWrapper[A, B]) Apply(ctx context.Context, a A) fx.Result[B] {
	r := elw.fx.Apply(ctx, a)
	if r.IsErr() {
		elw.p.record(r.UnwrapErr())
	}
	return r
}
