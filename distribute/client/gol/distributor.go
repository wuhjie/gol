package gol

import (
	"fmt"
	"net/rpc"
	"time"

	"uk.ac.bris.cs/gameoflife/client/util"
	"uk.ac.bris.cs/gameoflife/comm"
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

// Distributor imports read pgm file
func Distributor(p Params, c DistributorChannels) {

	// establish rpc connection
	client, _ := rpc.Dial("tcp", "127.0.0.1:8030")
	defer client.Close()

	ticker := time.NewTicker(2 * time.Second)

	brokerQStatus := new(comm.QStatus)
	req := comm.CommonMsg{Msg: "getting q status"}
	client.Call("Broker.GetQStatus", req, brokerQStatus)
	initialsent := comm.BrokerRequest{}
	turns := 0
	world := initialisedWorld(0, 0)
	qStatus := false

	switch brokerQStatus.Status {
	case true:
		savedStatus := new(comm.BrokerSaved)
		req := comm.CommonMsg{Msg: "getting broker status"}
		client.Call("Broker.GetBrokerStatus", req, savedStatus)
		turns = savedStatus.Turns
		initialsent = comm.BrokerRequest{
			World:       savedStatus.World,
			Threads:     savedStatus.Threads,
			ImageWidth:  savedStatus.ImageWidth,
			ImageHeight: savedStatus.ImageHeight,
		}
	case false:
		turns = p.Turns
		// variables that need all the time
		world = InputWorldImage(p, c)
		initialsent = comm.BrokerRequest{
			World:       world,
			Threads:     1,
			ImageWidth:  p.ImageWidth,
			ImageHeight: p.ImageHeight,
		}
		msgIfWorldReceived := new(comm.ResponseOnReceivedWorld)
		client.Call("Broker.WorldReceived", initialsent, msgIfWorldReceived)
	}

	for turns > 0 {
		localsent := comm.Localsent{
			Turns: turns,
		}
		BrokerReturn := new(comm.BrokerReturn)
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
				brokerreply := new(comm.CommonMsg)
				sentStr := comm.CommonMsg{Msg: "request about change q status"}
				client.Call("Broker.ModifyQStatus", sentStr, brokerreply)
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
				kQuitMsg := new(comm.KQuitting)
				kstatus := comm.KStatus{
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
