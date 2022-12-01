package dslx

//
// Closing anything that can be closed.
//

import (
	"context"
	"io"
)

// Close returns a function that closes any closeable state.
func Close[T io.Closer]() Function[*ErrorOr[T], *ErrorOr[T]] {
	return &closer[T]{}
}

// closer is the type returned by Close.
type closer[T io.Closer] struct{}

// Apply implements Function
func (c *closer[T]) Apply(ctx context.Context, maybeState *ErrorOr[T]) *ErrorOr[T] {
	if maybeState.Error() == nil {
		_ = maybeState.Unwrap().Close()
	}
	return maybeState
}
