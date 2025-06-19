package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"code.haedhutner.dev/mvv/LastMUD/internal/server"
)

func main() {
	fmt.Println(`\\\\---------------------////`)
	fmt.Println(`||||   LastMUD  Server   ||||`)
	fmt.Println(`////---------------------\\\\`)

	lastMudServer, err := server.CreateServer(":8000")

	if err != nil {
		log.Fatal(err)
	}

	go lastMudServer.Listen()

	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")

		if strings.Compare("exit", text) == 0 {
			lastMudServer.Stop()
			return
		}
	}
}
