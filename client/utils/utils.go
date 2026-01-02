package utils

func Ternary[T any](cond bool, v1 T, v2 T) T {
	if cond {
		return v1
	} else {
		return v2
	}
}
