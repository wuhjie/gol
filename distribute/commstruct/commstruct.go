package commstruct

// Cell is the same as local and remote cell
type Cell struct {
	X, Y int
}

// Localsent contains things that needed from the remote server
type Localsent struct {
	Turns int
}

// RemoteReply is what the local machine need
type RemoteReply struct {
	aliveCellsCount int
	completedTurns  int
	AliveCells      []Cell
	World           [][]byte
}

// InitialToRemote is used to sending world at first
type InitialToRemote struct {
	World       [][]byte
	Threads     int
	ImageWidth  int
	ImageHeight int
}

// ResponseOnReceivedWorld is the response when receicing the world
type ResponseOnReceivedWorld struct {
	Msg string
}

// KStatus is if K is sent on local machine
type KStatus struct {
	Status bool
}

// KQuiting is if K is received on the aws machine
type KQuitting struct {
	Msg string
}
