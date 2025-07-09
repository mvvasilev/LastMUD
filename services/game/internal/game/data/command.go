package data

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
)

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
	TokenUnknown             TokenType = "UNKOWN"
	TokenNumber              TokenType = "NUMBER"
	TokenDecimal             TokenType = "DECIMAL"
	TokenIdentifier          TokenType = "IDENTIFIER"
	TokenBracketedIdentifier TokenType = "BRACKET_IDENTIFER"
	TokenWhitespace          TokenType = "WHITESPACE"
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

type ArgName = string

const (
	ArgMessageContent  ArgName = "messageContent"
	ArgAccountName     ArgName = "accountName"
	ArgAccountPassword ArgName = "accountPassword"
)

type Arg struct {
	Value any
}

type ArgsMap = map[ArgName]Arg

type ArgsComponent struct {
	Args ArgsMap
}

func (ac ArgsComponent) Type() ecs.ComponentType {
	return TypeCommandArgs
}

type Command string

const (
	CommandNoOp     Command = "noop"
	CommandSay              = "say"
	CommandQuit             = "quit"
	CommandLogin            = "login"
	CommandRegister         = "register"
)

type CommandComponent struct {
	Cmd Command
}

func (cc CommandComponent) Type() ecs.ComponentType {
	return TypeCommand
}
