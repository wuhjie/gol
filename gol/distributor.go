package gol

import (
	"strconv"
	"strings"

	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events    chan<- Event
	ioCommand chan<- ioCommand
	ioIdle    <-chan bool
	// adding filename into distributor channel
	filename chan<- string
}

//calculation
func mod(x, m int) int {
	return (x + m) % m
}

//used to calculate the alive neighbors
func calculateNeighbors(p Params, x, y int, world [][]byte) int {
	neighbors := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i != 0 || j != 0 {
				if world[mod(y+i, p.ImageHeight)][mod(x+j, p.ImageWidth)] == alive {
					neighbors++
				}
			}
		}
	}
	return neighbors
}

// calculate the world after changing
func calculateNextStage(p Params, world [][]byte) [][]byte {
	newWorld := make([][]byte, p.ImageHeight)
	for i := range newWorld {
		newWorld[i] = make([]byte, p.ImageWidth)
	}

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			neighbors := calculateNeighbors(p, x, y, world)
			if world[y][x] == alive {
				if neighbors == 2 || neighbors == 3 {
					newWorld[y][x] = alive
				} else {
					newWorld[y][x] = dead
				}
			}
			if world[y][x] == dead {
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

func calculateAliveCells(p Params, world [][]byte) []util.Cell {
	// util.ReadAliveCells()
	aliveCells := []util.Cell{}

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			if world[y][x] == alive {
				aliveCells = append(aliveCells, util.Cell{X: x, Y: y})
			}
		}
	}
	return aliveCells
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	// TODO: Create a 2D slice to store the world.
	//width length
	world := make([][]byte, p.ImageHeight)
	for i := range world {
		world[i] = make([]byte, p.ImageWidth)
	}

	tempWorld := make([][]byte, p.ImageHeight)
	for i := range tempWorld {
		tempWorld[i] = make([]byte, p.ImageWidth)
	}

	//for implementing the ioinput
	c.ioCommand <- ioInput

	var turnCount = 0
	//Execute all turns of the Game of Life.
	// confusing about if the next stage means we only calculate turns-1
	for turn := 0; turn < p.Turns; turn++ {
		// caltulate the changes in each iteration
		tempWorld = calculateNextStage(p, world)

		world = tempWorld
		turnCount = turn
	}

	//calculate the alive cells
	calculateAliveCells(p, tempWorld)

	// extract the defined filename in each iteration && pass the filename to the iochannel
	fileName := strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight)}, "x")
	c.filename <- fileName

	// TODO: Send correct Events when required, e.g. CellFlipped, TurnComplete and FinalTurnComplete.
	//		 See event.go for a list of all events.

	// Make sure that the Io has finished any output before exiting.

	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turnCount, Quitting}
	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
