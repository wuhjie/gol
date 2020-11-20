package gol

import (
	"strconv"
	"strings"
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

// Run starts the processing of Game of Life. It should initialise channels and goroutines.
func Run(p Params, events chan<- Event, keyPresses <-chan rune) {

	ioCommand := make(chan ioCommand)
	ioIdle := make(chan bool)

	distributorChannels := distributorChannels{
		events,
		ioCommand,
		ioIdle,
	}
	go distributor(p, distributorChannels)

	fileName := strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight)}, "x")

	// ioFilename := make(chan string)

	//inicialise input and output channel
	ioInput := make(chan uint8)
	ioOutput := make(chan uint8)

	ioChannels := ioChannels{
		command:  ioCommand,
		idle:     ioIdle,
		filename: nil,
		output:   ioOutput,
		input:    ioInput,
	}

	ioChannels.filename <- fileName

	go startIo(p, ioChannels)

}
