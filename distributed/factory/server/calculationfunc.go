package server

import (
	"uk.ac.bris.cs/gameoflife/client/gol"
	"uk.ac.bris.cs/gameoflife/comm"
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

//CalculateAliveCells calculates the alive cells in current round
func CalculateAliveCells(p gol.Params, world [][]byte) []comm.Cell {
	var aliveCells []comm.Cell

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			if world[y][x] == alive {
				aliveCells = append(aliveCells, comm.Cell{X: x, Y: y})
			}
		}
	}
	return aliveCells
}

//CalculateNeighbors is used to calculate the alive neighbors
func CalculateNeighbors(p gol.Params, x, y int, world func(y, x int) byte) int {
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

// CalculateNextStage implements basic calculation on aws
func CalculateNextStage(startY, endY, startX, endX int, p gol.Params, world func(y, x int) byte) (calculated Calculated) {

	var AliveCells []comm.Cell

	newWorld := make([][]byte, endY-startY)
	for i := range newWorld {
		newWorld[i] = make([]byte, p.ImageWidth)
	}

	// width and height of current piece
	height := endY - startY
	width := endX - startX

	// calculate world in current piece; the cell need to compare with the cell in the original world
	for y := 0; y < height; y++ {
		absoluteY := y + startY
		for x := 0; x < width; x++ {
			neighbors := CalculateNeighbors(p, x, absoluteY, world)
			switch world(absoluteY, x) {
			case alive:
				if neighbors == 2 || neighbors == 3 {
					newWorld[y][x] = alive
				} else {
					newWorld[y][x] = dead
					renewCell := comm.Cell{X: x, Y: absoluteY}
					AliveCells = append(AliveCells, renewCell)
				}
			case dead:
				if neighbors == 3 {
					newWorld[y][x] = alive
					renewCell := comm.Cell{X: x, Y: absoluteY}
					AliveCells = append(AliveCells, renewCell)

				} else {
					newWorld[y][x] = dead
				}
			}
		}
	}

	calculated.CalculatedWorld = newWorld
	calculated.AliveCellsCount = AliveCells

	return calculated
}
