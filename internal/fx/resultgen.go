package fx

//
// Auto-generated code for Result[T]
//

// ComposeResult4 is ComposeResult for N=4.
func ComposeResult4[
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
	return ComposeResult(f0, ComposeResult3(f1, f2, f3))
}

// ComposeResult5 is ComposeResult for N=5.
func ComposeResult5[
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
	return ComposeResult(f0, ComposeResult4(f1, f2, f3, f4))
}

// ComposeResult6 is ComposeResult for N=6.
func ComposeResult6[
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
	return ComposeResult(f0, ComposeResult5(f1, f2, f3, f4, f5))
}

// ComposeResult7 is ComposeResult for N=7.
func ComposeResult7[
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
	return ComposeResult(f0, ComposeResult6(f1, f2, f3, f4, f5, f6))
}

// ComposeResult8 is ComposeResult for N=8.
func ComposeResult8[
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
	return ComposeResult(f0, ComposeResult7(f1, f2, f3, f4, f5, f6, f7))
}

// ComposeResult9 is ComposeResult for N=9.
func ComposeResult9[
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
	return ComposeResult(f0, ComposeResult8(f1, f2, f3, f4, f5, f6, f7, f8))
}

// ComposeResult10 is ComposeResult for N=10.
func ComposeResult10[
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
	return ComposeResult(f0, ComposeResult9(f1, f2, f3, f4, f5, f6, f7, f8, f9))
}

// ComposeResult11 is ComposeResult for N=11.
func ComposeResult11[
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
	return ComposeResult(f0, ComposeResult10(f1, f2, f3, f4, f5, f6, f7, f8, f9, f10))
}

// ComposeResult12 is ComposeResult for N=12.
func ComposeResult12[
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
	return ComposeResult(f0, ComposeResult11(f1, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11))
}

// ComposeResult13 is ComposeResult for N=13.
func ComposeResult13[
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
	return ComposeResult(f0, ComposeResult12(f1, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11, f12))
}

// ComposeResult14 is ComposeResult for N=14.
func ComposeResult14[
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
	return ComposeResult(f0, ComposeResult13(f1, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11, f12, f13))
}
