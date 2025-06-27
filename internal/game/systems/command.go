package systems

import (
	"fmt"
	"regexp"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
)

type CommandHandler func(world *ecs.World, delta time.Duration, player ecs.Entity, args map[string]data.Arg) (err error)

func commandQuery(command data.Command) func(comp data.CommandComponent) bool {
	return func(comp data.CommandComponent) bool {
		return comp.Cmd == command
	}
}

func CreateCommandHandler(command data.Command, handler CommandHandler) ecs.SystemExecutor {
	return func(world *ecs.World, delta time.Duration) (err error) {
		commands := ecs.QueryEntitiesWithComponent(world, commandQuery(command))
		processedCommands := []ecs.Entity{}

		for c := range commands {
			logging.Debug("Handling command of type ", command)

			player, _ := ecs.GetComponent[data.PlayerComponent](world, c)
			args, _ := ecs.GetComponent[data.ArgsComponent](world, c)

			err := handler(world, delta, player.Player, args.Args)

			if err != nil {
				logging.Info("Issue while handling command ", command, ": ", err)

				connId, _ := ecs.GetComponent[data.ConnectionIdComponent](world, player.Player)

				data.CreateGameOutput(world, connId.ConnectionId, []byte(err.Error()), false)
			}

			ecs.SetComponent(world, c, data.CommandStateComponent{State: data.CommandStateExecuted})

			// data.CreateCommandExecutedEvent(world, c) // Not needed right now

			processedCommands = append(processedCommands, c)
		}

		ecs.DeleteEntities(world, processedCommands...)

		return
	}
}

type commandError struct {
	err string
}

func createCommandError(v ...any) *commandError {
	return &commandError{
		err: fmt.Sprint("Error handling command: ", v),
	}
}

func (e *commandError) Error() string {
	return e.err
}

func handlePlayerCommandEvent(world *ecs.World, event ecs.Entity) (err error) {
	commandString, ok := ecs.GetComponent[data.CommandStringComponent](world, event)

	if !ok {
		return createCommandError("Unable to handle command, no command string found for event")
	}

	eventConnId, ok := ecs.GetComponent[data.ConnectionIdComponent](world, event)

	if !ok {
		return createCommandError("Unable to handle command, no connection id found for event")
	}

	player := ecs.NilEntity()

	for p := range ecs.IterateEntitiesWithComponent[data.IsPlayerComponent](world) {
		playerConnId, ok := ecs.GetComponent[data.ConnectionIdComponent](world, p)

		if ok && playerConnId.ConnectionId == eventConnId.ConnectionId {
			player = p
			break
		}
	}

	if player == ecs.NilEntity() {
		return createCommandError("Unable to find valid player with provided connection id")
	}

	tokens, err := tokenize(commandString.Command)

	if err != nil {
		return createCommandError("Error with tokenization: ", err)
	}

	command := data.CreateTokenizedCommand(world, player, commandString.Command, tokens)
	data.CreateParseCommandEvent(world, command)

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

func parseCommand(world *ecs.World, event ecs.Entity) (err error) {
	command, ok := ecs.GetComponent[data.EntityComponent](world, event)

	if !ok {
		return createCommandError("Unable to parse command: no command entity provided in event")
	}

	tokens, ok := ecs.GetComponent[data.TokensComponent](world, command.Entity)

	if !ok {
		return createCommandError("Unable to parse command: no tokens provided in command entity")
	}

	var foundMatch bool

	for cmd, parser := range commandParsers {
		match, args := parser(tokens.Tokens)

		if !match {
			continue
		}

		ecs.SetComponent(world, command.Entity, data.ArgsComponent{Args: args})
		ecs.SetComponent(world, command.Entity, data.CommandComponent{Cmd: cmd})
		ecs.SetComponent(world, command.Entity, data.CommandStateComponent{State: data.CommandStateParsed})

		foundMatch = true
	}

	player, _ := ecs.GetComponent[data.PlayerComponent](world, command.Entity)
	connectionId, _ := ecs.GetComponent[data.ConnectionIdComponent](world, player.Player)

	if !foundMatch {
		data.CreateGameOutput(world, connectionId.ConnectionId, []byte("Unknown command"), false)
		ecs.DeleteEntity(world, command.Entity)
	}

	return
}
