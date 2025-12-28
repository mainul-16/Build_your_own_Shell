package main

import (
	"bufio"
	"fmt"
	"io"
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

func builtinCd(args []string) {
	if len(args) == 0 {
		return
	}

	path := args[0]

	if path == "~" {
		home := os.Getenv("HOME")
		if home == "" {
			fmt.Printf("cd: ~: No such file or directory\n")
			return
		}
		path = home
	}

	if err := os.Chdir(path); err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", args[0])
	}
}

func main() {
	pathEnv := os.Getenv("PATH")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				return
			}
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		cmd := fields[0]
		args := fields[1:]

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
			builtinCd(args)

		case "type":
			if len(args) == 0 {
				continue
			}

			name := args[0]

			if builtins[name] {
				fmt.Printf("%s is a shell builtin\n", name)
				continue
			}

			if full, err := locateExecutable(name, pathEnv); err == nil {
				fmt.Printf("%s is %s\n", name, full)
			} else {
				fmt.Printf("%s: not found\n", name)
			}

		default:
			full, err := locateExecutable(cmd, pathEnv)
			if err != nil {
				fmt.Printf("%s: command not found\n", cmd)
				continue
			}

			command := &exec.Cmd{
				Path:   full,
				Args:   append([]string{cmd}, args...), // 🔥 FIX HERE
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
		if isExecutable(full) {
			return full, nil
		}
	}
	return "", fmt.Errorf("not found")
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	return info.Mode()&0111 != 0
}
