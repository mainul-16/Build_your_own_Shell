package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var builtins = map[string]bool{
	"echo": true,
	"exit": true,
	"type": true,
	"pwd":  true,
	"cd":   true,
}

func autocomplete(line string) (string, bool) {
	line = strings.TrimSuffix(line, "\n")

	if strings.HasSuffix(line, "\t") {
		prefix := strings.TrimSuffix(line, "\t")

		if strings.HasPrefix("echo", prefix) {
			return "echo ", true
		}
		if strings.HasPrefix("exit", prefix) {
			return "exit ", true
		}
	}
	return "", false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	pathEnv := os.Getenv("PATH")

	for {
		fmt.Print("$ ")

		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		// 🔥 TAB AUTOCOMPLETE (Codecrafters-style)
		if completed, ok := autocomplete(line); ok {
			fmt.Printf("$ %s\n", completed)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		cmd := parts[0]
		args := parts[1:]

		if cmd == "exit" {
			return
		}

		switch cmd {
		case "echo":
			fmt.Println(strings.Join(args, " "))

		case "pwd":
			dir, _ := os.Getwd()
			fmt.Println(dir)

		case "cd":
			if len(args) > 0 {
				if err := os.Chdir(args[0]); err != nil {
					fmt.Printf("cd: %s: No such file or directory\n", args[0])
				}
			}

		case "type":
			if builtins[args[0]] {
				fmt.Printf("%s is a shell builtin\n", args[0])
				continue
			}
			if full, err := locateExecutable(args[0], pathEnv); err == nil {
				fmt.Printf("%s is %s\n", args[0], full)
			} else {
				fmt.Printf("%s: not found\n", args[0])
			}

		default:
			full, err := locateExecutable(cmd, pathEnv)
			if err != nil {
				fmt.Printf("%s: command not found\n", cmd)
				continue
			}
			command := &exec.Cmd{
				Path:   full,
				Args:   append([]string{cmd}, args...),
				Stdin:  os.Stdin,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}
			_ = command.Run()
		}
	}
}

func locateExecutable(cmd, pathEnv string) (string, error) {
	for _, dir := range strings.Split(pathEnv, ":") {
		full := filepath.Join(dir, cmd)
		if info, err := os.Stat(full); err == nil && info.Mode()&0111 != 0 {
			return full, nil
		}
	}
	return "", fmt.Errorf("not found")
}
