package remoteutil

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
	event           RemoteEvent
	ioCommand       ioCommand
	P               Params
	World           [][]byte
}

// RemoteReply is what the local machine need
type RemoteReply struct {
	aliveCellsCount int
	completedTurns  int
}
