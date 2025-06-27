package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"net/http"
	_ "net/http/pprof"

	"code.haedhutner.dev/mvv/LastMUD/internal/server"
	"golang.org/x/term"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	defer wg.Wait()
	defer cancel()

	_, err := server.CreateServer(ctx, &wg, ":8000")

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	processInput()
}

func processInput() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		panic(err)
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	buf := make([]byte, 1)

	for {
		// If interrupt received, stop
		select {
		case <-sigChan:
			return
		default:
		}

		// TODO: Proper TUI for the server
		os.Stdin.Read(buf)

		if buf[0] == 'q' {
			return
		}

		time.Sleep(50 * time.Millisecond)
	}
}
