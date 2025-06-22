package game

import (
	"context"
	"sync"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/game/command"
	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
	"github.com/google/uuid"
)

const TickRate = time.Duration(50 * time.Millisecond)

type GameOutput struct {
	connId   uuid.UUID
	contents []byte
}

func (game *LastMUDGame) CreateOutput(connId uuid.UUID, contents []byte) GameOutput {
	return GameOutput{
		connId:   connId,
		contents: contents,
	}
}

func (g GameOutput) Id() uuid.UUID {
	return g.connId
}

func (g GameOutput) Contents() []byte {
	return g.contents
}

type LastMUDGame struct {
	ctx context.Context
	wg  *sync.WaitGroup

	commandRegistry *command.CommandRegistry
	world           *World

	eventBus *EventBus

	output chan GameOutput
}

func CreateGame(ctx context.Context, wg *sync.WaitGroup) (game *LastMUDGame) {
	game = &LastMUDGame{
		wg:       wg,
		ctx:      ctx,
		eventBus: CreateEventBus(),
		output:   make(chan GameOutput, 10),
		world:    CreateWorld(),
	}

	game.commandRegistry = game.CreateGameCommandRegistry()

	wg.Add(1)
	go game.start()

	return
}

func (game *LastMUDGame) start() {
	defer game.wg.Done()
	defer game.shutdown()

	logging.Info("Starting LastMUD...")

	lastTick := time.Now()

	for {
		now := time.Now()

		if game.shouldStop() {
			break
		}

		game.tick(now.Sub(lastTick))

		// Tick at regular intervals
		if time.Since(lastTick) < TickRate {
			time.Sleep(TickRate - time.Since(lastTick))
		}

		lastTick = now
	}
}

func (game *LastMUDGame) shutdown() {
	logging.Info("Stopping LastMUD...")
	close(game.output)
	game.eventBus.close()
}

func (game *LastMUDGame) shouldStop() bool {
	select {
	case <-game.ctx.Done():
		return true
	default:
		return false
	}
}

func (game *LastMUDGame) EnqueueEvent(event GameEvent) {
	game.eventBus.Push(event)
}

func (game *LastMUDGame) enqeueOutput(output GameOutput) {
	game.output <- output
}

func (game *LastMUDGame) ConsumeNextOutput() *GameOutput {
	select {
	case output := <-game.output:
		return &output
	default:
		return nil
	}
}

func (game *LastMUDGame) CommandRegistry() *command.CommandRegistry {
	return game.commandRegistry
}

func (g *LastMUDGame) tick(delta time.Duration) {
	for {
		event := g.eventBus.Pop()

		if event == nil {
			return
		}

		event.Handle(g, delta)
	}
}
