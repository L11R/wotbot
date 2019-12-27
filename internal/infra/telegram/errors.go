package telegram

type humanReadableError interface {
	error
	Human() string
	Cause() error
}

// Human-readable Error
type hrError struct {
	human string
	error error
}

func newHRError(human string, err error) humanReadableError {
	return &hrError{human: human, error: err}
}

// Just to complain error interface, it should be named String() I guess
func (e *hrError) Error() string {
	return e.error.Error()
}

func (e *hrError) Human() string {
	return e.human
}

func (e *hrError) Cause() error {
	return e.error
}
