package gol

import (
	"net/rpc"
	"strconv"
	"strings"
	"time"

	"uk.ac.bris.cs/gameoflife/util"
)

// DistributorChannels contains things that need for parallel calculation
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

// InputWorldImage
func InputWorldImage(p Params, c DistributorChannels) [][]byte {

	world := util.InitialisedWorld(p.ImageHeight, p.ImageWidth)
	c.IoCommand <- ioInput
	c.IoFilename <- strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight)}, "x")

	//adding the values in ioInput channel to initialised world inside distributor
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
	return world
}

// OutputWorldImage sends the world into the IoOutput channel
func OutputWorldImage(c DistributorChannels, p Params, world [][]byte) {
	c.IoCommand <- ioOutput
	filename := strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight), strconv.Itoa(c.CompletedTurns)}, "x")
	c.IoFilename <- filename

	for m := 0; m < p.ImageHeight; m++ {
		for n := 0; n < p.ImageWidth; n++ {
			c.IoOutput <- world[m][n]
		}
	}
	c.Events <- ImageOutputComplete{c.CompletedTurns, filename}
}

// CalculateAliveCells the alive cells in current round
func CalculateAliveCells(p Params, world [][]byte) []util.Cell {
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

// Distributor imports read pgm file
func Distributor(p Params, c DistributorChannels) {

	// establish rpc connection
	client, _ := rpc.Dial("tcp", "127.0.0.1:8030")
	defer client.Close()

	ticker := time.NewTicker(2 * time.Second)
	turns := p.Turns

	world := InputWorldImage(p, c)

	for turns > 0 {
		localsent := Localsent{
			CompletedTurns: p.Turns,
			World:          world,
			Threads:        p.Threads,
			ImageWidth:     p.ImageWidth,
			ImageHeight:    p.ImageHeight,
		}
		remotereply := new(RemoteReply)

		client.Call("Remote.CalculateNextTurn", localsent, remotereply)

		remoteAliveCells := remotereply.AliveCells
		for _, aCells := range remoteAliveCells {
			c.Events <- CellFlipped{
				CompletedTurns: c.CompletedTurns,
				Cell: util.Cell{
					X: aCells.X,
					Y: aCells.Y,
				},
			}
		}

		world := remotereply.World

		turns--
		c.CompletedTurns = p.Turns - turns
		c.Events <- TurnComplete{c.CompletedTurns}

		// different conditions
		select {
		case <-ticker.C:
			c.Events <- AliveCellsCount{c.CompletedTurns, len(CalculateAliveCells(p, world))}

		default:
		}
		// for quiting the programme: q
	}

	OutputWorldImage(c, p, world)

	c.IoCommand <- ioCheckIdle
	<-c.IoIdle

	c.Events <- FinalTurnComplete{c.CompletedTurns, CalculateAliveCells(p, world)}
	c.Events <- StateChange{c.CompletedTurns, Quitting}
	close(c.Events)
}
