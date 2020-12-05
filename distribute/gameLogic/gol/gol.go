package gol

import (
	"flag"
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

func gameLogicRunning(p Params, events chan<- Event, keyPresses <-chan rune, c DistributorChannels, io ioChannels) {

	go distributor(p, c)

	go startIo(p, io)

}

// establish rpc connection, as client/user
func userNetworkConnectionRelated(p Params, events chan<- Event, keyPresses <-chan rune, c DistributorChannels, io ioChannels) {
	server := flag.String("server", "127.0.0.1:8030", "IP: port string to connect to as server")
	flag.Parse()
	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()

	GameRunning := true

	//todo adding things to return when the game is supposed to end
	for GameRunning == true {
		// running logic as a user
		gameLogicRunning(p, events, keyPresses, c, io)
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

	userNetworkConnectionRelated(p, events, keyPresses, distributorChannels, ioChannels)

}
