package app

import (
	"errors"
	"fmt"
)

const bufferSize = 16

type Game struct {
	initialState [][]bool

	// stores next stages
	future chan [][]bool

	// buffer channel for data reuse
	buffer chan [][]bool
}

func NewGame(cells [][]bool, playOptions ...PlayOption) (*Game, error) {
	if err := checkCells(cells); err != nil {
		return nil, err
	}

	newGame := &Game{
		future:       make(chan [][]bool, 16),
		buffer:       make(chan [][]bool, 16),
		initialState: cells,
	}

	for _, opt := range playOptions {
		opt(newGame)
	}

	go newGame.begin()

	return newGame, nil
}

// Next returns the next stae in the Game.
// Returned state is not safe for use
// when this function is called again
func (g *Game) Next() [][]bool {
	nextState := <-g.future
	g.buffer <- nextState
	return nextState
}

func (g *Game) begin() {
	width := len(g.initialState)
	length := len(g.initialState[0])

	// preallocations
	allocs := make([]bool, bufferSize*width*length)

	states := make([][][]bool, bufferSize)
	for i := range states {
		states[i] = make([][]bool, width)
		gridStart := i * width * length
		for j := range states[i] {
			rowStart := gridStart + (j * length)
			states[i][j] = allocs[rowStart : rowStart+length]
		}
	}

	// use all the preallocations to calculate next N
	// where N is the bufferSize
	currentStage := g.initialState
	for _, state := range states {
		calculateNextGrid(state, currentStage)
		g.future <- state
		currentStage = state
	}

	// calculate the rest
	for {
		buffer, ok := <-g.buffer
		if !ok {
			break
		}
		calculateNextGrid(buffer, currentStage)
		g.future <- buffer
		currentStage = buffer
	}
}

// checkCells checks verifies the cells
// returns error if incorrect (e.g. not all rows are same len)
func checkCells(cells [][]bool) error {
	if len(cells) == 0 {
		return errors.New("checkCells: no rows")
	}

	if len(cells[0]) == 0 {
		return errors.New("checkCells: no data in rows")
	}

	firstElemlength := len(cells[0])
	for i, c := range cells[:len(cells)-1] {
		if len(c) != firstElemlength {
			return fmt.Errorf("length in %vth row:%v is not equal to %v", i, len(c), firstElemlength)
		}
	}

	return nil
}

func calculateNextGrid(des, src [][]bool) {
	for i, row := range src {
		for j, c := range row {
			des[i][j] = calculateDeadOrAlive(src, c, i, j)
		}
	}
}

// calculateDeadOrAlive determine if a particular cell will be dead or alive
// in the next state
func calculateDeadOrAlive(grid [][]bool, cellIsAlive bool, posX, posY int) bool {
	surroundAliveCount := calculativeSurroundingAlive(grid, posX, posY)

	if cellIsAlive {
		// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
		if surroundAliveCount < 2 {
			return false
		}

		// Any live cell with more than three live neighbours dies, as if by overpopulation.
		if surroundAliveCount > 3 {
			return false
		}

		// Any live cell with two or three live neighbours lives on to the next generation.
		return true
	}

	// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
	if surroundAliveCount == 3 {
		return true
	}
	return false
}

// calculate the number of surround alive cell
func calculativeSurroundingAlive(src [][]bool, posX, posY int) int {
	var result int

	lastXIdx := len(src) - 1
	lastYIdx := len(src[0]) - 1

	if posX > 1 { // if there is rows on above
		if src[posX-1][posY] { // top
			result++
		}

		if posY > 0 { // if there is column to the left
			if src[posX-1][posY-1] { // top left
				result++
			}
		}

		if posY < lastYIdx { // if there is column to the right
			if src[posX-1][posY+1] { // top right
				result++
			}
		}
	}

	// Middle row
	if posY > 0 { // if there is column to the left
		if src[posX][posY-1] { // bottom left
			result++
		}
	}

	if posY < lastYIdx { // if there is column to the right
		if src[posX][posY+1] { // bottom right
			result++
		}
	}

	if posX < lastXIdx { // if there is rows below
		if src[posX+1][posY] { // bottom
			result++
		}

		if posY > 0 { // if there is column to the left
			if src[posX+1][posY-1] { // bottom left
				result++
			}
		}

		if posY < lastYIdx { // if there is column to the right
			if src[posX+1][posY+1] { // bottom right
				result++
			}
		}
	}

	return result
}
