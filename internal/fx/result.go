package fx

//
// Result[T] monad
//

import (
	"context"
	"errors"
)

// Result either contains T or an error.
type Result[T any] interface {
	// IsErr returns whether Result contains an error.
	IsErr() bool

	// Unwrap returns the underlying object or panics
	// if Result actually contains an error.
	Unwrap() T

	// UnwrapErr returns the underlying error or panics
	// if Result actually contains a T.
	UnwrapErr() error
}

// Err constructs a Result containing an error.
func Err[T any](err error) Result[T] {
	return &result[T]{
		err: err,
		t:   *new(T), // zero value
	}
}

// Ok constructs a Result containing a T.
func Ok[T any](t T) Result[T] {
	return &result[T]{
		err: nil,
		t:   t,
	}
}

// result is the private implementation of Result.
type result[T any] struct {
	err error
	t   T
}

// IsErr implements Result
func (r *result[T]) IsErr() bool {
	return r.err != nil
}

// Unwrap implements Result
func (r *result[T]) Unwrap() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.t
}

// ErrNoError indicates that Result[T] does not contain an error.
var ErrNoError = errors.New("fx: Result[T] does not contain an error")

// UnwrapErr implements Result
func (r *result[T]) UnwrapErr() error {
	if r.err == nil {
		panic(ErrNoError)
	}
	return r.err
}

// ComposeResult composes f: A -> Result[B] with g: B -> Result[C]. The
// composition rule is such that, if f returns an error, we construct a Result[C]
// from such an error and return it without invoking g.
func ComposeResult[A, B, C any](f Func[A, Result[B]], g Func[B, Result[C]]) Func[A, Result[C]] {
	return &composeResultFunc[A, B, C]{f: f, g: g}
}

// composeResultFunc[A, B, C] is the type returned by ComposeResult.
type composeResultFunc[A, B, C any] struct {
	// f is the first function to compose.
	f Func[A, Result[B]]

	// g is the second functions to compose.
	g Func[B, Result[C]]
}

// Apply implements Func
func (f *composeResultFunc[A, B, C]) Apply(ctx context.Context, a A) Result[C] {
	r := f.f.Apply(ctx, a)
	if r.IsErr() {
		return Err[C](r.UnwrapErr()) // as documented
	}
	v := r.Unwrap()
	return f.g.Apply(ctx, v)
}

// ComposeResult3 is ComposeResult for N=3.
func ComposeResult3[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
](
	f0 Func[T0, Result[T1]],
	f1 Func[T1, Result[T2]],
	f2 Func[T2, Result[T3]],
) Func[T0, Result[T3]] {
	return ComposeResult(f0, ComposeResult(f1, f2))
}

// Unit returns a function that calls [Ok] on its argument.
func Unit[A any]() Func[A, Result[A]] {
	return Lambda(func(ctx context.Context, a A) Result[A] {
		return Ok(a)
	})
}

//
// Let's pause for a second and check how close we are
// to having actually created a Monad.
//
// We will use as reference "Monads: Programmer's Definition" in
// "Cathegory Theory for Computer Programmers" by Bartosz
// Milewski. See https://github.com/hmemcpy/milewski-ctfp-pdf.
//
// The definition of Monad (pp. 292) is:
//
//     class Monad m where
//       (>=>) :: (a -> m b) -> (b -> m c) -> (a -> m c)
//       return :: a -> m a
//
// Our fish (`>=>`) operator is [ComposeResult]. Our `return`
// equivalent is the [Ok] function.
//
// Then we have to verify monadic laws:
//
//     (f >=> g) >=> h = f >=> (g >=> h)  -- associativity
//     return >=> f    = f                -- left unit
//     f >=> return    = f                -- right unit
//
// Associativity holds because we can write ComposeResult3 as:
//
//     return ComposeResult(ComposeResult(f0, f1), f2)
//
// as well as:
//
//     return ComposeResult(f0, ComposeResult(f1, f2))
//
// Left unit holds because we can write:
//
//     var (
//       f Function[A, Result[B]]
//       a Result[A]
//     )
//
//     g := ComposeResult(Unit[A](), f)
//
// where g is equivalent to f.
//
// Right unit holds because we can write:
//
//     var (
//       f Function[A, Result[B]]
//       a Result[A]
//     )
//
//     g := ComposeResult(f, Unit[B]())
//
// where g is equivalent to f.
//
// We conclude that Result[T] _is_ a monad.
//
