package eno

type Error struct {
	Line    int
	Message string
	Err     error
}

func (e *Error) Error() string {
	return e.Message
}

func newError(line int, message string) *Error {
	return &Error{Line: line, Message: message}
}

func wrapError(line int, err error) *Error {
	return &Error{Line: line, Message: err.Error(), Err: err}
}
