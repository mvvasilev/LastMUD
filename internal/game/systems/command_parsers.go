package systems

import (
	"strings"

	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
)

func oneOf(value string, tests ...string) bool {
	for _, t := range tests {
		if value == t {
			return true
		}
	}

	return false
}

type commandParser = func(tokens []data.Token) (matches bool, args map[string]data.Arg)

var commandParsers = map[data.Command]commandParser{
	data.CommandSay: func(tokens []data.Token) (matches bool, args map[string]data.Arg) {
		matches = len(tokens) > 1
		matches = matches && oneOf(tokens[0].Lexeme, "say", "lc", "localchat")

		if !matches {
			return
		}

		lexemes := []string{}

		for _, t := range tokens[1:] {
			lexemes = append(lexemes, t.Lexeme)
		}

		args = map[string]data.Arg{
			"messageContent": {Value: strings.Join(lexemes, " ")},
		}

		return
	},
	data.CommandQuit: func(tokens []data.Token) (matches bool, args map[string]data.Arg) {
		matches = len(tokens) >= 1
		matches = matches && oneOf(tokens[0].Lexeme, "quit", "disconnect", "q", "leave")

		return
	},
}
