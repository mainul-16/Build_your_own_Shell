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

/* ---- RAW MODE (ABSOLUTE MINIMUM) ---- */

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

/* ---- MAIN ---- */

func main() {
	restore := rawMode()
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
