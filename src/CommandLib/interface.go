package commandlib

type Parameter interface {
	Value() any
}

type Command interface {
	Name() string
	Parameters() []Parameter
}
