package commandlib

type ArgType byte

const (
	StringArg = iota
	IntArg
	FloatArg
)

func (a ArgType) String() string {
	switch a {
	case StringArg:
		return "word"
	case IntArg:
		return "number"
	case FloatArg:
		return "decimal number"
	default:
		return "unknown"
	}
}

type arg struct {
	name     string
	help     string
	optional bool

	argType    ArgType
	validators []ArgValidator
}

func CreateArg(
	name string,
	help string,
	optional bool,
	argType ArgType,
	validators ...ArgValidator,
) (res *arg, err error) {
	res = &arg{
		name:       name,
		help:       help,
		argType:    argType,
		validators: validators,
	}

	return
}

func CreateStringArg(name string, help string, validators ...ArgValidator) (res *arg) {
	res = &arg{
		name:       name,
		help:       help,
		argType:    StringArg,
		validators: append([]ArgValidator{StringArgTypeValidator}, validators...),
	}

	return
}

func (arg *arg) ArgType() ArgType {
	return arg.argType
}

func (arg *arg) Name() string {
	return arg.name
}

func (arg *arg) IsOptional() bool {
	return arg.optional
}

func (arg *arg) Validate(value any) (valid bool, feedback []error) {
	feedback = []error{}

	for _, validate := range arg.validators {
		err := validate(value)

		if err != nil {
			valid = false
			feedback = append(feedback, err)
		}
	}

	return
}
