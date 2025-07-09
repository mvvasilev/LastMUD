package main

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game"
	"flag"
	_ "net/http/pprof"
)

var enableDiagnostics bool = false

func main() {
	flag.BoolVar(&enableDiagnostics, "d", false, "Enable pprof server ( port :6060 ). Disabled by default.")

	flag.Parse()

	game.LaunchGameServer(enableDiagnostics)
}
