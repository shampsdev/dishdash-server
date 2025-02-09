package algo

func Map[T, U any](ts []T, mapFunc func(t T) U) []U {
	result := make([]U, 0, len(ts))
	for _, t := range ts {
		result = append(result, mapFunc(t))
	}
	return result
}

func Filter[T any](ts []T, f func(T) bool) []T {
	result := make([]T, 0)
	for _, t := range ts {
		if f(t) {
			result = append(result, t)
		}
	}
	return result
}

func Count[T any](ts []T, f func(T) bool) int {
	result := 0
	for _, t := range ts {
		if f(t) {
			result++
		}
	}
	return result
}
