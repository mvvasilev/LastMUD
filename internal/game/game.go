package game

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/game/logic/world"
	"context"
	"sync"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/logic"
	"code.haedhutner.dev/mvv/LastMUD/internal/logging"

	"github.com/google/uuid"
)

const TickRate = 50 * time.Millisecond

type Output struct {
	connId          uuid.UUID
	contents        []byte
	closeConnection bool
}

func (g Output) Id() uuid.UUID {
	return g.connId
}

func (g Output) Contents() []byte {
	return g.contents
}

func (g Output) ShouldCloseConnection() bool {
	return g.closeConnection
}

type Game struct {
	ctx context.Context
	wg  *sync.WaitGroup

	world *World

	output chan Output
}

func CreateGame(ctx context.Context, wg *sync.WaitGroup) (game *Game) {
	game = &Game{
		wg:     wg,
		ctx:    ctx,
		output: make(chan Output),
		world:  CreateGameWorld(),
	}

	ecs.RegisterSystems(game.world.World, logic.CreateSystems()...)

	wg.Add(1)
	go game.start()

	return
}

// ConsumeNextOutput will block if no output present
func (game *Game) ConsumeNextOutput() *Output {
	select {
	case output := <-game.output:
		return &output
	default:
		return nil
	}
}

func (game *Game) ConnectPlayer(connectionId uuid.UUID) {
	world.CreatePlayerConnectEvent(game.world.World, connectionId)
}

func (game *Game) DisconnectPlayer(connectionId uuid.UUID) {
	world.CreatePlayerDisconnectEvent(game.world.World, connectionId)
}

func (game *Game) SendPlayerCommand(connectionId uuid.UUID, command string) {
	world.CreatePlayerCommandEvent(game.world.World, connectionId, command)
}

func (game *Game) start() {
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

func (game *Game) consumeOutputs() {
	entities := ecs.FindEntitiesWithComponents(game.world.World, data.TypeConnectionId, data.TypeContents)

	for _, entity := range entities {
		output := Output{}

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

func (game *Game) shutdown() {
	logging.Info("Stopping LastMUD...")
	close(game.output)
}

func (game *Game) shouldStop() bool {
	select {
	case <-game.ctx.Done():
		return true
	default:
		return false
	}
}

func (game *Game) enqeueOutput(output Output) {
	game.output <- output
}

func (game *Game) tick(delta time.Duration) {
	game.world.Tick(delta)
	game.consumeOutputs()
}
