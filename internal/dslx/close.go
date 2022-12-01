package dslx

//
// Closing anything that can be closed.
//

import (
	"context"
	"io"
)

// Close returns a function that closes any closeable state.
func Close[T io.Closer]() Function[T, ErrorOr[T]] {
	return &closer[T]{}
}

// closer is the type returned by Close.
type closer[T io.Closer] struct{}

// Apply implements Function
func (c *closer[T]) Apply(ctx context.Context, state T) ErrorOr[T] {
	_ = state.Close()
	return NewErrorOr(state, nil)
}
