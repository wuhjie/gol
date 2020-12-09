package util

// InitialisedWorld makes new 2-Dimension array of byte
func InitialisedWorld(height, width int) [][]byte {
	world := make([][]byte, height)
	for i := range world {
		world[i] = make([]byte, width)
	}
	return world
}

// RemoteCalculation is calling methods, only the method is needed
var RemoteCalculation = "Remote.CalculateNextTurn"

// Localsent contains things that needed from the remote server
type Localsent struct {
	Turns       int
	World       [][]byte
	Threads     int
	ImageWidth  int
	ImageHeight int
}

// RemoteReply is what the local machine need
type RemoteReply struct {
	aliveCellsCount int
	completedTurns  int
	AliveCells      []Cell
	World           [][]byte
}
