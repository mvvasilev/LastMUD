package commandlib

import "strings"

type TokenType = byte

const (
	TokenEOF = iota

	TokenUnknown

	TokenNumber
	TokenDecimal
	TokenIdentifier
	TokenBracketedIdentifier

	TokenDirection
	TokenCommand
	TokenSelf

	TokenPunctuation
)

var tokenPatterns = map[TokenType]string{
	TokenNumber:              `\b\d+\b`,
	TokenDecimal:             `\b\d+\.\d+\b`,
	TokenIdentifier:          `\b[a-zA-Z][a-zA-Z0-9]*\b`,
	TokenBracketedIdentifier: `\[[a-zA-Z][a-zA-Z0-9]*\]`,
	TokenDirection:           `\b(north|south|east|west|up|down)\b`,
	TokenSelf:                `\bself\b`,
	TokenPunctuation:         `[,.!?'/":;\-\[\]\(\)]`,
	TokenUnknown:             `.`,
}

type Token struct {
	token  TokenType
	lexeme string
	index  int
}

func CreateToken(token TokenType, lexeme string, index int) Token {
	return Token{
		token:  token,
		lexeme: lexeme,
		index:  index,
	}
}

func (t Token) Token() TokenType {
	return t.token
}

func (t Token) Lexeme() string {
	return t.lexeme
}

func (t Token) Index() int {
	return t.index
}

type tokenizer struct {
	commandNameTokenRegex string
}

func CreateTokenizer(commandNames []string) *tokenizer {
	return &tokenizer{
		commandNameTokenRegex: `\b(` + strings.Join(commandNames, "|") + `)\b`,
	}
}

func (t *tokenizer) Tokenize(commandMsg string) (tokens []Token) {
	tokens = []Token{}
	pos := 0
	inputLen := len(commandMsg)

	for pos < inputLen {
		matched := false

		for tokenType, pattern := range tokenPatterns {

		}
	}
}
