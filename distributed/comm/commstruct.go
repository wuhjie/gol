package comm

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

// BrokerRequest is used to sending world at first
type BrokerRequest struct {
	World       [][]byte
	Threads     int
	ImageWidth  int
	ImageHeight int
}

// BrokerConnection creates the connections with address and function name
type BrokerConnection struct {
	SentInfo string
}

// BrokerReturn is returned struct from broker
type BrokerReturn struct {
	World        [][]byte
	ChangedCells []Cell
}

// ResponseOnReceivedWorld is the response when receicing the world
type ResponseOnReceivedWorld struct {
	Msg string
}

// BrokerSaved is the saving variables on broker
type BrokerSaved struct {
	World       [][]byte
	Threads     int
	ImageWidth  int
	ImageHeight int
	Turns       int
}

// QStatus gives the status of is q is pressed
type QStatus struct {
	Status bool
}

// KStatus is if K is sent on local machine
type KStatus struct {
	Status bool
}

// KQuitting is if K is received on the aws machine
type KQuitting struct {
	Msg string
}

type CommonMsg struct {
	Msg string
}

// WorkerRequest is what worker need in stage 1
type WorkerRequest struct {
	StartX int
	EndX   int
	StartY int
	EndY   int
	World  [][]byte
}

// WorkerReply is returned struct from each worker
type WorkerReply struct {
	PartWorld    [][]byte
	ChangedCells []Cell
}
