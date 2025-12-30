package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var builtins = []string{
	"echo",
	"exit",
	"type",
	"pwd",
	"cd",
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		var buf []rune

		for {
			ch, err := reader.ReadByte()
			if err != nil {
				return
			}

			// ENTER
			if ch == '\n' {
				fmt.Println()
				break
			}

			// TAB
			if ch == '\t' {
				input := string(buf)
				matches := []string{}

				for _, b := range builtins {
					if strings.HasPrefix(b, input) {
						matches = append(matches, b)
					}
				}

				if len(matches) == 1 {
					rest := matches[0][len(input):]
					for _, r := range rest {
						buf = append(buf, r)
						fmt.Printf("%c", r)
					}
					buf = append(buf, ' ')
					fmt.Print(" ")
				} else {
					// 🔔 IMPORTANT FIX: real bell
					os.Stdout.Write([]byte{7})
				}
				continue
			}

			// BACKSPACE
			if ch == 127 {
				if len(buf) > 0 {
					buf = buf[:len(buf)-1]
					fmt.Print("\b \b")
				}
				continue
			}

			// NORMAL CHAR
			buf = append(buf, rune(ch))
			fmt.Printf("%c", ch)
		}

		line := strings.TrimSpace(string(buf))
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

		fmt.Println(line + ": command not found")
	}
}

