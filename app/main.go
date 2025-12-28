package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func autocomplete(input string) string {
	if strings.HasPrefix("echo", input) {
		return "echo "
	}
	if strings.HasPrefix("exit", input) {
		return "exit "
	}
	return input
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		var input strings.Builder

		for {
			ch, _ := reader.ReadByte()

			// ENTER
			if ch == '\n' {
				fmt.Println()
				break
			}

			// TAB
			if ch == '\t' {
				completed := autocomplete(input.String())
				fmt.Print("\r$ " + completed)
				input.Reset()
				input.WriteString(completed)
				continue
			}

			input.WriteByte(ch)
			fmt.Print(string(ch))
		}

		cmd := strings.TrimSpace(input.String())

		if cmd == "exit" {
			return
		}

		if strings.HasPrefix(cmd, "echo ") {
			fmt.Println(strings.TrimPrefix(cmd, "echo "))
			continue
		}

		fmt.Println(cmd + ": command not found")
	}
}
 