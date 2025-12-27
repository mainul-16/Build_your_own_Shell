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

			if arg == "exit" || arg == "echo" || arg == "type" {
				fmt.Printf("%s is a shell builtin\n", arg)
				continue
			}

			path, err := exec.LookPath(arg)
			if err != nil {
				fmt.Printf("%s not found\n", arg)
			} else {
				fmt.Printf("%s is %s\n", arg, path)
			}
			continue
		}

		// -----------------------------
		// external commands (FIXED)
		// -----------------------------
		if line != "" {
			parts := strings.Fields(line)

			path, err := exec.LookPath(parts[0])
			if err != nil {
				fmt.Printf("%s: command not found\n", parts[0])
				continue
			}

			cmd := &exec.Cmd{
				Path:   path,
				Args:   append([]string{parts[0]}, parts[1:]...), // argv[0] FIXED
				Stdin:  os.Stdin,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}

			_ = cmd.Run()
		}
	}
}
