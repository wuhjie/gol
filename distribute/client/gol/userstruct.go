package gol

const alive = 255
const dead = 0

// Params provides the details of how to run the Game of Life and which image to load.
type Params struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}

// User contains user related part
type User struct {
	InitialWorld    chan [][]byte
	aliveCellsCount int
}

// Localsent contains things that needed from the remote server
type Localsent struct {
	localWorld [][]byte
}

// RemoteReply is what the local machine need
type RemoteReply struct {
	aliveCellsCount int
	completedTurns  int
	// event
}
