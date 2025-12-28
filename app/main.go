package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
)

/* -------- RAW MODE (stdlib only) -------- */

func enableRawMode() func() {
	fd := int(os.Stdin.Fd())

	oldState, err := syscall.IoctlGetTermios(fd, syscall.TCGETS)
	if err != nil {
		return func() {}
	}

	newState := *oldState
	newState.Lflag &^= syscall.ICANON | syscall.ECHO
	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0

	_ = syscall.IoctlSetTermios(fd, syscall.TCSETS, &newState)

	return func() {
		_ = syscall.IoctlSetTermios(fd, syscall.TCSETS, oldState)
	}
}

/* -------- AUTOCOMPLETE -------- */

func autocomplete(s string) (string, bool) {
	if strings.HasPrefix("echo", s) {
		return "echo ", true
	}
	if strings.HasPrefix("exit", s) {
		return "exit ", true
	}
	return s, false
}

/* -------- MAIN -------- */

func main() {
	restore := enableRawMode()
	defer restore()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")

		var buf strings.Builder

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
				text := buf.String()
				if completed, ok := autocomplete(text); ok {
					fmt.Print("\r$ ")
					fmt.Print(completed)
					buf.Reset()
					buf.WriteString(completed)
				}
				continue
			}

			// NORMAL CHAR
			buf.WriteByte(ch)
			fmt.Printf("%c", ch)
		}

		line := strings.TrimSpace(buf.String())
		if line == "" {
			continue
		}

		if line == "exit" {
			return
		}

		if strings.HasPrefix(line, "echo ") {
			fmt.Println(strings.TrimPrefix(line, "echo "))
			continue
		}

		fmt.Printf("%s: command not found\n", line)
	}
}
