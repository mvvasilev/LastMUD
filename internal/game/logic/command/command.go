package command

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/game/logic/world"
	"fmt"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
)

type commandError struct {
	err string
}

func createCommandError(v ...any) *commandError {
	return &commandError{
		err: fmt.Sprint("Command error: ", v),
	}
}

func (e *commandError) Error() string {
	return e.err
}

type Handler func(w *ecs.World, delta time.Duration, player ecs.Entity, args map[string]data.Arg) (err error)

func commandQuery(command data.Command) func(comp data.CommandComponent) bool {
	return func(comp data.CommandComponent) bool {
		return comp.Cmd == command
	}
}

func CreateHandler(command data.Command, handler Handler) ecs.SystemExecutor {
	return func(w *ecs.World, delta time.Duration) (err error) {
		commands := ecs.QueryEntitiesWithComponent(w, commandQuery(command))
		var processedCommands []ecs.Entity

		for c := range commands {
			logging.Debug("Handling command of type ", command)

			player, _ := ecs.GetComponent[data.PlayerComponent](w, c)
			args, _ := ecs.GetComponent[data.ArgsComponent](w, c)

			err := handler(w, delta, player.Player, args.Args)

			if err != nil {
				logging.Info("Issue while handling command ", command, ": ", err)

				connId, _ := ecs.GetComponent[data.ConnectionIdComponent](w, player.Player)

				world.CreateGameOutput(w, connId.ConnectionId, err.Error())
			}

			ecs.SetComponent(w, c, data.CommandStateComponent{State: data.CommandStateExecuted})

			// data.CreateCommandExecutedEvent(world, c) // Not needed right now

			processedCommands = append(processedCommands, c)
		}

		ecs.DeleteEntities(w, processedCommands...)

		return
	}
}

func arg[T any](args data.ArgsMap, name data.ArgName) (val T, err error) {
	uncast, ok := args[name]

	if !ok || uncast.Value == nil {
		err = createCommandError("No arg named '", name, "' found or value is empty")
		return
	}

	val, ok = uncast.Value.(T)

	if !ok {
		err = createCommandError("Arg value of '", name, "' cannot be cast")
		return
	}

	return
}
