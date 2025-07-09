package event

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/game/data"
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/game/logic/command"
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/game/logic/world"
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

func HandleSubmitInput(w *ecs.World, event ecs.Entity) (err error) {
	commandString, ok := ecs.GetComponent[data.CommandStringComponent](w, event)

	if !ok {
		return createCommandParseError("No command string found for event")
	}

	eventConnId, ok := ecs.GetComponent[data.ConnectionIdComponent](w, event)

	if !ok {
		return createCommandParseError("No connection id found for event")
	}

	player := world.FindPlayerByConnectionId(w, eventConnId.ConnectionId)

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
	cmdEnt := ecs.QueryFirstEntityWithComponent(w, world.QueryParentEvent(event))

	if cmdEnt == ecs.NilEntity() {
		return createCommandParseError("No command entity provided in event")
	}

	tokens, ok := ecs.GetComponent[data.TokensComponent](w, cmdEnt)

	if !ok {
		return createCommandParseError("No tokens provided in command entity")
	}

	cmd, args, ok := command.ParseCommand(tokens.Tokens...)

	if !ok {
		player, _ := ecs.GetComponent[data.ParentComponent](w, cmdEnt)
		connectionId, _ := ecs.GetComponent[data.ConnectionIdComponent](w, player.Entity)

		world.CreateGameOutput(w, connectionId.ConnectionId, "Unknown command")
		ecs.DeleteEntity(w, cmdEnt)
		return createCommandParseError("Not a valid command")
	}

	ecs.SetComponent(w, cmdEnt, data.ArgsComponent{Args: args})
	ecs.SetComponent(w, cmdEnt, data.CommandComponent{Cmd: cmd})
	ecs.SetComponent(w, cmdEnt, data.CommandStateComponent{State: data.CommandStateParsed})

	return
}
