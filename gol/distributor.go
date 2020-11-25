package gol

import (
	"strconv"
	"strings"
	"time"
	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events    chan<- Event //events is what communicate with SDL
	ioCommand chan<- ioCommand
	ioIdle    <-chan bool
	ioFilename chan<- string
	aliveCellsCount chan<- []util.Cell
	ioInput   <-chan uint8
	ioOutput   chan<- uint8
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

//calculate the alive cells in current round
func calculateAliveCells(p Params, world [][]byte) []util.Cell {
	// util.ReadAliveCells()
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

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels, io ioChannels) {

	//width length
	world := make([][]byte, p.ImageHeight)
	for i := range world {
		world[i] = make([]byte, p.ImageWidth)
	}

	// initialised ticker for sending alive cells
	ticker := time.NewTicker(2 * time.Second)

	//for implementing the ioinput
	c.ioCommand <- ioInput
	c.ioFilename <- strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight)}, "x")

	//adding the vals in ioInput channel to initialised world inside distributor
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			val := <-c.ioInput
			world[y][x] = val
		}
	}

	//Execute all turns of the Game of Life.
	turns := p.Turns
	for turns > 0 {
		// calculate the changes in each iteration
		tempWorld := calculateNextStage(p, world)
		world = tempWorld
		//turnCount = turn
		turns--
		//ticker related
		select {
		case <-ticker.C:
			c.events <- AliveCellsCount{p.Turns-turns, len(calculateAliveCells(p, world))}
		default:
		}
	}

	//outputting the events
	c.ioCommand <- ioOutput
	c.ioFilename <- strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight)}, "x")

	for m := 0; m < p.ImageHeight; m++ {
		for n := 0; n < p.ImageWidth; n++ {
			c.ioOutput <- world[m][n]
		}
	}

	c.ioCommand <- ioCheckIdle
	<-c.ioIdle



	c.events <- FinalTurnComplete{p.Turns-turns, calculateAliveCells(p, world)}
	c.events <- StateChange{p.Turns-turns, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)

	// TODO: Send correct Events when required, e.g. CellFlipped, TurnComplete and FinalTurnComplete.
	//		 See event.go for a list of all events.

}
