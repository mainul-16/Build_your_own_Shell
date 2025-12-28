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

/* ---------------- AUTOCOMPLETE ---------------- */

func autocomplete(prefix string) (string, bool) {
	if strings.HasPrefix("echo", prefix) {
		return "echo ", true
	}
	if strings.HasPrefix("exit", prefix) {
		return "exit ", true
	}
	return "", false
}

/* ---------------- PARSER ---------------- */

func parseCommand(line string) []string {
	var args []string
	var current strings.Builder
	inSingle, inDouble := false, false

	for i := 0; i < len(line); i++ {
		ch := line[i]

		if ch == '\\' {
			if !inSingle && !inDouble && i+1 < len(line) {
				current.WriteByte(line[i+1])
				i++
				continue
			}
			current.WriteByte('\\')
			continue
		}

		switch ch {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			} else {
				current.WriteByte(ch)
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			} else {
				current.WriteByte(ch)
			}
		case ' ', '\t':
			if inSingle || inDouble {
				current.WriteByte(ch)
			} else if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(ch)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

/* ---------------- BUILTINS ---------------- */

func builtinCd(args []string) {
	if len(args) == 0 {
		return
	}
	path := args[0]
	if path == "~" {
		path = os.Getenv("HOME")
	}
	if err := os.Chdir(path); err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", args[0])
	}
}

/* ---------------- MAIN ---------------- */

func main() {
	reader := bufio.NewReader(os.Stdin)
	pathEnv := os.Getenv("PATH")

	for {
		fmt.Print("$ ")

		buffer := ""

		for {
			ch, _, err := reader.ReadRune()
			if err != nil {
				return
			}

			// ENTER
			if ch == '\n' {
				fmt.Print("\n")
				break
			}

			// TAB → autocomplete
			if ch == '\t' {
				if completed, ok := autocomplete(buffer); ok {
					// 🔥 clear line + redraw
					fmt.Print("\r\033[K$ ")
					fmt.Print(completed)
					buffer = strings.TrimSpace(completed)
				}
				continue
			}

			// Normal character
			buffer += string(ch)
			fmt.Print(string(ch))
		}

		line := strings.TrimSpace(buffer)
		if line == "" {
			continue
		}

		fields := parseCommand(line)
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

/* ---------------- HELPERS ---------------- */

func locateExecutable(cmd, pathEnv string) (string, error) {
	for _, dir := range strings.Split(pathEnv, ":") {
		full := filepath.Join(dir, cmd)
		if info, err := os.Stat(full); err == nil && info.Mode()&0111 != 0 {
			return full, nil
		}
	}
	return "", fmt.Errorf("not found")
}
