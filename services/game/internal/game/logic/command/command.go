package command

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/game/data"
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/game/logic/world"
	"code.haedhutner.dev/mvv/LastMUD/shared/log"
	"fmt"
	"time"
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

func CommandSystem(w *ecs.World, delta time.Duration) (err error) {
	commands := ecs.FindEntitiesWithComponents(w, data.TypeCommand, data.TypeParent, data.TypeCommandArgs, data.TypeCommandState)

	var processedCommands []ecs.Entity

	for _, c := range commands {
		command, _ := ecs.GetComponent[data.CommandComponent](w, c)
		player, _ := ecs.GetComponent[data.ParentComponent](w, c)
		args, _ := ecs.GetComponent[data.ArgsComponent](w, c)
		state, _ := ecs.GetComponent[data.CommandStateComponent](w, c)

		if state.State != data.CommandStateParsed {
			continue
		}

		if player.Entity == ecs.NilEntity() {
			return createCommandError("Nil entity as command parent")
		}

		handler, ok := handlers[command.Cmd]

		if !ok {
			return createCommandError("No handler available for command of type ", command.Cmd)
		}

		err := handler(w, delta, player.Entity, args.Args)

		if err != nil {
			log.Info("Issue while handling command ", command, ": ", err)

			connId, _ := ecs.GetComponent[data.ConnectionIdComponent](w, player.Entity)

			world.CreateGameOutput(w, connId.ConnectionId, err.Error())
		}

		ecs.SetComponent(w, c, data.CommandStateComponent{State: data.CommandStateExecuted})

		processedCommands = append(processedCommands, c)
	}

	ecs.DeleteEntities(w, processedCommands...)
	return
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
