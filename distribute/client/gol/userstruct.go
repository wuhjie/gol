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

// Localsent contains things that needed from the remote server
type Localsent struct {
	CompletedTurns int
	World          [][]byte
	Threads        int
	ImageWidth     int
	ImageHeight    int
}

// RemoteReply is what the local machine need
type RemoteReply struct {
	aliveCellsCount int
	completedTurns  int
	AliveCells      []util.Cell
	World           [][]byte
}
