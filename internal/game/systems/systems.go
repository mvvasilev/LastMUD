package systems

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
)

type SystemPriorityOffset = int

const (
	EventOffset   = 0
	CommandOffset = 10000
)

func CreateSystems() []*ecs.System {
	return []*ecs.System{
		ecs.CreateSystem("PlayerConnectEventHandler", EventOffset+0, CreateEventHandler(data.EventPlayerConnect, handlePlayerConnectEvent)),
		ecs.CreateSystem("PlayerDisconnectEventHandler", EventOffset+1, CreateEventHandler(data.EventPlayerDisconnect, handlePlayerDisconnectEvent)),
		ecs.CreateSystem("PlayerCommandEventHandler", EventOffset+2, CreateEventHandler(data.EventPlayerCommand, handlePlayerCommandEvent)),
		ecs.CreateSystem("ParseCommandEventHandler", EventOffset+4, CreateEventHandler(data.EventParseCommand, parseCommand)),
		ecs.CreateSystem("SayCommandHandler", CommandOffset+0, CreateCommandHandler(data.CommandSay, handleSayCommand)),
		ecs.CreateSystem("QuitCommandHandler", CommandOffset+1, CreateCommandHandler(data.CommandQuit, handleQuitCommand)),
	}
}
