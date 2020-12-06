package server

const alive = 255
const dead = 0

// Cell is used as the return type for the testing framework.
type Cell struct {
	X, Y int
}

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
	ioOutput ioCommand = iota
	ioInput
	ioCheckIdle
)

// DistributorChannels is calculation related
type DistributorChannels struct {
	Events          chan<- Event
	IoCommand       chan<- ioCommand
	IoIdle          <-chan bool
	IoFilename      chan<- string
	AliveCellsCount chan<- []Cell
	IoInput         <-chan uint8
	IoOutput        chan<- uint8
	CompletedTurns  int
	KeyPresses      <-chan rune
}

// Localsent struct which is the same as localmachine
type Localsent struct {
	aliveCellsCount int
	completedTurns  int
	event           Event
	ioCommand       ioCommand
	P               Params
	World           [][]byte
}

// RemoteReply is what the local machine need
type RemoteReply struct {
	aliveCellsCount int
	completedTurns  int
	// event
}
