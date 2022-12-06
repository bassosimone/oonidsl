package fx

//
// Lambda constructor
//

import "context"

// Lambda takes in input a lambda and constructs the equivalent Func.
func Lambda[A, B any](fx func(context.Context, A) B) Func[A, B] {
	return &lambda[A, B]{fx}
}

// lambda is the type returned by Lambda.
type lambda[A, B any] struct {
	fun func(context.Context, A) B
}

// Apply implements Func
func (f *lambda[A, B]) Apply(ctx context.Context, a A) B {
	return f.fun(ctx, a)
}
