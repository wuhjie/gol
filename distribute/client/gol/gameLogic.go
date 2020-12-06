package gol

import (
	util2 "gol/distribute/client/util"
	"strconv"
	"strings"
)

//calculation-related
func mod(x, m int) int {
	return (x + m) % m
}

// initialising new 2-Dimension array of byte
func InitialisedWorld(height, width int) [][]byte {
	world := make([][]byte, height)
	for i := range world {
		world[i] = make([]byte, width)
	}
	return world
}

// making immutable world for calculating, prevent race condition
func MakeImmutableWorld(world [][]byte) func(y, x int) byte {
	return func(y, x int) byte {
		return world[y][x]
	}
}

//used to calculate the alive neighbors
func CalculateNeighbors(p Params, x, y int, world func(y, x int) byte) int {
	neighbors := 0
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i != 0 || j != 0 {
				if world(mod(y+i, p.ImageHeight), mod(x+j, p.ImageWidth)) == alive {
					neighbors++
				}
			}
		}
	}
	return neighbors
}

// Calculate the world after changing
func CalculateNextStage(startY, endY, startX, endX int, p Params, world func(y, x int) byte, c DistributorChannels) [][]byte {
	newWorld := make([][]byte, endY-startY)

	// width and height of current piece
	height := endY - startY
	width := endX - startX

	for i := range newWorld {
		newWorld[i] = make([]byte, p.ImageWidth)
	}

	// calculate world in current piece; the cell need to compare with the cell in the original world
	for y := 0; y < height; y++ {
		absoluteY := y + startY

		for x := 0; x < width; x++ {
			neighbors := CalculateNeighbors(p, x, absoluteY, world)
			if world(absoluteY, x) == alive {
				if neighbors == 2 || neighbors == 3 {
					newWorld[y][x] = alive
				} else {
					newWorld[y][x] = dead
					c.Events <- CellFlipped{CompletedTurns: c.CompletedTurns, Cell: util2.Cell{X: x, Y: absoluteY}}
				}
			}
			if world(absoluteY, x) == dead {
				if neighbors == 3 {
					newWorld[y][x] = alive
					c.Events <- CellFlipped{CompletedTurns: c.CompletedTurns, Cell: util2.Cell{X: x, Y: absoluteY}}
				} else {
					newWorld[y][x] = dead
				}
			}
		}
	}
	return newWorld
}

//calculate the alive cells in current round
func CalculateAliveCells(p Params, world [][]byte) []util2.Cell {
	var aliveCells []util2.Cell

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			if world[y][x] == alive {
				aliveCells = append(aliveCells, util2.Cell{X: x, Y: y})
			}
		}
	}
	return aliveCells
}

// for sending the world into the IoOutput channel
func OutputWorldImage(c DistributorChannels, p Params, world [][]byte) {
	c.IoCommand <- ioOutput
	filename := strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight), strconv.Itoa(c.CompletedTurns)}, "x")
	c.IoFilename <- filename

	for m := 0; m < p.ImageHeight; m++ {
		for n := 0; n < p.ImageWidth; n++ {
			c.IoOutput <- world[m][n]
		}
	}
	c.Events <- ImageOutputComplete{c.CompletedTurns, filename}

}
