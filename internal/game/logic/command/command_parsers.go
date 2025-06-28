package command

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

type commandParser = func(tokens []data.Token) (matches bool, args data.ArgsMap)

var Parsers = map[data.Command]commandParser{
	data.CommandSay: func(tokens []data.Token) (matches bool, args data.ArgsMap) {
		matches = len(tokens) > 1
		matches = matches && oneOf(tokens[0].Lexeme, "say", "lc", "localchat")

		if !matches {
			return
		}

		var lexemes []string

		for _, t := range tokens[1:] {
			lexemes = append(lexemes, t.Lexeme)
		}

		args = data.ArgsMap{
			data.ArgMessageContent: {Value: strings.Join(lexemes, " ")},
		}

		return
	},
	data.CommandQuit: func(tokens []data.Token) (matches bool, args data.ArgsMap) {
		matches = len(tokens) >= 1
		matches = matches && oneOf(tokens[0].Lexeme, "quit", "disconnect", "q", "leave")

		return
	},
	data.CommandRegister: func(tokens []data.Token) (matches bool, args data.ArgsMap) {
		matches = len(tokens) >= 3
		matches = matches && oneOf(tokens[0].Lexeme, "register", "signup")

		if !matches {
			return
		}

		accountName := tokens[1].Lexeme
		accountPassword := tokens[2].Lexeme

		args = data.ArgsMap{
			data.ArgAccountName:     {Value: accountName},
			data.ArgAccountPassword: {Value: accountPassword},
		}

		return
	},
	data.CommandLogin: func(tokens []data.Token) (matches bool, args data.ArgsMap) {
		matches = len(tokens) >= 3
		matches = matches && oneOf(tokens[0].Lexeme, "login", "signin")

		if !matches {
			return
		}

		accountName := tokens[1].Lexeme
		accountPassword := tokens[2].Lexeme

		args = data.ArgsMap{
			data.ArgAccountName:     {Value: accountName},
			data.ArgAccountPassword: {Value: accountPassword},
		}

		return
	},
}
