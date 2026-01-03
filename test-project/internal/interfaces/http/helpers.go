package http

// enumPtr is a helper function to convert *string to *T for enum types.
// This is useful for handling optional enum fields in update requests.
func EnumPtr[T any](v *string, parse func(string) T) *T {
	if v == nil {
		return nil
	}
	parsed := parse(*v)
	return &parsed
}
