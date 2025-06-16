package commandlib

import "fmt"

type ArgValidator = func(value any) (err error)

type argValidationError struct {
	message string
}

func CreateArgValidationError(template string, args ...any) *argValidationError {
	return &argValidationError{
		message: fmt.Sprintf(template, args...),
	}
}

func (err *argValidationError) Error() string {
	return err.message
}

func StringArgTypeValidator(value any) (err error) {
	_, valid := value.(string)

	if !valid {
		err = CreateArgValidationError("Invalid argument type, expected %v", StringArg)
	}

	return
}

func IntArgTypeValidator(value any) (err error) {
	_, valid := value.(int32)

	if !valid {
		err = CreateArgValidationError("Invalid type, expected %v", IntArg)
	}

	return
}

func FloatArgTypeValidator(value any) (err error) {
	_, valid := value.(float32)

	if !valid {
		err = CreateArgValidationError("Invalid type, expected %v", FloatArg)
	}

	return
}
