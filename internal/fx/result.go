package fx

//
// The Result[T] monad
//

import (
	"context"
	"errors"
)

// Result is a monad containing either an object
// in the category, T, or an error.
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

// FlatMap transforms f: A -> Result[B] in f': Result[A] -> Result[B].
func FlatMap[A, B any](f Func[A, Result[B]]) Func[Result[A], Result[B]] {
	return &flatMapFunc[A, B]{f: f}
}

// flatMapFunc is the type returned by FlatMap.
type flatMapFunc[A, B any] struct {
	f Func[A, Result[B]]
}

// Apply implements Func
func (f *flatMapFunc[A, B]) Apply(ctx context.Context, a Result[A]) Result[B] {
	if a.IsErr() {
		return Err[B](a.UnwrapErr())
	}
	return f.f.Apply(ctx, a.Unwrap())
}

//
// Let's pause for a second and check how close we are to
// having actually created a Maybe-like Monad.
//
// The following is the definition of Monad in Haskell[1]
//
//     class Monad m where
//       (>>=)  :: m a -> (  a -> m b) -> m b
//       (>>)   :: m a ->  m b         -> m b
//       return ::   a                 -> m a
//
// Where `m a` is written `Result[A]`.
//
// Our `>>=` operator is FlatMap().Apply().
//
// Our `return` operator is Ok.
//
// The `>>` (sequence) operator is implicit in the fact that
// golang is an imperative language.
//
// Then we have the three monad laws [2] to satisfy.
//
// Left identity is satisfied because:
//
//     f(Ok(x).Unwrap())           ===   f(x)
//
// Right identity is satisfied because:
//
//     Ok(m.Unwrap())              ===   m
//
// For associativity, given:
//
//     var (
//       f Func[A, Result[B]]
//       g Func[B, Result[C]]
//       h Func[C, Result[D]]
//     )
//
// we have:
//
//     Compose(Compose(FlapMap(f), FlapMap(g)), FlatMap(h)) ===
//     Compose(FlapMap(f), Compose(FlapMap(g), FlapMap(h)))
//
// .. [1] https://wiki.haskell.org/Monad
//
// .. [2] https://wiki.haskell.org/Monad_laws
//
// To conclude, it seems Result[T] is close to being a monad.
//

// ComposeFlat composes f with FlatMap(g).
func ComposeFlat[A, B, C any](f Func[A, Result[B]], g Func[B, Result[C]]) Func[A, Result[C]] {
	return Compose(f, FlatMap(g))
}

// ComposeFlat3 composes-flat three functions together.
func ComposeFlat3[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
](
	f0 Func[T0, Result[T1]],
	f1 Func[T1, Result[T2]],
	f2 Func[T2, Result[T3]],
) Func[T0, Result[T3]] {
	return Compose(f0, FlatMap(ComposeFlat(f1, f2)))
}
