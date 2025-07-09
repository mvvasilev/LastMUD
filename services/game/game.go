package game

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/server"
	"code.haedhutner.dev/mvv/LastMUD/shared/log"
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func LaunchGameServer(enablePprof bool) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	log.DefaultLogger = log.NewLogger("lastmud", log.DBG, log.DBG, "./log/lastmud.log")

	defer wg.Wait()
	defer cancel()

	_, err := server.CreateServer(ctx, &wg, ":8000")

	if err != nil {
		log.Fatal(err)
	}

	if enablePprof {
		go func() {
			log.Error(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	processInput()
}

func processInput() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		// If interrupt received, stop
		select {
		case <-sigChan:
			return
		default:
		}

		time.Sleep(50 * time.Millisecond)
	}
}
