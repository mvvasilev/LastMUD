package commandlib

type ArgumentBase interface {
	Name() string
	ArgType() ArgType
	IsOptional() bool
	Validate(value any) (valid bool, feedback []error)
}

type ArgumentValue interface {
	Value() any
}

type command struct {
	name    string
	altname string

	args []ArgumentBase

	work func(argValues []ArgumentValue) (err error)
}

func CreateCommand(
	name string,
	altname string,
	work func(argValues []ArgumentValue) (err error),
	arguments ...ArgumentBase,
) (cmd *command, err error) {
	var onlyAcceptingOptionals = false

	for _, v := range arguments {
		if !v.IsOptional() && onlyAcceptingOptionals {
			// Optional arguments can only be placed after non-optional ones
			err = CreateCommandLibError(name, "Cannot define non-optional arguments after optional ones.")
			cmd = nil

			return
		}

		if v.IsOptional() {
			onlyAcceptingOptionals = true
		}
	}

	cmd = new(command)

	cmd.name = name
	cmd.altname = altname
	cmd.work = work

	return
}

func (cmd *command) Name() string {
	return cmd.name
}

func (cmd *command) Execute(argValues []ArgumentValue) (err error) {

	for i, v := range cmd.args {
		if i > len(argValues)-1 {
			if !v.IsOptional() {
				return CreateCommandLibError(cmd.name, "Not enough arguments, found %d, expected more", len(argValues))
			} else {
				break // There are no more arg values to process, and the remaining arguments are all optional anyway
			}
		}
	}

	return cmd.work(argValues)
}
