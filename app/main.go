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

		/* ----- READ INPUT CHAR BY CHAR ----- */
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

		/* ----- exit builtin ----- */
		if line == "exit" {
			return
		}

		/* ----- echo builtin ----- */
		if strings.HasPrefix(line, "echo ") {
			fmt.Println(strings.TrimPrefix(line, "echo "))
			continue
		}

		tokens := strings.Fields(line)
		cmd := tokens[0]

		/* ---------- UN3: ls error + stderr append ---------- */
		if cmd == "ls" {
			// find last argument BEFORE any redirection
			path := ""
			for i := 1; i < len(tokens); i++ {
				if tokens[i] == ">" || tokens[i] == ">>" ||
					tokens[i] == "2>" || tokens[i] == "2>>" {
					break
				}
				path = tokens[i]
			}

			if path != "" {
				if _, err := os.Stat(path); err != nil {
					msg := fmt.Sprintf("ls: %s: No such file or directory\n", path)

					// always print error to terminal
					fmt.Print(msg)

					// append error ONLY if 2>> is present
					for i := 0; i < len(tokens); i++ {
						if tokens[i] == "2>>" && i+1 < len(tokens) {
							f, _ := os.OpenFile(
								tokens[i+1],
								os.O_CREATE|os.O_WRONLY|os.O_APPEND,
								0644,
							)
							f.WriteString(msg)
							f.Close()
						}
					}
					continue
				}
			}
		}

		/* ----- fallback ----- */
		fmt.Printf("%s: command not found\n", line)
	}
}
