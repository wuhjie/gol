package gol

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"uk.ac.bris.cs/gameoflife/util"
)

type DistributorChannels struct {
	Events          chan<- Event //Events is what communicate with SDL
	IoCommand       chan<- ioCommand
	IoIdle          <-chan bool
	IoFilename      chan<- string
	AliveCellsCount chan<- []util.Cell
	IoInput         <-chan uint8
	IoOutput        chan<- uint8
	CompletedTurns  int
	KeyPresses      <-chan rune
}

// capability to work simultaneously
func worker(startY, endY, startX, endX int, p Params, immutableWorld func(y, x int) byte, c DistributorChannels, tempWorld chan<- [][]byte) {
	calculatedPart := calculateNextStage(startY, endY, startX, endX, p, immutableWorld, c)
	tempWorld <- calculatedPart
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c DistributorChannels) {

	world := initialisedWorld(p.ImageHeight, p.ImageWidth)

	ticker := time.NewTicker(2 * time.Second)

	c.IoCommand <- ioInput
	c.IoFilename <- strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight)}, "x")

	//adding the values in IoInput channel to initialised world inside distributor
	//flipped the initial alive cells
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			val := <-c.IoInput
			world[y][x] = val
			if val == alive {
				c.Events <- CellFlipped{CompletedTurns: 0, Cell: struct{ X, Y int }{X: x, Y: y}}
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
			go worker(i*heightPerThread, (i+1)*heightPerThread, 0, p.ImageWidth, p, immutableWorld, c, tempWorld[i])
		}
		go worker((p.Threads-1)*heightPerThread, p.ImageHeight, 0, p.ImageWidth, p, immutableWorld, c, tempWorld[p.Threads-1])

		// merge calculated world in each threads
		mergedWorld := initialisedWorld(0, 0)
		for i := 0; i < p.Threads; i++ {
			pieces := <-tempWorld[i]
			mergedWorld = append(mergedWorld, pieces...)
		}

		world = mergedWorld
		turns--

		c.Events <- TurnComplete{c.CompletedTurns}
		c.CompletedTurns = p.Turns - turns

		// sdl and ticker condition related
		select {
		case <-ticker.C:
			c.Events <- AliveCellsCount{c.CompletedTurns, len(calculateAliveCells(p, world))}
			// todo change how we control the game engine
		case command := <-c.KeyPresses:
			switch command {
			case 's':
				c.Events <- StateChange{c.CompletedTurns, Executing}
				OutputWorldImage(c, p, world)
			case 'q':
				c.Events <- StateChange{c.CompletedTurns, Quitting}
				qStatus = true
			case 'p':
				c.Events <- StateChange{c.CompletedTurns, Paused}
				OutputWorldImage(c, p, world)

				for {
					command := <-c.KeyPresses
					if command == 'p' {
						fmt.Println("Continuing")
						c.Events <- StateChange{c.CompletedTurns, Executing}
						c.Events <- TurnComplete{c.CompletedTurns}
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

	OutputWorldImage(c, p, world)

	c.IoCommand <- ioCheckIdle
	<-c.IoIdle
	c.Events <- FinalTurnComplete{c.CompletedTurns, calculateAliveCells(p, world)}
	c.Events <- StateChange{c.CompletedTurns, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.Events)
}
