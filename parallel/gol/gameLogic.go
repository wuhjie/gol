package gol

import (
	"strconv"
	"strings"
	"uk.ac.bris.cs/gameoflife/util"
)

//calculation
func mod(x, m int) int {
	return (x + m) % m
}

func initialisedWorld(height, width int) [][]byte {
	world := make([][]byte, height)
	for i:= range world {
		world[i] = make([]byte, width)
	}
	return world
}

// making immutable world for calculating, prevent race condition
func makeImmutableWorld(world [][]byte) func(y, x int) byte {
	return func(y, x int) byte {
		return world[y][x]
	}
}

//used to calculate the alive neighbors
func calculateNeighbors(p Params, x, y int, world func(y, x int) byte) int {
	neighbors := 0
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i != 0 || j != 0 {
				if world(mod(y+i, p.ImageHeight),mod(x+j, p.ImageWidth)) == alive {
					neighbors++
				}
			}
		}
	}
	return neighbors
}

// calculate the world after changing
func calculateNextStage(startY, endY, startX, endX int, p Params, world func(y, x int) byte, c distributorChannels) [][]byte {
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
			neighbors := calculateNeighbors(p, x, absoluteY, world)
			if world(absoluteY, x) == alive {
				if neighbors == 2 || neighbors == 3 {
					newWorld[y][x] = alive
				} else {
					newWorld[y][x] = dead
					c.events <- CellFlipped{CompletedTurns: c.completedTurns, Cell: util.Cell{X: x, Y: absoluteY}}
				}
			}
			if world(absoluteY, x) == dead {
				if neighbors == 3 {
					newWorld[y][x] = alive
					c.events <- CellFlipped{CompletedTurns: c.completedTurns, Cell: util.Cell{X: x, Y: absoluteY}}
				} else {
					newWorld[y][x] = dead
				}
			}
		}
	}
	return newWorld
}

//calculate the alive cells in current round
func calculateAliveCells(p Params, world [][]byte) []util.Cell {
	var aliveCells []util.Cell

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			if world[y][x] == alive {
				aliveCells = append(aliveCells, util.Cell{X: x, Y: y})
			}
		}
	}
	return aliveCells
}

// for sending the world into the ioOutput channel
func outputWorldImage(c distributorChannels, p Params, world [][]byte) {
	c.ioCommand <- ioOutput
	filename := strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight), strconv.Itoa(c.completedTurns)}, "x")
	c.ioFilename <- filename

	for m := 0; m < p.ImageHeight; m++ {
		for n := 0; n < p.ImageWidth; n++ {
			c.ioOutput <- world[m][n]
		}
	}
	c.events <- ImageOutputComplete{c.completedTurns, filename}
}