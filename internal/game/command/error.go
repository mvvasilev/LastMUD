package command

import "fmt"

type commandError struct {
	cmdName string
	message string
}

func createCommandError(cmdName string, msg string, msgArgs ...any) *commandError {
	return &commandError{
		cmdName: cmdName,
		message: fmt.Sprintf(msg, msgArgs...),
	}
}

func (cmdErr *commandError) Error() string {
	return "Error with command '" + cmdErr.cmdName + "': " + cmdErr.message
}

type commandContextError struct {
	err string
}

func createCommandContextError(err string) *commandContextError {
	return &commandContextError{
		err: err,
	}
}

func (cce *commandContextError) Error() string {
	return cce.err
}
