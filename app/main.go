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

/* ---------- RAW MODE (LINUX ONLY) ---------- */

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

		/* ----- READ INPUT ----- */
		for {
			ch, err := reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					fmt.Println()
					return
				}
				continue
			}

			if ch == '\n' {
				fmt.Println()
				break
			}

			// TAB AUTOCOMPLETE (QP2)
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

		/* ---------- UN3: APPEND STDERR ---------- */

		tokens := strings.Fields(line)
		for i := 0; i < len(tokens); i++ {
			if tokens[i] == "2>>" && i+1 < len(tokens) {
				cmd := tokens[0]
				outFile := tokens[i+1]

				if cmd == "ls" {
					path := tokens[i-1]

					if _, err := os.Stat(path); err != nil {
						msg := fmt.Sprintf("ls: %s: No such file or directory\n", path)

						// print error to terminal
						fmt.Print(msg)

						// append error to file
						f, _ := os.OpenFile(
							outFile,
							os.O_CREATE|os.O_WRONLY|os.O_APPEND,
							0644,
						)
						f.WriteString(msg)
						f.Close()
						goto PROMPT
					}
				}
			}
		}

		fmt.Printf("%s: command not found\n", line)

	PROMPT:
		continue
	}
}
