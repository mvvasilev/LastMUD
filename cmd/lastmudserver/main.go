package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"net/http"
	_ "net/http/pprof"

	"code.haedhutner.dev/mvv/LastMUD/internal/server"
)

var enableDiagnostics bool = false

func main() {
	flag.BoolVar(&enableDiagnostics, "d", false, "Enable pprof server ( port :6060 ). Disabled by default.")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	defer wg.Wait()
	defer cancel()

	_, err := server.CreateServer(ctx, &wg, ":8000")

	if err != nil {
		log.Fatal(err)
	}

	if enableDiagnostics {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
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
