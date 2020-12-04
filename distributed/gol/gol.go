package gol

import "uk.ac.bris.cs/gameoflife/util"

const alive = 255
const dead = 0

// Params provides the details of how to run the Game of Life and which image to load.
type Params struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
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

	go distributor(p, distributorChannels)

	go startIo(p, ioChannels)
}


