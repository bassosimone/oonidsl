package fx

//
// Func definition and operations
//

import "context"

// Func is a morphism inside a given category where
// A and B are objects in the category.
type Func[A, B any] interface {
	Apply(ctx context.Context, a A) (b B)
}

// Compose composes f: A -> B with g: B -> C.
func Compose[A, B, C any](f Func[A, B], g Func[B, C]) Func[A, C] {
	return &composeFunc[A, B, C]{f: f, g: g}
}

// composeFunc[A, B, C] is the type returned by Compose.
type composeFunc[A, B, C any] struct {
	f Func[A, B]
	g Func[B, C]
}

// Apply implements Func
func (f *composeFunc[A, B, C]) Apply(ctx context.Context, a A) C {
	return f.g.Apply(ctx, f.f.Apply(ctx, a))
}

// Compose3 composes three functions together.
func Compose3[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
](
	f0 Func[T0, T1],
	f1 Func[T1, T2],
	f2 Func[T2, T3],
) Func[T0, T3] {
	return Compose(f0, Compose(f1, f2))
}
