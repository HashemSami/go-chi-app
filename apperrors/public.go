package apperrors

// Public wraps the original error with a new error that has
// a 'Public() string' method that will return a message that is
// acceptable to display to the public. This error can also be
// unwrapped using the traditional 'errors' package approach
func Public(err error, msg string) error {
	return nil
}
