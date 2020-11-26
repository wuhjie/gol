package gol

import (
	"fmt"
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
	completedTurns int
	keyPresses <-chan rune
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
func calculateNextStage(p Params, world [][]byte, c distributorChannels) [][]byte {
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
					c.events <- CellFlipped{CompletedTurns: c.completedTurns, Cell: util.Cell{X: y, Y: x}}
				}
			}
			if world[y][x] == dead {
				if neighbors == 3 {
					newWorld[y][x] = alive
					c.events <- CellFlipped{CompletedTurns: c.completedTurns, Cell: util.Cell{X: y, Y: x}}
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
func distributor(p Params, c distributorChannels) {

	//width length
	world := make([][]byte, p.ImageHeight)
	for i := range world {
		world[i] = make([]byte, p.ImageWidth)
	}

	// initialised ticker for sending alive cells
	ticker := time.NewTicker(2 * time.Second)

	//for implementing the ioInput
	c.ioCommand <- ioInput
	c.ioFilename <- strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight)}, "x")

	//adding the vals in ioInput channel to initialised world inside distributor
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			val := <-c.ioInput
			world[y][x] = val
			//flipped the initial alive cells
			if val == alive {
				c.events <- CellFlipped{CompletedTurns: 0, Cell: struct{ X, Y int }{X:x, Y:y}}
			}
		}
	}

	//Execute all turns of the Game of Life.
	turns := p.Turns
	qStatus := false

	for turns > 0 {
		// calculate the changes in each iteration
		tempWorld := calculateNextStage(p, world, c)
		world = tempWorld
		//turnCount = turn
		turns--
		// for output pic into the window
		c.events <- TurnComplete{c.completedTurns}
		c.completedTurns = p.Turns-turns

		select {
		//ticker related
		case <-ticker.C:
			c.events <- AliveCellsCount{c.completedTurns, len(calculateAliveCells(p, world))}
		case command := <-c.keyPresses:
			// If s is pressed, generate a PGM file with the current state of the board.

			if command == 's' {
				c.events <- StateChange {c.completedTurns, Executing}

				//for print out image
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
			// If q is pressed, generate a PGM file with the current state of the board and then terminate the program.
			// Your program should not continue to execute all turns set in gol.Params.Turns.
			if command == 'q' {
				c.events <- StateChange {c.completedTurns, Quitting}
				qStatus = true
			}
			// If p is pressed, pause the processing and print the current turn that is being processed.
			// If p is pressed again resume the processing and print "Continuing".
			// It is not necessary for q and s to work while the execution is paused.
			if command == 'p' {
				c.events <- StateChange {c.completedTurns, Paused}

				//for print out image
				c.ioCommand <- ioOutput
				filename := strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight), strconv.Itoa(c.completedTurns)}, "x")
				c.ioFilename <- filename

				for m := 0; m < p.ImageHeight; m++ {
					for n := 0; n < p.ImageWidth; n++ {
						c.ioOutput <- world[m][n]
					}
				}
				c.events <- ImageOutputComplete{c.completedTurns, filename}

				// waiting for the next p
				for {
					command := <-c.keyPresses
					if command == 'p' {
						fmt.Println("Continuing")
						c.events <- StateChange {c.completedTurns, Executing}
						c.events <- TurnComplete{c.completedTurns}
					}
					break
				}
			}
		default:
		}
		// for quiting the programme: q
		if qStatus == true {
			break
		}
	}

	//outputting the events
	c.ioCommand <- ioOutput
	filename := strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight), strconv.Itoa(c.completedTurns)}, "x")
	c.ioFilename <- filename

	for m := 0; m < p.ImageHeight; m++ {
		for n := 0; n < p.ImageWidth; n++ {
			c.ioOutput <- world[m][n]
		}
	}

	c.events <- ImageOutputComplete{c.completedTurns, filename}

	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- FinalTurnComplete{c.completedTurns, calculateAliveCells(p, world)}
	c.events <- StateChange{c.completedTurns, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
