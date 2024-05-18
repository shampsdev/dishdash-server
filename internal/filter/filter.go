package filter

func Map[T, U any](ts []T, mapFunc func(t T) U) []U {
	result := make([]U, 0, len(ts))
	for _, t := range ts {
		result = append(result, mapFunc(t))
	}
	return result
}
