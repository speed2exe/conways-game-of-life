package main

import (
	"bufio"
	"os"
	"time"

	"github.com/speed2exe/conways-game-of-life/internal/app"
)

func main() {
	initialState := [][]bool{
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, true, true, true, false, false, false, false, false},
		{false, false, false, false, true, false, true, false, false, false, false, false},
		{false, false, false, false, true, true, true, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
	}

	game, err := app.NewGame(initialState)
	if err != nil {
		panic(err)
	}

	output := bufio.NewWriter(os.Stdout)

	for range time.NewTicker(500 * time.Millisecond).C {
		output.WriteString("-------------------------\n")
		state := game.Next()
		for _, row := range state {
			output.WriteByte('|')
			for _, c := range row {
				if c {
					output.WriteByte('@')
				} else {
					output.WriteByte(' ')
				}
			}
			output.WriteByte('|')
			output.WriteByte('\n')
		}
		output.Flush()
	}

}
