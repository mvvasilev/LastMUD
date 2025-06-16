package commandlib

type Command interface {
	Name() string
	DoWork(argValues []ArgumentValue) (err error)
}

type commandRegistry struct {
	commands []Command
}
