package server

const alive = 255
const dead = 0

// Params is with the same structure with the local machine
type Params struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}

// RemoteReply contains things that needed from the remote server
type RemoteReply struct {
	aliveCellsCount int
	completedTurns  int
	// event
}

// ioCommand related
type ioCommand uint8

const (
	ioOutput ioCommand = iota
	ioInput
	ioCheckIdle
)
