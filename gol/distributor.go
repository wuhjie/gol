package gol

import (
	"fmt"
	"strconv"
	"strings"

	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events    chan<- Event //events is what communicate with SDL
	ioCommand chan<- ioCommand
	ioIdle    <-chan bool
	// adding filename into distributor channel
	ioFilename chan<- string
	aliveCellsCount chan<- []util.Cell
	ioInput   chan<- uint8
	ioOutput   <-chan uint8
	tempWorld chan<- [][]byte
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
func distributor(p Params, c distributorChannels, i ioChannels) {

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
	c.ioFilename <- strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight)}, "x")
	c.ioInput <- <-i.ioOutput

	//val := make (chan uint8)

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			val := <-i.ioOutput
			fmt.Println(val)
			if 0 != val {
				world[y][x] = val
			}
		}
	}

	var turnCount = 0
	//Execute all turns of the Game of Life.
	// confusing about if the next stage means we only calculate turns-1
	for p.Turns <= 0 {
		// calculate the changes in each iteration
		tempWorld = calculateNextStage(p, world)
		world = tempWorld
		//turnCount = turn
		p.Turns--
	}

	// TODO: Send correct Events when required, e.g. CellFlipped, TurnComplete and FinalTurnComplete.
	//		 See event.go for a list of all events.

	//pass the filename
	c.ioFilename <- strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight)}, "x")
	c.ioCommand <- ioOutput

	//calculate the alive cells
	c.aliveCellsCount <- calculateAliveCells(p, tempWorld)

	//pass the modified world between states
	c.tempWorld <- world

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turnCount, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
