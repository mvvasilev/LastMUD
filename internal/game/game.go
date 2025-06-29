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

type InputType = string

const (
	Connect    InputType = "Connect"
	Disconnect           = "Disconnect"
	Command              = "Command"
)

type Input struct {
	connId    uuid.UUID
	inputType InputType
	command   string
}

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

	input  chan Input
	output chan Output

	stop context.CancelFunc
}

func CreateGame(ctx context.Context, wg *sync.WaitGroup) (game *Game) {
	ctx, cancel := context.WithCancel(ctx)

	game = &Game{
		wg:     wg,
		ctx:    ctx,
		input:  make(chan Input, 1000),
		output: make(chan Output, 1000),
		world:  CreateGameWorld(),
		stop:   cancel,
	}

	ecs.RegisterSystems(game.world.World, logic.CreateSystems()...)

	wg.Add(1)
	go game.start()

	return
}

// ConsumeNextOutput will block if no output present
func (game *Game) ConsumeNextOutput() *Output {
	if game.shouldStop() {
		return nil
	}

	select {
	case output := <-game.output:
		return &output
	default:
		return nil
	}
}

func (game *Game) Connect(connectionId uuid.UUID) {
	if game.shouldStop() {
		return
	}

	game.input <- Input{
		inputType: Connect,
		connId:    connectionId,
	}
}

func (game *Game) Disconnect(connectionId uuid.UUID) {
	if game.shouldStop() {
		return
	}

	game.input <- Input{
		inputType: Disconnect,
		connId:    connectionId,
	}
}

func (game *Game) SendCommand(connectionId uuid.UUID, cmd string) {
	if game.shouldStop() {
		return
	}

	game.input <- Input{
		inputType: Command,
		connId:    connectionId,
		command:   cmd,
	}
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
	entities := ecs.FindEntitiesWithComponents(game.world.World, data.TypeIsOutput, data.TypeConnectionId, data.TypeContents)

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
	close(game.input)
}

func (game *Game) shouldStop() bool {
	select {
	case <-game.ctx.Done():
		return true
	default:
		return false
	}
}

func (game *Game) nextInput() *Input {
	select {
	case input := <-game.input:
		return &input
	default:
		return nil
	}
}

func (game *Game) enqeueOutput(output Output) {
	game.output <- output
}

func (game *Game) tick(delta time.Duration) {
	for {
		input := game.nextInput()

		if input == nil {
			break
		}

		switch input.inputType {
		case Connect:
			world.CreatePlayerConnectEvent(game.world.World, input.connId)
		case Disconnect:
			world.CreatePlayerDisconnectEvent(game.world.World, input.connId)
		case Command:
			world.CreateSubmitInputEvent(game.world.World, input.connId, input.command)
		}
	}

	game.world.Tick(delta)

	game.consumeOutputs()
}
