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
		fmt.Print("$ ")

		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)

		// exit builtin
		if line == "exit" {
			return
		}

		// echo builtin
		if strings.HasPrefix(line, "echo") {
			if len(line) > 4 {
				fmt.Println(strings.TrimPrefix(line, "echo "))
			} else {
				fmt.Println()
			}
			continue
		}

		// unknown command
		if line != "" {
			cmd := strings.Fields(line)[0]
			fmt.Printf("%s: command not found\n", cmd)
		}
	}
}
