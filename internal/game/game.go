package game

import (
	"context"
	"sync"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/systems"
	"code.haedhutner.dev/mvv/LastMUD/internal/logging"

	"github.com/google/uuid"
)

const TickRate = time.Duration(50 * time.Millisecond)

const MaxEnqueuedOutputPerTick = 100

type GameOutput struct {
	connId          uuid.UUID
	contents        []byte
	closeConnection bool
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

func (g GameOutput) ShouldCloseConnection() bool {
	return g.closeConnection
}

type LastMUDGame struct {
	ctx context.Context
	wg  *sync.WaitGroup

	world *data.GameWorld

	output chan GameOutput
}

func CreateGame(ctx context.Context, wg *sync.WaitGroup) (game *LastMUDGame) {
	game = &LastMUDGame{
		wg:     wg,
		ctx:    ctx,
		output: make(chan GameOutput, MaxEnqueuedOutputPerTick),
		world:  data.CreateGameWorld(),
	}

	ecs.RegisterSystems(game.world.World, systems.CreateSystems()...)

	wg.Add(1)
	go game.start()

	return
}

func (game *LastMUDGame) ConsumeNextOutput() *GameOutput {
	select {
	case output := <-game.output:
		return &output
	default:
		return nil
	}
}

func (game *LastMUDGame) ConnectPlayer(connectionId uuid.UUID) {
	data.CreatePlayerConnectEvent(game.world.World, connectionId)
}

func (game *LastMUDGame) DisconnectPlayer(connectionId uuid.UUID) {
	data.CreatePlayerDisconnectEvent(game.world.World, connectionId)
}

func (game *LastMUDGame) SendPlayerCommand(connectionId uuid.UUID, command string) {
	data.CreatePlayerCommandEvent(game.world.World, connectionId, command)
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

func (game *LastMUDGame) consumeOutputs() {
	entities := ecs.FindEntitiesWithComponents(game.world.World, data.TypeConnectionId, data.TypeContents)

	for _, entity := range entities {
		output := GameOutput{}

		connId, _ := ecs.GetComponent[data.ConnectionIdComponent](game.world.World, entity)
		output.connId = connId.ConnectionId

		contents, hasContents := ecs.GetComponent[data.ContentsComponent](game.world.World, entity)

		if hasContents {
			output.contents = contents.Contents
		}

		_, shouldClose := ecs.GetComponent[data.CloseConnectionComponent](game.world.World, entity)
		output.closeConnection = shouldClose

		game.enqeueOutput(output)
	}

	ecs.DeleteEntities(game.world.World, entities...)
}

func (game *LastMUDGame) shutdown() {
	logging.Info("Stopping LastMUD...")
	close(game.output)
}

func (game *LastMUDGame) shouldStop() bool {
	select {
	case <-game.ctx.Done():
		return true
	default:
		return false
	}
}

func (game *LastMUDGame) enqeueOutput(output GameOutput) {
	game.output <- output
}

func (g *LastMUDGame) tick(delta time.Duration) {
	g.world.Tick(delta)
	g.consumeOutputs()
}
