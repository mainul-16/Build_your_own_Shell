package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	builtins := []string{"echo", "exit", "type", "pwd", "cd"}

	rl, err := readline.NewEx(&readline.Config{
		Prompt: "$ ",
		AutoComplete: readline.AutoCompleterFunc(
			func(line []rune, pos int) ([][]rune, int) {
				input := string(line[:pos])
				matches := [][]rune{}

				for _, b := range builtins {
					if strings.HasPrefix(b, input) {
						matches = append(matches, []rune(b[len(input):]+" "))
					}
				}

				// ✅ NO MATCH → ring bell, return nothing
				if len(matches) == 0 {
					readline.Bell()
					return nil, 0
				}

				// 1 match → autocomplete
				if len(matches) == 1 {
					return matches, len(input)
				}

				// >1 match → do nothing (future stage)
				return nil, 0
			},
		),
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
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

