package dslx

//
// Functional extensions (auto-generated code)
//

// Compose3 is Compose for N=3.
func Compose3[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
](
	f0 Func[T0, *Result[T1]],
	f1 Func[T1, *Result[T2]],
	f2 Func[T2, *Result[T3]],
) Func[T0, *Result[T3]] {
	return Compose2(f0, Compose2(f1, f2))
}

// Compose4 is Compose for N=4.
func Compose4[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
	T4 any,
](
	f0 Func[T0, *Result[T1]],
	f1 Func[T1, *Result[T2]],
	f2 Func[T2, *Result[T3]],
	f3 Func[T3, *Result[T4]],
) Func[T0, *Result[T4]] {
	return Compose2(f0, Compose3(f1, f2, f3))
}

// Compose5 is Compose for N=5.
func Compose5[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
	T4 any,
	T5 any,
](
	f0 Func[T0, *Result[T1]],
	f1 Func[T1, *Result[T2]],
	f2 Func[T2, *Result[T3]],
	f3 Func[T3, *Result[T4]],
	f4 Func[T4, *Result[T5]],
) Func[T0, *Result[T5]] {
	return Compose2(f0, Compose4(f1, f2, f3, f4))
}

// Compose6 is Compose for N=6.
func Compose6[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
	T4 any,
	T5 any,
	T6 any,
](
	f0 Func[T0, *Result[T1]],
	f1 Func[T1, *Result[T2]],
	f2 Func[T2, *Result[T3]],
	f3 Func[T3, *Result[T4]],
	f4 Func[T4, *Result[T5]],
	f5 Func[T5, *Result[T6]],
) Func[T0, *Result[T6]] {
	return Compose2(f0, Compose5(f1, f2, f3, f4, f5))
}

// Compose7 is Compose for N=7.
func Compose7[
	T0 any,
	T1 any,
	T2 any,
	T3 any,
	T4 any,
	T5 any,
	T6 any,
	T7 any,
](
	f0 Func[T0, *Result[T1]],
	f1 Func[T1, *Result[T2]],
	f2 Func[T2, *Result[T3]],
	f3 Func[T3, *Result[T4]],
	f4 Func[T4, *Result[T5]],
	f5 Func[T5, *Result[T6]],
	f6 Func[T6, *Result[T7]],
) Func[T0, *Result[T7]] {
	return Compose2(f0, Compose6(f1, f2, f3, f4, f5, f6))
}

// Compose8 is Compose for N=8.
func Compose8[
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
	f0 Func[T0, *Result[T1]],
	f1 Func[T1, *Result[T2]],
	f2 Func[T2, *Result[T3]],
	f3 Func[T3, *Result[T4]],
	f4 Func[T4, *Result[T5]],
	f5 Func[T5, *Result[T6]],
	f6 Func[T6, *Result[T7]],
	f7 Func[T7, *Result[T8]],
) Func[T0, *Result[T8]] {
	return Compose2(f0, Compose7(f1, f2, f3, f4, f5, f6, f7))
}

// Compose9 is Compose for N=9.
func Compose9[
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
	f0 Func[T0, *Result[T1]],
	f1 Func[T1, *Result[T2]],
	f2 Func[T2, *Result[T3]],
	f3 Func[T3, *Result[T4]],
	f4 Func[T4, *Result[T5]],
	f5 Func[T5, *Result[T6]],
	f6 Func[T6, *Result[T7]],
	f7 Func[T7, *Result[T8]],
	f8 Func[T8, *Result[T9]],
) Func[T0, *Result[T9]] {
	return Compose2(f0, Compose8(f1, f2, f3, f4, f5, f6, f7, f8))
}

// Compose10 is Compose for N=10.
func Compose10[
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
	f0 Func[T0, *Result[T1]],
	f1 Func[T1, *Result[T2]],
	f2 Func[T2, *Result[T3]],
	f3 Func[T3, *Result[T4]],
	f4 Func[T4, *Result[T5]],
	f5 Func[T5, *Result[T6]],
	f6 Func[T6, *Result[T7]],
	f7 Func[T7, *Result[T8]],
	f8 Func[T8, *Result[T9]],
	f9 Func[T9, *Result[T10]],
) Func[T0, *Result[T10]] {
	return Compose2(f0, Compose9(f1, f2, f3, f4, f5, f6, f7, f8, f9))
}

// Compose11 is Compose for N=11.
func Compose11[
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
	f0 Func[T0, *Result[T1]],
	f1 Func[T1, *Result[T2]],
	f2 Func[T2, *Result[T3]],
	f3 Func[T3, *Result[T4]],
	f4 Func[T4, *Result[T5]],
	f5 Func[T5, *Result[T6]],
	f6 Func[T6, *Result[T7]],
	f7 Func[T7, *Result[T8]],
	f8 Func[T8, *Result[T9]],
	f9 Func[T9, *Result[T10]],
	f10 Func[T10, *Result[T11]],
) Func[T0, *Result[T11]] {
	return Compose2(f0, Compose10(f1, f2, f3, f4, f5, f6, f7, f8, f9, f10))
}

// Compose12 is Compose for N=12.
func Compose12[
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
	f0 Func[T0, *Result[T1]],
	f1 Func[T1, *Result[T2]],
	f2 Func[T2, *Result[T3]],
	f3 Func[T3, *Result[T4]],
	f4 Func[T4, *Result[T5]],
	f5 Func[T5, *Result[T6]],
	f6 Func[T6, *Result[T7]],
	f7 Func[T7, *Result[T8]],
	f8 Func[T8, *Result[T9]],
	f9 Func[T9, *Result[T10]],
	f10 Func[T10, *Result[T11]],
	f11 Func[T11, *Result[T12]],
) Func[T0, *Result[T12]] {
	return Compose2(f0, Compose11(f1, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11))
}

// Compose13 is Compose for N=13.
func Compose13[
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
	f0 Func[T0, *Result[T1]],
	f1 Func[T1, *Result[T2]],
	f2 Func[T2, *Result[T3]],
	f3 Func[T3, *Result[T4]],
	f4 Func[T4, *Result[T5]],
	f5 Func[T5, *Result[T6]],
	f6 Func[T6, *Result[T7]],
	f7 Func[T7, *Result[T8]],
	f8 Func[T8, *Result[T9]],
	f9 Func[T9, *Result[T10]],
	f10 Func[T10, *Result[T11]],
	f11 Func[T11, *Result[T12]],
	f12 Func[T12, *Result[T13]],
) Func[T0, *Result[T13]] {
	return Compose2(f0, Compose12(f1, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11, f12))
}

// Compose14 is Compose for N=14.
func Compose14[
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
	f0 Func[T0, *Result[T1]],
	f1 Func[T1, *Result[T2]],
	f2 Func[T2, *Result[T3]],
	f3 Func[T3, *Result[T4]],
	f4 Func[T4, *Result[T5]],
	f5 Func[T5, *Result[T6]],
	f6 Func[T6, *Result[T7]],
	f7 Func[T7, *Result[T8]],
	f8 Func[T8, *Result[T9]],
	f9 Func[T9, *Result[T10]],
	f10 Func[T10, *Result[T11]],
	f11 Func[T11, *Result[T12]],
	f12 Func[T12, *Result[T13]],
	f13 Func[T13, *Result[T14]],
) Func[T0, *Result[T14]] {
	return Compose2(f0, Compose13(f1, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11, f12, f13))
}
