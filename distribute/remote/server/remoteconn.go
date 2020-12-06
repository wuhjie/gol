package server

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

// Server is related to local machine
type Server struct {
	// gameStatus bool
	World [][]byte
}

// capability to work simultaneously
func worker(startY, endY, startX, endX int, p Params, immutableWorld func(y, x int) byte, c DistributorChannels, tempWorld chan<- [][]byte) {
	calculatedPart := CalculateNextStage(startY, endY, startX, endX, p, immutableWorld)
	tempWorld <- calculatedPart
}
