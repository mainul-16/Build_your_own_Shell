package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func tryAutocomplete(line string) (string, bool) {
	line = strings.TrimRight(line, "\n")

	trimmed := strings.TrimRight(line, " ")
	spaces := len(line) - len(trimmed)

	if spaces == 0 {
		return "", false
	}

	if strings.HasPrefix("echo", trimmed) {
		return "echo ", true
	}
	if strings.HasPrefix("exit", trimmed) {
		return "exit ", true
	}
	return "", false
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")

		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		// ✅ AUTOCOMPLETE HANDLING
		if completed, ok := tryAutocomplete(line); ok {
			fmt.Printf("$ %s\n", completed)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if line == "exit" {
			return
		}

		if strings.HasPrefix(line, "echo ") {
			fmt.Println(strings.TrimPrefix(line, "echo "))
		}
	}
}
