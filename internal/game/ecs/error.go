package ecs

import "fmt"

type ecsError struct {
	err string
}

func newECSError(v ...any) *ecsError {
	return &ecsError{
		err: fmt.Sprint(v...),
	}
}

func newFormattedECSError(format string, v ...any) *ecsError {
	return &ecsError{
		err: fmt.Sprintf(format, v...),
	}
}

func (err *ecsError) Error() string {
	return err.err
}
