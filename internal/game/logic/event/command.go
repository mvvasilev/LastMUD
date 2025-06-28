package event

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/logic/command"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/logic/world"
	"fmt"
	"regexp"
)

type commandParseError struct {
	err string
}

func createCommandParseError(v ...any) *commandParseError {
	return &commandParseError{
		err: fmt.Sprint("Error parsing command: ", v),
	}
}

func (e *commandParseError) Error() string {
	return e.err
}

func HandlePlayerCommand(w *ecs.World, event ecs.Entity) (err error) {
	commandString, ok := ecs.GetComponent[data.CommandStringComponent](w, event)

	if !ok {
		return createCommandParseError("No command string found for event")
	}

	eventConnId, ok := ecs.GetComponent[data.ConnectionIdComponent](w, event)

	if !ok {
		return createCommandParseError("No connection id found for event")
	}

	player := ecs.NilEntity()

	for p := range ecs.IterateEntitiesWithComponent[data.IsPlayerComponent](w) {
		playerConnId, ok := ecs.GetComponent[data.ConnectionIdComponent](w, p)

		if ok && playerConnId.ConnectionId == eventConnId.ConnectionId {
			player = p
			break
		}
	}

	if player == ecs.NilEntity() {
		return createCommandParseError("Unable to find valid player with provided connection id")
	}

	tokens, err := tokenize(commandString.Command)

	if err != nil {
		return createCommandParseError("Error with tokenization: ", err)
	}

	cmd := world.CreateTokenizedCommand(w, player, commandString.Command, tokens)
	world.CreateParseCommandEvent(w, cmd)

	return
}

func tokenize(commandString string) (tokens []data.Token, err error) {
	tokens = []data.Token{}
	pos := 0
	inputLen := len(commandString)

	// Continue iterating until we reach the end of the input
	for pos < inputLen {
		matched := false
		remaining := commandString[pos:]

		// Iterate through each token type and test its pattern
		for tokenType, pattern := range data.TokenPatterns {
			// If the token pattern doesn't compile, panic ( why do we have invalid patterns?! )
			tokenPattern := regexp.MustCompile(pattern)

			// If the location of the match isn't nil, that means we've found a match
			if loc := tokenPattern.FindStringIndex(remaining); loc != nil {
				lexeme := remaining[loc[0]:loc[1]]

				pos += loc[1]
				matched = true

				// Skip whitespace
				if tokenType == data.TokenWhitespace {
					break
				}

				tokens = append(
					tokens,
					data.Token{
						Type:   tokenType,
						Lexeme: lexeme,
						Index:  pos,
					},
				)

				break
			}
		}

		// Unknown tokens are still added
		if !matched {
			tokens = append(
				tokens,
				data.Token{
					Type:   data.TokenUnknown,
					Lexeme: commandString[pos : pos+1],
					Index:  pos,
				},
			)

			pos++
		}
	}

	// Mark the end of the tokens
	tokens = append(tokens, data.Token{Type: data.TokenEOF, Lexeme: "", Index: pos})

	return
}

func HandleParseCommand(w *ecs.World, event ecs.Entity) (err error) {
	cmdEnt, ok := ecs.GetComponent[data.EntityComponent](w, event)

	if !ok {
		return createCommandParseError("No command entity provided in event")
	}

	tokens, ok := ecs.GetComponent[data.TokensComponent](w, cmdEnt.Entity)

	if !ok {
		return createCommandParseError("No tokens provided in command entity")
	}

	var foundMatch bool

	for cmd, parser := range command.Parsers {
		match, args := parser(tokens.Tokens)

		if !match {
			continue
		}

		ecs.SetComponent(w, cmdEnt.Entity, data.ArgsComponent{Args: args})
		ecs.SetComponent(w, cmdEnt.Entity, data.CommandComponent{Cmd: cmd})
		ecs.SetComponent(w, cmdEnt.Entity, data.CommandStateComponent{State: data.CommandStateParsed})

		foundMatch = true
	}

	player, _ := ecs.GetComponent[data.PlayerComponent](w, cmdEnt.Entity)
	connectionId, _ := ecs.GetComponent[data.ConnectionIdComponent](w, player.Player)

	if !foundMatch {
		world.CreateGameOutput(w, connectionId.ConnectionId, []byte("Unknown command"))
		ecs.DeleteEntity(w, cmdEnt.Entity)
	}

	return
}
