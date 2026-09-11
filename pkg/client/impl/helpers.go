package impl

func ValueToPointer[T any](value T) *T {
	return &value
}

func BoolToString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}
