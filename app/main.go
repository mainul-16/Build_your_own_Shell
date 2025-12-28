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

/* ---------------- PARSER ---------------- */

func parseCommand(line string) []string {
	var args []string
	var current strings.Builder

	inSingle := false
	inDouble := false

	for i := 0; i < len(line); i++ {
		ch := line[i]

		// Backslash handling
		if ch == '\\' {
			if !inSingle && !inDouble {
				if i+1 < len(line) {
					current.WriteByte(line[i+1])
					i++
				}
				continue
			}

			if inDouble {
				if i+1 < len(line) {
					next := line[i+1]
					if next == '"' || next == '\\' {
						current.WriteByte(next)
						i++
					} else {
						current.WriteByte('\\')
					}
				} else {
					current.WriteByte('\\')
				}
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

/* ---------------- MAIN ---------------- */

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

		fields := parseCommand(line)
		if len(fields) == 0 {
			continue
		}

		// --- Handle output redirection ---
		var outFile *os.File
		cleanArgs := []string{}

		for i := 0; i < len(fields); i++ {
			if fields[i] == ">" || fields[i] == "1>" {
				if i+1 < len(fields) {
					f, err := os.Create(fields[i+1])
					if err != nil {
						fmt.Println("error creating file")
						continue
					}
					outFile = f
					i++ // skip filename
				}
			} else {
				cleanArgs = append(cleanArgs, fields[i])
			}
		}

		cmd := cleanArgs[0]
		args := cleanArgs[1:]

		if cmd == "exit" {
			return
		}

		// Default stdout
		stdout := os.Stdout
		if outFile != nil {
			stdout = outFile
		}

		switch cmd {

		case "echo":
			fmt.Fprintln(stdout, strings.Join(args, " "))

		case "pwd":
			dir, _ := os.Getwd()
			fmt.Fprintln(stdout, dir)

		case "cd":
			builtinCd(args)

		case "type":
			if len(args) == 0 {
				continue
			}
			name := args[0]

			if builtins[name] {
				fmt.Fprintf(stdout, "%s is a shell builtin\n", name)
				continue
			}

			if full, err := locateExecutable(name, pathEnv); err == nil {
				fmt.Fprintf(stdout, "%s is %s\n", name, full)
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
				Args:   append([]string{cmd}, args...),
				Stdin:  os.Stdin,
				Stdout: stdout,
				Stderr: os.Stderr,
			}
			_ = command.Run()
		}

		if outFile != nil {
			outFile.Close()
		}
	}
}

/* ---------------- HELPERS ---------------- */

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
