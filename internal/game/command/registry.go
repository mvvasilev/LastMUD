package command

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
)

type TokenMatcher func(tokens []Token) bool

type ParameterParser func(tokens []Token) []Parameter

type CommandWork func(parameters ...Parameter) (err error)

type CommandDefinition struct {
	name            string
	tokenMatcher    TokenMatcher
	parameterParser ParameterParser
}

func CreateCommandDefinition(
	name string,
	tokenMatcher TokenMatcher,
	parameterParser ParameterParser,
) CommandDefinition {
	return CommandDefinition{
		name:            name,
		tokenMatcher:    tokenMatcher,
		parameterParser: parameterParser,
	}
}

func (def CommandDefinition) Name() string {
	return def.name
}

func (def CommandDefinition) Match(tokens []Token) bool {
	return def.tokenMatcher(tokens)
}

func (def CommandDefinition) ParseParameters(tokens []Token) []Parameter {
	return def.parameterParser(tokens)
}

type CommandRegistry struct {
	commandDefinitions []CommandDefinition
}

func CreateCommandRegistry(commandDefinitions ...CommandDefinition) *CommandRegistry {
	return &CommandRegistry{
		commandDefinitions: commandDefinitions,
	}
}

func (comReg *CommandRegistry) Match(tokens []Token) (comDef *CommandDefinition) {
	for _, v := range comReg.commandDefinitions {
		if v.Match(tokens) {
			logging.Debug("Found match", v.Name())
			return &v
		}
	}

	return
}
