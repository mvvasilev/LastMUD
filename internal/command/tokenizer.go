package command

import (
	"fmt"
	"regexp"
)

type TokenType byte

const (
	TokenEOF TokenType = iota

	TokenUnknown

	TokenNumber
	TokenDecimal
	TokenIdentifier
	TokenBracketedIdentifier
	TokenText

	TokenDirection
	TokenCommand
	TokenSelf

	TokenWhitespace
)

func (tt TokenType) String() string {
	switch tt {
	case TokenEOF:
		return "EOF"
	case TokenUnknown:
		return "Unknown"
	case TokenNumber:
		return "Number"
	case TokenDecimal:
		return "Decimal"
	case TokenIdentifier:
		return "Identifier"
	case TokenBracketedIdentifier:
		return "BracketedIdentifier"
	case TokenText:
		return "Text"
	case TokenDirection:
		return "Direction"
	case TokenCommand:
		return "Command"
	case TokenSelf:
		return "Self"
	case TokenWhitespace:
		return "Whitespace"
	default:
		return fmt.Sprintf("TokenType(%d)", byte(tt))
	}
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

func (t Token) String() string {
	return fmt.Sprintf("%3d %16v: %q", t.index, t.token, t.lexeme)
}

type tokenPattern struct {
	tokenType TokenType
	pattern   string
}

// Used to tokenize a string input.
// This is the starting point for parsing a command string.
// Create with [CreateTokenizer]
type Tokenizer struct {
	tokenPatterns []tokenPattern
}

func CreateTokenizer() *Tokenizer {
	return &Tokenizer{
		tokenPatterns: []tokenPattern{
			{tokenType: TokenDecimal, pattern: `\b\d+\.\d+\b`},
			{tokenType: TokenNumber, pattern: `\b\d+\b`},
			{tokenType: TokenDirection, pattern: `\b(north|south|east|west|up|down)\b`},
			{tokenType: TokenBracketedIdentifier, pattern: `\[[ a-zA-Z0-9'-][ a-zA-Z0-9'-]*\]`},
			{tokenType: TokenSelf, pattern: `\bself\b`},
			{tokenType: TokenIdentifier, pattern: `\b[a-zA-Z'-][a-zA-Z0-9'-]*\b`},
			{tokenType: TokenWhitespace, pattern: `\s+`},
			{tokenType: TokenUnknown, pattern: `.`},
		},
	}
}

// Tokenize a command string
func (t *Tokenizer) Tokenize(commandString string) (tokens []Token, err error) {
	tokens = []Token{}
	pos := 0
	inputLen := len(commandString)

	// Continue iterating until we reach the end of the input
	for pos < inputLen {
		matched := false
		remaining := commandString[pos:]

		// Iterate through each token type and test its pattern
		for _, pattern := range t.tokenPatterns {
			// All patterns are case-insensitive and must match the beginning of the input (^)
			re, regexError := regexp.Compile(`(?i)^` + pattern.pattern)

			// If we encounter a regex error, stop tokenization ( wrongly defined pattern? )
			if regexError != nil {
				tokens = nil
				err = regexError
				return
			}

			// If the location of the match isn't nil, that means we've found a match
			if loc := re.FindStringIndex(remaining); loc != nil {
				lexeme := remaining[loc[0]:loc[1]]

				pos += loc[1]
				matched = true

				tokens = append(tokens, CreateToken(pattern.tokenType, lexeme, pos))
				break
			}
		}

		// Unknown tokens are still added
		if !matched {
			tokens = append(tokens, CreateToken(TokenUnknown, commandString[pos:pos+1], pos))
			pos++
		}
	}

	// Mark the end of the tokens
	tokens = append(tokens, CreateToken(TokenEOF, "", pos))

	return
}
