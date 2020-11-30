package gol

import (
	"fmt"
	"math"
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
func calculateNeighbors(p Params, x, y int,  world func(y, x int) byte) int {
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
	
	// width and height in current piece
	height := endY - startY
	width := endX - startX

	// making world with given width of pic
	for i := range newWorld {
		newWorld[i] = make([]byte, p.ImageWidth)
	}

	// calculate world in current piece
	for y := 0; y < height; y++ {
		// calculate the absolute coordinate
		absoluteY := y + startY

		for x := 0; x < width; x++ {
			neighbors := calculateNeighbors(p, x, absoluteY, world)
			if world(y, x) == alive {
				if neighbors == 2 || neighbors == 3 {
					newWorld[y][x] = alive
				} else {
					newWorld[y][x] = dead
					c.events <- CellFlipped{CompletedTurns: c.completedTurns, Cell: util.Cell{X: x, Y: absoluteY}}
				}
			}
			if world(y, x) == dead {
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

// capability to work simultaneously
func worker (startY, endY, startX, endX int, p Params, immutableWorld func(y, x int) byte, c distributorChannels, tempWorld chan<- [][]byte) {
	calculatedPart := calculateNextStage(startY, endY, startX, endX, p, immutableWorld, c)
	tempWorld <- calculatedPart
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

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	world := initialisedWorld(p.ImageHeight, p.ImageWidth)

	// initialised ticker for sending alive cells
	ticker := time.NewTicker(2 * time.Second)

	//for implementing the ioInput
	c.ioCommand <- ioInput
	c.ioFilename <- strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight)}, "x")

	//adding the values in ioInput channel to initialised world inside distributor
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
		immutableWorld := makeImmutableWorld(world)

		// tempWorld for each iteration
		tempWorld := make([]chan [][]byte, p.Threads)
		// adding channels
		for i := range tempWorld {
			tempWorld[i] = make(chan [][]byte)
		}

		// worker functions for different numbers of threads
		if p.ImageHeight % p.Threads == 0 {
			heightPerThread := p.ImageHeight / p.Threads
			for i := 0; i < p.Threads; i++ {
				fmt.Print("calculated height")
				fmt.Print(i*heightPerThread)
				fmt.Print("-")
				fmt.Println((i+1)*heightPerThread)
				go worker(i*heightPerThread, (i+1)*heightPerThread,0 , p.ImageWidth, p, immutableWorld, c, tempWorld[i])
			}
		}
		if p.ImageHeight % p.Threads != 0 {
			heightPerThread := int(math.Floor(float64(p.ImageHeight / p.Threads)))
			//fmt.Print("height-----")
			//fmt.Println(heightPerThread)

			for i := 0; i < p.Threads-1; i++ {
				fmt.Print("calculated height-----")
				fmt.Print(i*heightPerThread)
				fmt.Print("-")
				fmt.Println((i+1)*heightPerThread)
				go worker(i*heightPerThread, (i+1)*heightPerThread,0 , p.ImageWidth, p, immutableWorld, c, tempWorld[i])
			}
			fmt.Print("calculated height-----")
			fmt.Print((p.Threads-1) * heightPerThread)
			fmt.Print("-")
			fmt.Println(p.ImageHeight)
			//for the rest of pictures
			go worker((p.Threads-1) * heightPerThread, p.ImageHeight,0 , p.ImageWidth, p, immutableWorld, c, tempWorld[p.Threads-1])
		}

		//merging components together with initialised new empty world
		mergedWorld := initialisedWorld(0,0)

		// merge calculated world in each threads
		for i:= 0; i < p.Threads; i++ {
			pieces := <-tempWorld[i]
			mergedWorld = append(mergedWorld, pieces...)
		}
		fmt.Print("length of mergedWorld------")
		fmt.Println(len(mergedWorld))

		world = mergedWorld
		turns--
		// for output pic into the window
		c.events <- TurnComplete{c.completedTurns}
		c.completedTurns = p.Turns-turns

		// different conditions
		select {
		//ticker related
		case <-ticker.C:
			c.events <- AliveCellsCount{c.completedTurns, len(calculateAliveCells(p, world))}
		case command := <-c.keyPresses:

			// s---generate a PGM file with the current state of the board.
			if command == 's' {
				c.events <- StateChange {c.completedTurns, Executing}
				//for print out image
				outputWorldImage(c, p, world)
			}
			// q---generate a PGM file with the current state of the board and then terminate the program. Don't execute all turns set in gol.Params.Turns.
			if command == 'q' {
				c.events <- StateChange {c.completedTurns, Quitting}
				qStatus = true
			}
			// p---pause the processing and print the current turn that is being processed && resume the processing and print "Continuing" && q,s don't work
			if command == 'p' {
				c.events <- StateChange {c.completedTurns, Paused}

				//for print out image
				outputWorldImage(c, p, world)

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
	outputWorldImage(c, p, world)

	c.ioCommand <- ioCheckIdle
	<-c.ioIdle
	c.events <- FinalTurnComplete{c.completedTurns, calculateAliveCells(p, world)}
	c.events <- StateChange{c.completedTurns, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
