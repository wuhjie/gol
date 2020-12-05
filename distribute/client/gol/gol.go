package gol

import (
	"flag"
	"fmt"
	"net/rpc"

	"uk.ac.bris.cs/gameoflife/util"
)

const alive = 255
const dead = 0

// Params provides the details of how to run the Game of Life and which image to load.
type Params struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}

// UserChannels contains user related part
type UserChannels struct {
	initialWorld chan [][]byte
}

func gameLogicRunning(p Params, u UserChannels, io ioChannels, c DistributorChannels) {
	// essential goroutine running
	go Distributor(u, p, c)
	go startIo(p, io)

	select {
	// sending readed world to remote server
	case world := <-u.initialWorld:
		fmt.Println("initialWorld received")

	default:
	}
}

// establish rpc connection, as client/user
func userNetworkConnectionRelated(p Params, c DistributorChannels, io ioChannels) {
	server := flag.String("server", "127.0.0.1:8030", "IP: port string to connect to as server")
	flag.Parse()
	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()

	initialWorld := make(chan [][]byte)

	userChannels := UserChannels{
		initialWorld: initialWorld,
	}

	gameStatus := true

	//todo adding things to return when the game is supposed to end
	for gameStatus == true {
		// running logic as a user
		gameLogicRunning(p, userChannels, io, c)
	}

}

// Run starts the processing of Game of Life. It should initialise channels and goroutines.
func Run(p Params, events chan<- Event, keyPresses <-chan rune) {

	ioCommand := make(chan ioCommand)
	ioIdle := make(chan bool)

	// input and output channel
	ioInput := make(chan uint8)
	ioOutput := make(chan uint8)

	ioFilename := make(chan string)
	aliveCellsCount := make(chan []util.Cell)

	completedTurns := 0

	distributorChannels := DistributorChannels{
		events,
		ioCommand,
		ioIdle,
		ioFilename,
		aliveCellsCount,
		ioInput,
		ioOutput,
		completedTurns,
		keyPresses,
	}

	ioChannels := ioChannels{
		ioCommand,
		ioIdle,
		ioFilename,
		ioInput,
		ioOutput,
	}

	userNetworkConnectionRelated(p, distributorChannels, ioChannels)

}
