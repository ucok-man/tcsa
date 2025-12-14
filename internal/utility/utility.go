package utility

import "math"

type mapFunc[E any] func(E) E

func SlicesMap[S ~[]E, E any](s S, f mapFunc[E]) S {
	result := make(S, len(s))
	for i := range s {
		result[i] = f(s[i])
	}
	return result
}

func SetPtrValue[T any](v T) *T {
	return &v
}

func DerefOrDefault[T any](ptr *T, defaultVal T) T {
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

func Round2(v float64) float64 {
	return math.Round(v*100) / 100
}
