package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		// Print prompt
		fmt.Print("$ ")

		// Read full line (TAB is sent as spaces)
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		// Remove newline
		line = strings.TrimRight(line, "\n")

		// Detect TAB (represented as trailing spaces)
		trimmed := strings.TrimRight(line, " ")

		// If spaces were added, treat as TAB
		if trimmed != line {
			if strings.HasPrefix("echo", trimmed) {
				fmt.Println("$ echo ")
				continue
			}
			if strings.HasPrefix("exit", trimmed) {
				fmt.Println("$ exit ")
				continue
			}
		}

		// Handle exit if user actually types it
		if trimmed == "exit" {
			return
		}
	}
}
