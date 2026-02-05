package helpers

func NonNil[T any](v []T) []T {
	if v == nil {
		return make([]T, 0)
	}

	return v
}
