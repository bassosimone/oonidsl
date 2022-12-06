package fx

//
// Auto-generated code for Result[T]
//

// ComposeFlat4 composes-flat N=4 functions together
func ComposeFlat4[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
	T4 any,
](
	f0 Func[T0, Result[T1]],
	f1 Func[T1, Result[T2]],
	f2 Func[T2, Result[T3]],
	f3 Func[T3, Result[T4]],
) Func[T0, Result[T4]] {
	return Compose(f0, FlatMap(ComposeFlat3(f1, f2, f3)))
}

// ComposeFlat5 composes-flat N=5 functions together
func ComposeFlat5[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
	T4 any,
	T5 any,
](
	f0 Func[T0, Result[T1]],
	f1 Func[T1, Result[T2]],
	f2 Func[T2, Result[T3]],
	f3 Func[T3, Result[T4]],
	f4 Func[T4, Result[T5]],
) Func[T0, Result[T5]] {
	return Compose(f0, FlatMap(ComposeFlat4(f1, f2, f3, f4)))
}

// ComposeFlat6 composes-flat N=6 functions together
func ComposeFlat6[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
	T4 any,
	T5 any,
	T6 any,
](
	f0 Func[T0, Result[T1]],
	f1 Func[T1, Result[T2]],
	f2 Func[T2, Result[T3]],
	f3 Func[T3, Result[T4]],
	f4 Func[T4, Result[T5]],
	f5 Func[T5, Result[T6]],
) Func[T0, Result[T6]] {
	return Compose(f0, FlatMap(ComposeFlat5(f1, f2, f3, f4, f5)))
}

// ComposeFlat7 composes-flat N=7 functions together
func ComposeFlat7[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
	T4 any,
	T5 any,
	T6 any,
	T7 any,
](
	f0 Func[T0, Result[T1]],
	f1 Func[T1, Result[T2]],
	f2 Func[T2, Result[T3]],
	f3 Func[T3, Result[T4]],
	f4 Func[T4, Result[T5]],
	f5 Func[T5, Result[T6]],
	f6 Func[T6, Result[T7]],
) Func[T0, Result[T7]] {
	return Compose(f0, FlatMap(ComposeFlat6(f1, f2, f3, f4, f5, f6)))
}

// ComposeFlat8 composes-flat N=8 functions together
func ComposeFlat8[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
	T4 any,
	T5 any,
	T6 any,
	T7 any,
	T8 any,
](
	f0 Func[T0, Result[T1]],
	f1 Func[T1, Result[T2]],
	f2 Func[T2, Result[T3]],
	f3 Func[T3, Result[T4]],
	f4 Func[T4, Result[T5]],
	f5 Func[T5, Result[T6]],
	f6 Func[T6, Result[T7]],
	f7 Func[T7, Result[T8]],
) Func[T0, Result[T8]] {
	return Compose(f0, FlatMap(ComposeFlat7(f1, f2, f3, f4, f5, f6, f7)))
}

// ComposeFlat9 composes-flat N=9 functions together
func ComposeFlat9[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
	T4 any,
	T5 any,
	T6 any,
	T7 any,
	T8 any,
	T9 any,
](
	f0 Func[T0, Result[T1]],
	f1 Func[T1, Result[T2]],
	f2 Func[T2, Result[T3]],
	f3 Func[T3, Result[T4]],
	f4 Func[T4, Result[T5]],
	f5 Func[T5, Result[T6]],
	f6 Func[T6, Result[T7]],
	f7 Func[T7, Result[T8]],
	f8 Func[T8, Result[T9]],
) Func[T0, Result[T9]] {
	return Compose(f0, FlatMap(ComposeFlat8(f1, f2, f3, f4, f5, f6, f7, f8)))
}

// ComposeFlat10 composes-flat N=10 functions together
func ComposeFlat10[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
	T4 any,
	T5 any,
	T6 any,
	T7 any,
	T8 any,
	T9 any,
	T10 any,
](
	f0 Func[T0, Result[T1]],
	f1 Func[T1, Result[T2]],
	f2 Func[T2, Result[T3]],
	f3 Func[T3, Result[T4]],
	f4 Func[T4, Result[T5]],
	f5 Func[T5, Result[T6]],
	f6 Func[T6, Result[T7]],
	f7 Func[T7, Result[T8]],
	f8 Func[T8, Result[T9]],
	f9 Func[T9, Result[T10]],
) Func[T0, Result[T10]] {
	return Compose(f0, FlatMap(ComposeFlat9(f1, f2, f3, f4, f5, f6, f7, f8, f9)))
}

// ComposeFlat11 composes-flat N=11 functions together
func ComposeFlat11[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
	T4 any,
	T5 any,
	T6 any,
	T7 any,
	T8 any,
	T9 any,
	T10 any,
	T11 any,
](
	f0 Func[T0, Result[T1]],
	f1 Func[T1, Result[T2]],
	f2 Func[T2, Result[T3]],
	f3 Func[T3, Result[T4]],
	f4 Func[T4, Result[T5]],
	f5 Func[T5, Result[T6]],
	f6 Func[T6, Result[T7]],
	f7 Func[T7, Result[T8]],
	f8 Func[T8, Result[T9]],
	f9 Func[T9, Result[T10]],
	f10 Func[T10, Result[T11]],
) Func[T0, Result[T11]] {
	return Compose(f0, FlatMap(ComposeFlat10(f1, f2, f3, f4, f5, f6, f7, f8, f9, f10)))
}

// ComposeFlat12 composes-flat N=12 functions together
func ComposeFlat12[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
	T4 any,
	T5 any,
	T6 any,
	T7 any,
	T8 any,
	T9 any,
	T10 any,
	T11 any,
	T12 any,
](
	f0 Func[T0, Result[T1]],
	f1 Func[T1, Result[T2]],
	f2 Func[T2, Result[T3]],
	f3 Func[T3, Result[T4]],
	f4 Func[T4, Result[T5]],
	f5 Func[T5, Result[T6]],
	f6 Func[T6, Result[T7]],
	f7 Func[T7, Result[T8]],
	f8 Func[T8, Result[T9]],
	f9 Func[T9, Result[T10]],
	f10 Func[T10, Result[T11]],
	f11 Func[T11, Result[T12]],
) Func[T0, Result[T12]] {
	return Compose(f0, FlatMap(ComposeFlat11(f1, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11)))
}

// ComposeFlat13 composes-flat N=13 functions together
func ComposeFlat13[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
	T4 any,
	T5 any,
	T6 any,
	T7 any,
	T8 any,
	T9 any,
	T10 any,
	T11 any,
	T12 any,
	T13 any,
](
	f0 Func[T0, Result[T1]],
	f1 Func[T1, Result[T2]],
	f2 Func[T2, Result[T3]],
	f3 Func[T3, Result[T4]],
	f4 Func[T4, Result[T5]],
	f5 Func[T5, Result[T6]],
	f6 Func[T6, Result[T7]],
	f7 Func[T7, Result[T8]],
	f8 Func[T8, Result[T9]],
	f9 Func[T9, Result[T10]],
	f10 Func[T10, Result[T11]],
	f11 Func[T11, Result[T12]],
	f12 Func[T12, Result[T13]],
) Func[T0, Result[T13]] {
	return Compose(f0, FlatMap(ComposeFlat12(f1, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11, f12)))
}

// ComposeFlat14 composes-flat N=14 functions together
func ComposeFlat14[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
	T4 any,
	T5 any,
	T6 any,
	T7 any,
	T8 any,
	T9 any,
	T10 any,
	T11 any,
	T12 any,
	T13 any,
	T14 any,
](
	f0 Func[T0, Result[T1]],
	f1 Func[T1, Result[T2]],
	f2 Func[T2, Result[T3]],
	f3 Func[T3, Result[T4]],
	f4 Func[T4, Result[T5]],
	f5 Func[T5, Result[T6]],
	f6 Func[T6, Result[T7]],
	f7 Func[T7, Result[T8]],
	f8 Func[T8, Result[T9]],
	f9 Func[T9, Result[T10]],
	f10 Func[T10, Result[T11]],
	f11 Func[T11, Result[T12]],
	f12 Func[T12, Result[T13]],
	f13 Func[T13, Result[T14]],
) Func[T0, Result[T14]] {
	return Compose(f0, FlatMap(ComposeFlat13(f1, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11, f12, f13)))
}
