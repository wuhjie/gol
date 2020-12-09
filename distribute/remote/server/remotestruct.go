package server

import "uk.ac.bris.cs/gameoflife/commstruct"

// Params is with the same structure with the local machine
type Params struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}

// ioCommand related
type ioCommand uint8

const (
	IoOutput ioCommand = iota
	IoInput
	IoCheckIdle
)

// DistributorChannels is calculation related
type DistributorChannels struct {
	Events          chan<- RemoteEvent
	IoCommand       chan<- ioCommand
	IoIdle          <-chan bool
	IoFilename      chan<- string
	AliveCellsCount chan<- []commstruct.Cell
	IoInput         <-chan uint8
	IoOutput        chan<- uint8
	CompletedTurns  int
	KeyPresses      <-chan rune
}

// Calculated is returned from remote server
type Calculated struct {
	CalculatedWorld [][]byte
	AliveCellsCount []commstruct.Cell
}
