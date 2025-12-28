package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var builtins = map[string]bool{
	"exit": true,
	"echo": true,
	"type": true,
	"pwd":  true,
	"cd":   true,
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")

		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		cmd := parts[0]
		args := parts[1:]

		// exit
		if cmd == "exit" {
			return
		}

		// cd (absolute paths)
		if cmd == "cd" {
			if len(args) == 0 {
				continue
			}

			path := args[0]
			if err := os.Chdir(path); err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", path)
			}
			continue
		}

		// pwd
		if cmd == "pwd" {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("pwd: error retrieving current directory")
			} else {
				fmt.Println(cwd)
			}
			continue
		}

		// echo
		if cmd == "echo" {
			fmt.Println(strings.Join(args, " "))
			continue
		}

		// type
		if cmd == "type" {
			if len(args) == 0 {
				continue
			}

			target := args[0]

			if builtins[target] {
				fmt.Printf("%s is a shell builtin\n", target)
				continue
			}

			path, err := exec.LookPath(target)
			if err != nil {
				fmt.Printf("%s not found\n", target)
			} else {
				fmt.Printf("%s is %s\n", target, path)
			}
			continue
		}

		// external commands
		path, err := exec.LookPath(cmd)
		if err != nil {
			fmt.Printf("%s: command not found\n", cmd)
			continue
		}

		command := &exec.Cmd{
			Path:   path,
			Args:   append([]string{cmd}, args...),
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}

		_ = command.Run()
	}
}
