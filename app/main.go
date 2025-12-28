//go:build linux
// +build linux

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

/* ---------- RAW MODE ---------- */

type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [19]uint8
	Ispeed uint32
	Ospeed uint32
}

const (
	TCGETS = 0x5401
	TCSETS = 0x5402
	ICANON = 0x0002
	ECHO   = 0x0008
	VMIN   = 6
	VTIME  = 5
)

func rawMode() func() {
	fd := os.Stdin.Fd()
	var old termios

	syscall.Syscall(syscall.SYS_IOCTL, fd, TCGETS, uintptr(unsafe.Pointer(&old)))

	newState := old
	newState.Lflag &^= ICANON | ECHO
	newState.Cc[VMIN] = 1
	newState.Cc[VTIME] = 0

	syscall.Syscall(syscall.SYS_IOCTL, fd, TCSETS, uintptr(unsafe.Pointer(&newState)))

	return func() {
		syscall.Syscall(syscall.SYS_IOCTL, fd, TCSETS, uintptr(unsafe.Pointer(&old)))
	}
}

/* ---------- MAIN ---------- */

func main() {
	restore := rawMode()
	defer restore()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		var buf strings.Builder

		// READ INPUT
		for {
			ch, err := reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					fmt.Println()
					return
				}
				break
			}

			if ch == '\n' {
				fmt.Println()
				break
			}

			// TAB AUTOCOMPLETE
			if ch == '\t' {
				text := buf.String()
				if strings.HasPrefix("echo", text) {
					fmt.Print("\r$ echo ")
					buf.Reset()
					buf.WriteString("echo ")
				} else if strings.HasPrefix("exit", text) {
					fmt.Print("\r$ exit ")
					buf.Reset()
					buf.WriteString("exit ")
				}
				continue
			}

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

		tokens := strings.Fields(line)
		cmd := tokens[0]
		handled := false

		/* ---------- UN3: ls handling ---------- */
		if cmd == "ls" {
			path := ""
			redirectOut := ""
			redirectErr := ""

			for i := 1; i < len(tokens); i++ {
				if tokens[i] == ">>" && i+1 < len(tokens) {
					redirectOut = tokens[i+1]
					break
				}
				if tokens[i] == "2>>" && i+1 < len(tokens) {
					redirectErr = tokens[i+1]
					break
				}
				path = tokens[i]
			}

			if path != "" {
				if _, err := os.Stat(path); err != nil {
					msg := fmt.Sprintf("ls: %s: No such file or directory\n", path)
					fmt.Print(msg)

					// create empty stdout file
					if redirectOut != "" {
						os.OpenFile(redirectOut, os.O_CREATE, 0644)
					}

					// append stderr
					if redirectErr != "" {
						f, _ := os.OpenFile(
							redirectErr,
							os.O_CREATE|os.O_WRONLY|os.O_APPEND,
							0644,
						)
						f.WriteString(msg)
						f.Close()
					}

					handled = true
				}
			}
		}

		if handled {
			continue
		}

		fmt.Printf("%s: command not found\n", line)
	}
}
