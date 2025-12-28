package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")

		buffer := ""
		spaceCount := 0

		for {
			ch, err := reader.ReadByte()
			if err != nil {
				return
			}

			// ENTER → stop input
			if ch == '\n' {
				fmt.Print("\n")
				break
			}

			// SPACE (TAB is rendered as spaces)
			if ch == ' ' {
				spaceCount++

				// FIRST space after builtin prefix = autocomplete trigger
				if spaceCount == 1 {
					if buffer == "ech" {
						fmt.Print("\r$ echo ")
						buffer = "echo"
						continue
					}
					if buffer == "exi" {
						fmt.Print("\r$ exit ")
						buffer = "exit"
						continue
					}
				}

				// otherwise just echo space
				fmt.Print(" ")
				continue
			}

			// reset space counter on normal character
			spaceCount = 0
			buffer += string(ch)
			fmt.Print(string(ch))
		}

		// handle exit
		if buffer == "exit" {
			return
		}
	}
}
