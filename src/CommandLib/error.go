package commandlib

import "fmt"

type commandLibError struct {
	cmdName string
	message string
}

func CreateCommandLibError(cmdName string, msg string, msgArgs ...any) *commandLibError {
	return &commandLibError{
		cmdName: cmdName,
		message: fmt.Sprintf(msg, msgArgs...),
	}
}

func (cmdErr *commandLibError) Error() string {
	return "Error with command '" + cmdErr.cmdName + "': " + cmdErr.message
}
