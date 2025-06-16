package commandlib

import (
	"strconv"
	"strings"
)

func Tokenize(commandMsg string) []any {
	split := strings.Split(commandMsg, " ")

	tokens := []any{}

	for _, v := range split {
		valInt, err := strconv.ParseInt(v, 10, 32)

		if err == nil {
			tokens = append(tokens, valInt)
		}

		valFloat, err := strconv.ParseFloat(v, 32)

		if err == nil {
			tokens = append(tokens, valFloat)
		}

		tokens = append(tokens, v)
	}

	return tokens
}
