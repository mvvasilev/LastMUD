package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Command interface {
	Name() string
}

func main() {
	ln, err := net.Listen("tcp", ":8000")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening on port 8000")

	conn, err := ln.Accept()

	if err != nil {
		log.Fatal(err)
	}

	// cmdRegistry := commandlib.CreateCommandRegistry(
	// 	commandlib.CreateCommandDefinition(
	// 		"exit",
	// 		func(tokens []commandlib.Token) bool {
	// 			return tokens[0].Lexeme() == "exit"
	// 		},
	// 		func(tokens []commandlib.Token) []commandlib.Parameter {
	// 			return nil
	// 		},
	// 		func(parameters ...commandlib.Parameter) (err error) {
	// 			err = conn.Close()
	// 			return
	// 		},
	// 	),
	// )

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		response := ""

		// if err != nil {

		// 	if err == io.EOF {
		// 		fmt.Println("Client disconnected")
		// 		break
		// 	}

		// 	log.Println("Read error:", err)

		// 	continue
		// }

		conn.Write([]byte(message + "\n"))

		// cmdContext, err := commandlib.CreateCommandContext(cmdRegistry, message)

		// if err != nil {
		// 	log.Println(err)
		// 	response = err.Error()
		// } else {
		// 	// err = cmdContext.ExecuteCommand()

		// 	// if err != nil {
		// 	// 	log.Println(err)
		// 	// 	response = err.Error()
		// 	// }
		// }

		conn.Write([]byte(response + "\n> "))
	}

}
