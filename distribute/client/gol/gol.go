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

// User contains user related part
type User struct {
	InitialWorld    chan [][]byte
	aliveCellsCount int
}

// RemoteReply contains things that needed from the remote server
type RemoteReply struct {
	aliveCellsCount int
	completedTurns  int
	event           util.Event
	ioCommand       ioCommand
}

// establish rpc connection, as client/user
func userNetworkConnectionRelated(p Params, c DistributorChannels, io ioChannels) {

	serverAdd := "localhost"
	client, err := rpc.DialHTTP("tcp", serverAdd+":8080")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// initialising user struct
	initialWorld := make(chan [][]byte)
	acellsCount := 0

	u := User{
		InitialWorld:    initialWorld,
		aliveCellsCount: acellsCount,
	}

	// goroutine related
	go LocalFilesReading(u, p, c)
	go startIo(p, io)

	select {
	// sending readed world to remote server
	case localWorld := <-u.InitialWorld:
		fmt.Println("initialWorld received")

		worldToRemote := &

		// todo replay
		var reply RemoteReply
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
