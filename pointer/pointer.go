package pointer

// To returns a pointer to the given value.
func To[T any](value T) *T {
	return &value
}
