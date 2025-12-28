package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

/* ---------------- TERMINAL RAW MODE ---------------- */

func enableRawMode() func() {
	fd := int(os.Stdin.Fd())

	oldState, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return func() {}
	}

	newState := *oldState
	newState.Lflag &^= unix.ICANON | unix.ECHO
	newState.Cc[unix.VMIN] = 1
	newState.Cc[unix.VTIME] = 0

	_ = unix.IoctlSetTermios(fd, unix.TCSETS, &newState)

	return func() {
		_ = unix.IoctlSetTermios(fd, unix.TCSETS, oldState)
	}
}

/* ---------------- AUTOCOMPLETE ---------------- */

func autocomplete(input string) (string, bool) {
	builtins := []string{"echo", "exit"}

	for _, b := range builtins {
		if strings.HasPrefix(b, input) {
			return b + " ", true
		}
	}
	return input, false
}

/* ---------------- BUILTINS ---------------- */

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

/* ---------------- MAIN ---------------- */

func main() {
	restore := enableRawMode()
	defer restore()

	pathEnv := os.Getenv("PATH")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")

		var line strings.Builder

		// ----- READ INPUT BYTE BY BYTE -----
		for {
			ch, err := reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					fmt.Println()
					return
				}
				continue
			}

			// ENTER
			if ch == '\n' {
				fmt.Println()
				break
			}

			// TAB
			if ch == '\t' {
				current := line.String()
				completed, ok := autocomplete(current)
				if ok {
					fmt.Print("\r$ ")
					fmt.Print(completed)
					line.Reset()
					line.WriteString(completed)
				}
				continue
			}

			// NORMAL CHARACTER
			line.WriteByte(ch)
			fmt.Printf("%c", ch)
		}

		commandLine := strings.TrimSpace(line.String())
		if commandLine == "" {
			continue
		}

		fields := strings.Fields(commandLine)
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
