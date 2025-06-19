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
