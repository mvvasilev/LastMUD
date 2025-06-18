package commandlib

import "strconv"

type Parameter struct {
	value string
}

func CreateParameter(value string) Parameter {
	return Parameter{
		value: value,
	}
}

func (p Parameter) AsString() (res string, err error) {
	return p.value, nil
}

func (p Parameter) AsInteger() (res int, err error) {
	return strconv.Atoi(p.value)
}

func (p Parameter) AsDecimal() (res float64, err error) {
	return strconv.ParseFloat(p.value, 32)
}
