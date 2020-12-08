package gol

import (
	"uk.ac.bris.cs/gameoflife/util"
)

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

	go Distributor(p, distributorChannels)

	go startIo(p, ioChannels)

}
