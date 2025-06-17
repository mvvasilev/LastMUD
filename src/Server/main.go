package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	commandlib "code.haedhutner.dev/mvv/LastMUD/CommandLib"
)

type Command interface {
	Name() string
}

type argValue struct {
	value string
}

func main() {
	// testcmd, err := commandlib.CreateCommand(
	// 	"test",
	// 	"t",
	// 	func(argValues []commandlib.ArgumentValue) (err error) {
	// 		err = nil
	// 		return
	// 	},
	// 	commandlib.CreateStringArg("test", "test message"),
	// )

	tokenizer := commandlib.CreateTokenizer()

	ln, err := net.Listen("tcp", ":8000")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening on port 8000")

	conn, err := ln.Accept()

	if err != nil {
		log.Fatal(err)
	}

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		response := ""

		if err != nil {
			log.Fatal(err)
		}

		conn.Write([]byte(message + "\n"))

		tokens, err := tokenizer.Tokenize(message)

		if err != nil {
			response = err.Error()
		} else {
			lines := make([]string, len(tokens))

			for i, tok := range tokens {
				lines[i] = tok.String()
			}

			response = strings.Join(lines, "\n")
		}

		// if strings.HasPrefix(message, testcmd.Name()) {
		// 	tokens := commandlib.Tokenize(message)
		// 	args := []commandlib.ArgumentValue{}

		// 	for _, v := range tokens[1:] {
		// 		args = append(args, commandlib.CreateArgValue(v))
		// 	}

		// 	err := testcmd.DoWork(args)

		// 	if err != nil {
		// 		fmt.Print(err.Error())
		// 	}
		// } else {
		// 	fmt.Print("Message Received: ", string(message))

		// 	response = strings.ToUpper(message)
		// }

		conn.Write([]byte(response + "\n> "))
	}

}
