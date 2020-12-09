package commstruct

// Cell is the same as local and remote cell
type Cell struct {
	X, Y int
}

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
