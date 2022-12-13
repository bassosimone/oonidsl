// Package fx contains extensions for functional programming.
//
// The [Func] abstraction represents a function taking in input
// a context and an instance of A and returning a B. Once you have a
// [Func], you use Apply to call the function and obtain its return
// value. You can compose instances of Func using [Compose]. You can
// compose multiple functions using syntactic sugar calls such as
// [Compose3], [Compose4], and so forth. With [Lambda], you can create
// a [Func] on-the-fly from a Go lambda expression.
//
// The [Result] abstraction represents a type containing either an
// instance of type T or an error. As such, it models the result
// of computations where we have network errors. To construct a [Result]
// containing an instance of T, you use the [Ok] function; to construct
// a [Result] containing error, use the [Err] function.
//
// Instances of [Func] returning a [Result] can be composed using the
// [ComposeResult] function. The [Unit] function is the simplest [Func]
// returning [Result] that takes in input an A and returns a [Result]
// containing the same A. There's syntactic sugar for chaining [ComposeResult]
// calls such as [ComposeResult3], [ComposeResult4], etc.
//
// A [Streamable] is an abstraction wrapping a readable Go channel
// where the writer will close the channel when done. With
// the [Zip] function you can read from several [Streamable]. Use [Stream]
// to turn a list of values into a [Streamable].
//
// With [Map] you can apply the same [Func] to multiple values
// using configurable parallelism. The [MapAsync] function is like
// [Map] but uses [Streamable] as input and output. With [Parallel]
// you can apply N distinct [Func] to the same input value using
// configrable parallelism. The [ParallelAsync] function is like
// [Parallel] but uses [Streamable] as input and output.
package fx
