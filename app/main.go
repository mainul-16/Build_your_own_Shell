package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

		// type builtin
		if strings.HasPrefix(line, "type ") {
			arg := strings.TrimSpace(strings.TrimPrefix(line, "type "))

			// shell builtins
			if arg == "exit" || arg == "echo" || arg == "type" {
				fmt.Printf("%s is a shell builtin\n", arg)
				continue
			}

			// external command lookup
			path, err := exec.LookPath(arg)
			if err == nil {
				fmt.Printf("%s is %s\n", arg, path)
			} else {
				fmt.Printf("%s not found\n", arg)
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
