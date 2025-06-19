package command

import "log"

type TokenMatcher func(tokens []Token) bool

type ParameterParser func(tokens []Token) []Parameter

type CommandWork func(parameters ...Parameter) (err error)

type CommandDefinition struct {
	name            string
	tokenMatcher    TokenMatcher
	parameterParser ParameterParser
	work            CommandWork
}

func CreateCommandDefinition(
	name string,
	tokenMatcher TokenMatcher,
	parameterParser ParameterParser,
	work CommandWork,
) CommandDefinition {
	return CommandDefinition{
		name:            name,
		tokenMatcher:    tokenMatcher,
		parameterParser: parameterParser,
		work:            work,
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

func (def CommandDefinition) ExecuteFunc() CommandWork {
	return def.work
}

type CommandRegistry struct {
	commandDefinitions []CommandDefinition
}

func CreateCommandRegistry(commandDefinitions ...CommandDefinition) *CommandRegistry {
	return &CommandRegistry{
		commandDefinitions: commandDefinitions,
	}
}

func (comReg *CommandRegistry) Register(newCommandDefinitions ...CommandDefinition) {
	comReg.commandDefinitions = append(comReg.commandDefinitions, newCommandDefinitions...)
}

func (comReg *CommandRegistry) Match(tokens []Token) (comDef *CommandDefinition) {
	for _, v := range comReg.commandDefinitions {
		if v.Match(tokens) {
			log.Println("Found match", v.Name())
			return &v
		}
	}

	return
}
