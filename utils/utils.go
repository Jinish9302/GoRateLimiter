package utils

// NoUtilFunction does nothing and always returns 0 and nil error. It exists solely as a placeholder for tests where a util function is expected but not provided.
func NoUtilFunction() (int, error) {
	return 0, nil
}
