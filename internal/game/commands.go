package game

import "code.haedhutner.dev/mvv/LastMUD/internal/game/command"

type CommandType = string

const (
	SayCommand CommandType = "say"
)

func (game *LastMUDGame) CreateGameCommandRegistry() *command.CommandRegistry {
	return command.CreateCommandRegistry(
		command.CreateCommandDefinition(
			SayCommand,
			func(tokens []command.Token) bool {
				return len(tokens) > 1 && tokens[0].Lexeme() == "say"
			},
			func(tokens []command.Token) []command.Parameter {
				saying := ""

				for _, t := range tokens[1:] {
					saying += t.Lexeme()
				}

				return []command.Parameter{
					command.CreateParameter(saying),
				}
			},
		),
	)
}
