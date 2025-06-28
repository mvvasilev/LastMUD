package logic

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/logic/command"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/logic/event"
)

const (
	EventOffset   = 0
	CommandOffset = 10000
)

func CreateSystems() []*ecs.System {
	return []*ecs.System{
		// Event Handlers
		ecs.CreateSystem("PlayerConnectEventHandler", EventOffset+0, event.CreateHandler(data.EventPlayerConnect, event.HandlePlayerConnect)),
		ecs.CreateSystem("PlayerDisconnectEventHandler", EventOffset+1, event.CreateHandler(data.EventPlayerDisconnect, event.HandlePlayerDisconnect)),
		ecs.CreateSystem("PlayerCommandEventHandler", EventOffset+2, event.CreateHandler(data.EventPlayerCommand, event.HandlePlayerCommand)),
		ecs.CreateSystem("ParseCommandEventHandler", EventOffset+4, event.CreateHandler(data.EventParseCommand, event.HandleParseCommand)),

		// Command Handlers
		ecs.CreateSystem("SayCommandHandler", CommandOffset+0, command.CreateHandler(data.CommandSay, command.HandleSay)),
		ecs.CreateSystem("QuitCommandHandler", CommandOffset+1, command.CreateHandler(data.CommandQuit, command.HandleQuit)),
		ecs.CreateSystem("RegisterCommandHandler", CommandOffset+2, command.CreateHandler(data.CommandRegister, command.HandleRegister)),
	}
}
