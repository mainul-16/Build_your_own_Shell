package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Print shell prompt
	fmt.Print("$ ")

	// Read user input
	command, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	// Remove newline and print error
	fmt.Println(command[:len(command)-1] + ": command not found")
}
