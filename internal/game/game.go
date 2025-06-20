package game

import (
	"context"
	"sync"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
)

const TickRate = time.Duration(50 * time.Millisecond)

type GameSignal struct {
}

type LastMUDGame struct {
	ctx context.Context
	wg  *sync.WaitGroup
}

func CreateGame(ctx context.Context, wg *sync.WaitGroup) (game *LastMUDGame) {
	game = &LastMUDGame{
		wg:  wg,
		ctx: ctx,
	}

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
}

func (game *LastMUDGame) shouldStop() bool {
	select {
	case <-game.ctx.Done():
		return true
	default:
		return false
	}
}

func (g *LastMUDGame) tick(delta time.Duration) {
	// logging.Debug("Tick")
	// TODO
}
