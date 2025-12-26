package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Read user input
	command, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	// Remove trailing newline and print error message
	fmt.Println(command[:len(command)-1] + ": command not found")
}
