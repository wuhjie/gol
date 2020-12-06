package gol

import (
	"fmt"

	"log"
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
	InitialWorld    chan [][]byte
	aliveCellsCount chan int
}

// establish rpc connection, as client/user
func userNetworkConnectionRelated(p Params, c DistributorChannels, io ioChannels) {

	serverAdd := "localhost"
	client, err := rpc.DialHTTP("tcp", serverAdd+":8080")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	initialWorld := make(chan [][]byte)

	u := UserChannels{
		InitialWorld: initialWorld,
	}

	// goroutine related
	go LocalFilesReading(u, p, c)
	go startIo(p, io)

	select {
	// sending readed world to remote server
	case localWorld := <-u.InitialWorld:
		fmt.Println("initialWorld received")
		world := &server.Server{localWorld}
		var reply int
		err = client.Call("Server.CalculationRunning", world, &reply)
		if err != nil {
			log.Fatalf("world error:", err)
		}

	default:
	}
}

// Run starts game of life of the user side
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
