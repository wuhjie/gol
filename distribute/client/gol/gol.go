package gol

import (
	"flag"
	"fmt"

	"net/rpc"

	"uk.ac.bris.cs/gameoflife/util"
)

// CallRemoteCalculatiois related to  all remote calcylation
func CallRemoteCalculation(client rpc.Client, localWorld [][]byte) {
	request := Localsent{localWorld}
	response := new(RemoteReply)
	client.Call(util.RemoteCalculation, request, response)
	fmt.Println("calculation called")
}

// establish rpc connection, as client/user
func userNetworkConnectionRelated(p Params, c DistributorChannels, io ioChannels) {

	server := flag.String("server", "127.0.0.1:8030", "IP: port string to connect to as server")
	flag.Parse()
	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()

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
		fmt.Println("world received from local file")

		// connect to remote server
		CallRemoteCalculation(*client, localWorld)

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
