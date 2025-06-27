package data

import "code.haedhutner.dev/mvv/LastMUD/internal/ecs"

type CommandStringComponent struct {
	Command string
}

func (cs CommandStringComponent) Type() ecs.ComponentType {
	return TypeCommandString
}

type CommandState byte

const (
	CommandStateTokenized CommandState = iota
	CommandStateParsed
	CommandStateExecuted
)

type CommandStateComponent struct {
	State CommandState
}

func (cs CommandStateComponent) Type() ecs.ComponentType {
	return TypeCommandState
}

type TokenType string

const (
	TokenEOF                 TokenType = "EOF"
	TokenUnknown                       = "UNKOWN"
	TokenNumber                        = "NUMBER"
	TokenDecimal                       = "DECIMAL"
	TokenIdentifier                    = "IDENTIFIER"
	TokenBracketedIdentifier           = "BRACKET_IDENTIFER"
	TokenWhitespace                    = "WHITESPACE"
)

var TokenPatterns = map[TokenType]string{
	TokenDecimal:             `(?i)^\b\d+\.\d+\b`,
	TokenNumber:              `(?i)^\b\d+\b`,
	TokenBracketedIdentifier: `\[[^\]]+\]`,
	TokenIdentifier:          `[^\s]+`,
	TokenWhitespace:          `(?i)^\s+`,
	TokenUnknown:             `(?i)^[^ \t\n\r\f\v]+`,
}

type Token struct {
	Type   TokenType
	Lexeme string
	Index  int
}

type TokensComponent struct {
	Tokens []Token
}

func (tc TokensComponent) Type() ecs.ComponentType {
	return TypeCommandTokens
}

type Arg struct {
	Value any
}

type ArgsComponent struct {
	Args map[string]Arg
}

func (ac ArgsComponent) Type() ecs.ComponentType {
	return TypeCommandArgs
}

type Command string

const (
	CommandSay  Command = "say"
	CommandQuit         = "quit"
)

type CommandComponent struct {
	Cmd Command
}

func (cc CommandComponent) Type() ecs.ComponentType {
	return TypeCommand
}

func CreateTokenizedCommand(world *ecs.World, player ecs.Entity, commandString string, tokens []Token) ecs.Entity {
	command := ecs.NewEntity()

	ecs.SetComponent(world, command, PlayerComponent{Player: player})
	ecs.SetComponent(world, command, CommandStringComponent{Command: commandString})
	ecs.SetComponent(world, command, TokensComponent{Tokens: tokens})
	ecs.SetComponent(world, command, CommandStateComponent{State: CommandStateTokenized})

	return command
}
