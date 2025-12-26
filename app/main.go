package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		// 1. Display prompt
		fmt.Print("$ ")

		// 2. Read user input
		command, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		// 3. Remove newline
		command = command[:len(command)-1]

		// 4. Print error message
		fmt.Println(command + ": command not found")
	}
}
