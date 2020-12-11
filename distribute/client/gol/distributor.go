package gol

import (
	"fmt"
	"net/rpc"
	"strconv"
	"strings"
	"time"

	"uk.ac.bris.cs/gameoflife/client/util"
	"uk.ac.bris.cs/gameoflife/commstruct"
)

// Server is used for flag related
var Server *string

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

// InputWorldImage is related to loading images from the io channel
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
	qStatus := false
	// numbers of workers
	numofnode := 1

	// variables that need all the time
	world := InputWorldImage(p, c)
	initialsent := commstruct.BrokerRequest{
		World:       world,
		Threads:     2,
		ImageWidth:  p.ImageWidth,
		ImageHeight: p.ImageHeight,
		NumOfNode:   numofnode,
	}

	//
	msgIfWorldReceived := new(commstruct.ResponseOnReceivedWorld)
	client.Call("Broker.WorldReceived", initialsent, msgIfWorldReceived)

	for turns > 0 {
		localsent := commstruct.BrokerConnection{
			SentInfo: "local sent to broker",
		}
		BrokerReturn := new(commstruct.BrokerReturn)
		client.Call("Broker.Calculate", localsent, BrokerReturn)

		remoteAliveCells := BrokerReturn.ChangedCells
		for _, aCells := range remoteAliveCells {
			c.Events <- CellFlipped{
				CompletedTurns: c.CompletedTurns,
				Cell:           util.Cell{X: aCells.X, Y: aCells.Y},
			}
		}

		world = BrokerReturn.World

		turns--
		c.CompletedTurns = p.Turns - turns
		c.Events <- TurnComplete{
			CompletedTurns: c.CompletedTurns}

		// different conditions
		select {
		case <-ticker.C:
			c.Events <- AliveCellsCount{
				CompletedTurns: c.CompletedTurns,
				CellsCount:     len(CalculateAliveCells(p, world))}

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
				pStatus := 0

				for {
					command := <-c.KeyPresses
					switch command {
					case 'p':
						fmt.Println("Continuing")
						c.Events <- StateChange{c.CompletedTurns, Executing}
						c.Events <- TurnComplete{c.CompletedTurns}
						pStatus = 1
					}
					if pStatus == 1 {
						break
					}
				}
			case 'k':
				OutputWorldImage(c, p, world)
				kQuitMsg := new(commstruct.KQuitting)
				kstatus := commstruct.KStatus{
					Status: true,
				}
				client.Call("Broker.QuittingBroker", kstatus, kQuitMsg)
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

	c.Events <- FinalTurnComplete{c.CompletedTurns, CalculateAliveCells(p, world)}
	c.Events <- StateChange{c.CompletedTurns, Quitting}
	close(c.Events)
}
