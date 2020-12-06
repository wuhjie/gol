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


// capability to work simultaneously
func worker (startY, endY, startX, endX int, p Params, immutableWorld func(y, x int) byte, c distributorChannels, tempWorld chan<- [][]byte) {
	calculatedPart := calculateNextStage(startY, endY, startX, endX, p, immutableWorld, c)
	tempWorld <- calculatedPart
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	world := initialisedWorld(p.ImageHeight, p.ImageWidth)

	ticker := time.NewTicker(2 * time.Second)

	c.ioCommand <- ioInput
	c.ioFilename <- strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight)}, "x")

	//adding the values in ioInput channel to initialised world inside distributor
	//flipped the initial alive cells
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			val := <-c.ioInput
			world[y][x] = val
			if val == alive {
				c.events <- CellFlipped{CompletedTurns: 0, Cell: struct{ X, Y int }{X:x, Y:y}}
			}
		}
	}

	turns := p.Turns
	qStatus := false

	for turns > 0 {
		immutableWorld := makeImmutableWorld(world)

		tempWorld := make([]chan [][]byte, p.Threads)
		for i := range tempWorld {
			tempWorld[i] = make(chan [][]byte)
		}

		heightPerThread := p.ImageHeight / p.Threads
		for i := 0; i < p.Threads-1; i++ {
			go worker(i*heightPerThread, (i+1)*heightPerThread,0 , p.ImageWidth, p, immutableWorld, c, tempWorld[i])
		}
		go worker((p.Threads-1)*heightPerThread, p.ImageHeight,0 , p.ImageWidth, p, immutableWorld, c, tempWorld[p.Threads-1])

		// merge calculated world in each threads
		mergedWorld := initialisedWorld(0,0)
		for i:= 0; i < p.Threads; i++ {
			pieces := <-tempWorld[i]
			mergedWorld = append(mergedWorld, pieces...)
		}

		world = mergedWorld
		turns--

		c.events <- TurnComplete{c.completedTurns}
		c.completedTurns = p.Turns-turns

		// different conditions
		select {
		case <-ticker.C:
			c.events <- AliveCellsCount{c.completedTurns, len(calculateAliveCells(p, world))}
		case command := <-c.keyPresses:
			switch command{
			case 's':
				c.events <- StateChange {c.completedTurns, Executing}
				outputWorldImage(c, p, world)
			case 'q':
				c.events <- StateChange {c.completedTurns, Quitting}
				qStatus = true
			case 'p':
				c.events <- StateChange {c.completedTurns, Paused}
				outputWorldImage(c, p, world)
				pStatus := 0

				for {
					command := <-c.keyPresses
					switch command{
					case 'p':
						fmt.Println("Continuing")
						c.events <- StateChange {c.completedTurns, Executing}
						c.events <- TurnComplete{c.completedTurns}
						pStatus = 1
					}
					if pStatus == 1 {
						break
					}
				}
			}
		default:
		}
		// for quiting the programme: q
		if qStatus == true {
			break
		}
	}

	outputWorldImage(c, p, world)

	c.ioCommand <- ioCheckIdle
	<-c.ioIdle
	c.events <- FinalTurnComplete{c.completedTurns, calculateAliveCells(p, world)}
	c.events <- StateChange{c.completedTurns, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
