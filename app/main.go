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
		// Print prompt
		fmt.Print("$ ")

		// Read input
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		// Clean input
		line = strings.TrimSpace(line)

		// ----------------
		// exit builtin
		// ----------------
		if line == "exit" {
			return
		}

		// ----------------
		// echo builtin
		// ----------------
		if strings.HasPrefix(line, "echo") {
			if len(line) > 4 {
				fmt.Println(strings.TrimPrefix(line, "echo "))
			} else {
				fmt.Println()
			}
			continue
		}

		// ----------------
		// type builtin
		// ----------------
		if strings.HasPrefix(line, "type ") {
			arg := strings.TrimSpace(strings.TrimPrefix(line, "type "))

			// shell builtins
			if arg == "exit" || arg == "echo" || arg == "type" {
				fmt.Printf("%s is a shell builtin\n", arg)
				continue
			}

			// look up external command
			path, err := exec.LookPath(arg)
			if err != nil {
				fmt.Printf("%s not found\n", arg)
			} else {
				fmt.Printf("%s is %s\n", arg, path)
			}
			continue
		}

		// ----------------
		// external commands
		// ----------------
		if line != "" {
			parts := strings.Fields(line)

			// find executable in PATH
			path, err := exec.LookPath(parts[0])
			if err != nil {
				fmt.Printf("%s: command not found\n", parts[0])
				continue
			}

			// run the program
			cmd := exec.Command(path, parts[1:]...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// IMPORTANT: do NOT print anything on Run error
			_ = cmd.Run()
		}
	}
}
