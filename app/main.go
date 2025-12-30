package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

type BuiltinCompleter struct {
	builtins []string
}

func (c *BuiltinCompleter) Do(line []rune, pos int) ([][]rune, int) {
	prefix := string(line[:pos])
	matches := [][]rune{}

	for _, b := range c.builtins {
		if strings.HasPrefix(b, prefix) {
			matches = append(matches, []rune(b[len(prefix):]+" "))
		}
	}

	// ❗ no match → ring bell, do nothing
	if len(matches) == 0 {
		os.Stdout.Write([]byte("\x07")) // bell
		return nil, 0
	}

	// single match → autocomplete
	if len(matches) == 1 {
		return matches, len(prefix)
	}

	return nil, 0
}

func main() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt: "$ ",
		AutoComplete: &BuiltinCompleter{
			builtins: []string{"echo", "exit", "type", "pwd", "cd"},
		},
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		}
		if err == io.EOF {
			return
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

