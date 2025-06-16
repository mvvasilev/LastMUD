package commandlib

type argValue struct {
	value any
}

func CreateArgValue(val any) *argValue {
	return &argValue{
		value: val,
	}
}

func (aVal *argValue) Value() any {
	return aVal.value
}
