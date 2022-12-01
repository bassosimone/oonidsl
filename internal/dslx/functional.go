package dslx

//
// Functional algorithms
//

import (
	"context"
	"sync"
)

// Function is a generic function from A to B.
type Function[A, B any] interface {
	// Apply applies ctx and A to the function to produce B.
	Apply(ctx context.Context, a A) B
}

// Compose2 composes a function from A to B with a function from B to C.
func Compose2[A, B, C any](f Function[A, B], g Function[B, C]) Function[A, C] {
	return &composer[A, B, C]{f, g}
}

// Compose3 composes three functions together.
func Compose3[A, B, C, D any](f Function[A, B], g Function[B, C], h Function[C, D]) Function[A, D] {
	return Compose2(Compose2(f, g), h)
}

// Compose4 composes four functions together.
func Compose4[A, B, C, D, E any](
	f Function[A, B], g Function[B, C], h Function[C, D], i Function[D, E]) Function[A, E] {
	return Compose2(Compose3(f, g, h), i)
}

// Compose5 composes five functions together.
func Compose5[A, B, C, D, E, F any](f Function[A, B], g Function[B, C], h Function[C, D],
	i Function[D, E], j Function[E, F]) Function[A, F] {
	return Compose2(Compose4(f, g, h, i), j)
}

// Compose6 composes six functions together.
func Compose6[A, B, C, D, E, F, G any](f Function[A, B], g Function[B, C], h Function[C, D],
	i Function[D, E], j Function[E, F], k Function[F, G]) Function[A, G] {
	return Compose2(Compose5(f, g, h, i, j), k)
}

// Compose7 composes seven functions together.
func Compose7[A, B, C, D, E, F, G, H any](f Function[A, B], g Function[B, C], h Function[C, D],
	i Function[D, E], j Function[E, F], k Function[F, G], l Function[G, H]) Function[A, H] {
	return Compose2(Compose6(f, g, h, i, j, k), l)
}

// Compose8 composes eight functions together.
func Compose8[A, B, C, D, E, F, G, H, I any](f Function[A, B], g Function[B, C], h Function[C, D],
	i Function[D, E], j Function[E, F], k Function[F, G], l Function[G, H],
	m Function[H, I]) Function[A, I] {
	return Compose2(Compose7(f, g, h, i, j, k, l), m)
}

// Compose9 composes nine functions together.
func Compose9[A, B, C, D, E, F, G, H, I, J any](f Function[A, B], g Function[B, C], h Function[C, D],
	i Function[D, E], j Function[E, F], k Function[F, G], l Function[G, H],
	m Function[H, I], n Function[I, J]) Function[A, J] {
	return Compose2(Compose8(f, g, h, i, j, k, l, m), n)
}

// Compose10 composes ten functions together.
func Compose10[A, B, C, D, E, F, G, H, I, J, K any](f Function[A, B], g Function[B, C], h Function[C, D],
	i Function[D, E], j Function[E, F], k Function[F, G], l Function[G, H],
	m Function[H, I], n Function[I, J], o Function[J, K]) Function[A, K] {
	return Compose2(Compose9(f, g, h, i, j, k, l, m, n), o)
}

// composer implements Compose.
type composer[A, B, C any] struct {
	f Function[A, B]
	g Function[B, C]
}

// Apply implements Function[A, C].
func (h *composer[A, B, C]) Apply(ctx context.Context, a A) C {
	return h.g.Apply(ctx, h.f.Apply(ctx, a))
}

// Parallelism is the type used to specify parallelism.
type Parallelism int

// Map applies the given function to a list of elements.
//
// Arguments:
//
// - ctx is the context;
//
// - parallelism is the number of goroutines to use (we'll use
// a single goroutine is parallelism is < 1);
//
// - fx is the function to apply;
//
// - as is the list on which to apply fx.
//
// The return value is the list [fx(a)] for every a in A.
func Map[A, B any](
	ctx context.Context,
	parallelism Parallelism,
	fx Function[A, B],
	as ...A,
) []B {
	return MapAsync(ctx, parallelism, fx, Stream(as...)).Collect()
}

// MapAsync is like Map but deals with possibly-very-long sequences.
func MapAsync[A, B any](
	ctx context.Context,
	parallelism Parallelism,
	fx Function[A, B],
	inputs *Streamable[A],
) *Streamable[B] {
	// create channel for returning results
	r := make(chan B)

	// spawn worker goroutines
	wg := &sync.WaitGroup{}
	if parallelism < 1 {
		parallelism = 1
	}
	for i := Parallelism(0); i < parallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for a := range inputs.C {
				r <- fx.Apply(ctx, a)
			}
		}()
	}

	// close channel when done
	go func() {
		defer close(r)
		wg.Wait()
	}()

	return &Streamable[B]{r}
}

// Parallel executes f1...fn functions in parallel over the same input.
//
// Arguments:
//
// - ctx is the context;
//
// - parallelism is the number of goroutines to use (we'll use
// a single goroutine is parallelism is < 1);
//
// - input is the functions' input;
//
// - fn is the list of functions.
//
// The return value is the list [fx(a)] for every fx in fn.
func Parallel[A, B any](
	ctx context.Context,
	parallelism Parallelism,
	input A,
	fn ...Function[A, B],
) []B {
	return ParallelAsync(ctx, parallelism, input, Stream(fn...)).Collect()
}

// ParallelAsync is like Parallel but deals with possibly-very-long sequences.
func ParallelAsync[A, B any](
	ctx context.Context,
	parallelism Parallelism,
	input A,
	funcs *Streamable[Function[A, B]],
) *Streamable[B] {
	// create channel for returning results
	r := make(chan B)

	// spawn worker goroutines
	wg := &sync.WaitGroup{}
	if parallelism < 1 {
		parallelism = 1
	}
	for i := Parallelism(0); i < parallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fx := range funcs.C {
				r <- fx.Apply(ctx, input)
			}
		}()
	}

	// close channel when done
	go func() {
		defer close(r)
		wg.Wait()
	}()

	return &Streamable[B]{r}
}

// Streamable wraps a channel that returns T and is closed
// by the producer when all input has been emitted.
type Streamable[T any] struct {
	C <-chan T
}

// Collect collects all the elements inside a stream.
func (s *Streamable[T]) Collect() (v []T) {
	for t := range s.C {
		v = append(v, t)
	}
	return
}

// Stream creates a Streamable out of static values.
func Stream[T any](ts ...T) *Streamable[T] {
	c := make(chan T)
	go func() {
		defer close(c)
		for _, t := range ts {
			c <- t
		}
	}()
	return &Streamable[T]{c}
}

// Lambda takes in input a lambda and constructs the equivalent Function.
func Lambda[A, B any](fx func(context.Context, A) B) Function[A, B] {
	return &lambda[A, B]{fx}
}

// lambda is the type returned by Lambda.
type lambda[A, B any] struct {
	fun func(context.Context, A) B
}

// Apply implements Function
func (f *lambda[A, B]) Apply(ctx context.Context, a A) B {
	return f.fun(ctx, a)
}

// Mutex[T] wraps T with a mutex.
type Mutex[T any] struct {
	mu sync.Mutex
	v  T
}

// Set sets the underlying wrapped value.
func (m *Mutex[T]) Set(v T) {
	defer m.mu.Unlock()
	m.mu.Lock()
	m.v = v
}

// Get returns the underlying wrapped value.
func (m *Mutex[T]) Get() (v T) {
	defer m.mu.Unlock()
	m.mu.Lock()
	v = m.v
	return
}

// ApplyAsync is equivalent to calling Apply but returns a Streamable.
func ApplyAsync[A, B any](ctx context.Context, fx Function[A, B], input A) *Streamable[B] {
	return MapAsync(ctx, Parallelism(1), fx, Stream(input))
}

// Zip zips together results from multiple streams.
func Zip[T any](sources ...*Streamable[T]) *Streamable[T] {
	r := make(chan T)
	wg := &sync.WaitGroup{}
	for _, src := range sources {
		wg.Add(1)
		go func(s *Streamable[T]) {
			defer wg.Done()
			for e := range s.C {
				r <- e
			}
		}(src)
	}
	go func() {
		defer close(r)
		wg.Wait()
	}()
	return &Streamable[T]{r}
}

// ZipAndCollect chains Zip and Collect.
func ZipAndCollect[T any](sources ...*Streamable[T]) []T {
	return Zip(sources...).Collect()
}

// ErrorOr[T] contains either an error or an instance of T.
type ErrorOr[T any] struct {
	// err is the error
	err error

	// v is the instance of T
	v T
}

// NewErrorOr constructs a new ErrorOr instance.
func NewErrorOr[T any](v T, err error) *ErrorOr[T] {
	if err != nil {
		return &ErrorOr[T]{
			err: err,
			v:   *new(T), // zero value
		}
	}
	return &ErrorOr[T]{
		err: nil,
		v:   v,
	}
}

// Error returns the error or nil
func (eo *ErrorOr[T]) Error() error {
	return eo.err
}

// Unwrap returns the value or calls panic with the underlying error.
func (eo *ErrorOr[T]) Unwrap() T {
	if eo.err != nil {
		panic(eo.err)
	}
	return eo.v
}
