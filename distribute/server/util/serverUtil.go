package util

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
