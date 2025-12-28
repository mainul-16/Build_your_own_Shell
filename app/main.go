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

func builtinCd(args string) {

	path := args
	if len(path) == 0 {
		return
	}

	if err := os.Chdir(path); err != nil {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", path)
		return
	}
}

func main() {
	// Capture PATH once at startup (consistent behavior).
	pathEnv := os.Getenv("PATH")

	in := bufio.NewReader(os.Stdin)

	for {
		// Prompt
		fmt.Fprint(os.Stdout, "$ ")

		// Read a line
		line, err := in.ReadString('\n')
		if err != nil {
			// Handle EOF (Ctrl-D) gracefully
			if err == io.EOF {
				fmt.Println()
				return
			}
			fmt.Fprintln(os.Stderr, "read error:", err)
			continue
		}

		// Trim and split fields (handles multiple spaces/tabs)
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		cmd := fields[0]

		// Global exit
		if cmd == "exit" {
			return
		}

		switch cmd {
		case "echo":
			// echo prints remaining args joined by spaces
			fmt.Println(strings.Join(fields[1:], " "))

		case "pwd":
			dir, _ := os.Getwd()

			fmt.Println(dir)

		case "cd":
			builtinCd(fields[1])

		case "type":
			// check operands
			if len(fields) < 2 {
				fmt.Println("type: missing operand")
				continue
			}
			name := fields[1]

			// builtin?
			if builtins[name] {
				fmt.Printf("%s is a shell builtin\n", name)
				continue
			}

			// locate executable in PATH
			if full, err := locateExecutable(name, pathEnv); err == nil && full != "" {
				fmt.Printf("%s is %s\n", name, full)
			} else {
				fmt.Printf("%s: not found\n", name)
			}

		default:
			// Unknown command (same UX as your original)
			// fmt.Println(fields)
			// fmt.Println(cmd)
			full, err := locateExecutable(cmd, pathEnv)
			// fmt.Println(full)
			if err == nil && full != "" {

				p_exec := exec.Command(cmd, fields[1:]...)
				output, _ := p_exec.CombinedOutput()

				// Print the output
				fmt.Print(string(output))

			} else {
				fmt.Printf("%s: command not found\n", line)
			}
		}
	}
}

// locateExecutable looks for cmd in the PATH directories and returns
// the first executable full path it finds, or "", error if none.
func locateExecutable(cmd string, pathEnv string) (string, error) {
	paths := strings.Split(pathEnv, ":")
	for _, p := range paths {
		full := filepath.Join(p, cmd)
		if isExecutable(full) {
			return full, nil
		}
	}
	return "", fmt.Errorf("not found")
}

// isExecutable returns true if the path exists and has any exec bit set.
// Directories are considered non-executable for the purpose of locating a binary.
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	mode := info.Mode()
	return mode&0111 != 0
}
