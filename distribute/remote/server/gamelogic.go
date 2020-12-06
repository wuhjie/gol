package server

import (
	"strconv"
	"strings"

	"uk.ac.bris.cs/gameoflife/remoteutil"
)

const alive = 255
const dead = 0

//calculation-related
func mod(x, m int) int {
	return (x + m) % m
}

// InitialisedWorld initialising new 2-Dimension array of byte
func InitialisedWorld(height, width int) [][]byte {
	world := make([][]byte, height)
	for i := range world {
		world[i] = make([]byte, width)
	}
	return world
}

// MakeImmutableWorld makes immutable world for calculating, prevent race condition
func MakeImmutableWorld(world [][]byte) func(y, x int) byte {
	return func(y, x int) byte {
		return world[y][x]
	}
}

//CalculateNeighbors is used to calculate the alive neighbors
func CalculateNeighbors(p remoteutil.Params, x, y int, world func(y, x int) byte) int {
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

// CalculateNextStage calculates the world after changing
func CalculateNextStage(startY, endY, startX, endX int, p remoteutil.Params, world func(y, x int) byte) [][]byte {
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
				}
			}
			if world(absoluteY, x) == dead {
				if neighbors == 3 {
					newWorld[y][x] = alive
				} else {
					newWorld[y][x] = dead
				}
			}
		}
	}
	return newWorld
}

//CalculateAliveCells calculates the alive cells in current round
func CalculateAliveCells(p remoteutil.Params, world [][]byte) []remoteutil.Cell {
	var aliveCells []remoteutil.Cell

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			if world[y][x] == alive {
				aliveCells = append(aliveCells, remoteutil.Cell{X: x, Y: y})
			}
		}
	}
	return aliveCells
}

// OutputWorldImage sends the world into the IoOutput channel
func OutputWorldImage(c remoteutil.DistributorChannels, p remoteutil.Params, world [][]byte) {
	c.IoCommand <- remoteutil.IoOutput
	filename := strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight), strconv.Itoa(c.CompletedTurns)}, "x")
	c.IoFilename <- filename

	for m := 0; m < p.ImageHeight; m++ {
		for n := 0; n < p.ImageWidth; n++ {
			c.IoOutput <- world[m][n]
		}
	}
	c.Events <- remoteutil.ImageOutputComplete{c.CompletedTurns, filename}
}
