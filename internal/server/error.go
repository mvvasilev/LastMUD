package server

type inputEmptyError struct {
	msg string
}

func newInputEmptyError() *inputEmptyError {
	return &inputEmptyError{
		msg: "No input available at this moment",
	}
}

func (err *inputEmptyError) Error() string {
	return err.msg
}
