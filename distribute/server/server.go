package server

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"time"

	"uk.ac.bris.cs/gameoflife/util"
)

type DistributorChannels struct {
	Events chan<- Event //Events is what communicate with SDL
	// IoCommand       chan<- ioCommand
	IoIdle          <-chan bool
	IoFilename      chan<- string
	AliveCellsCount chan<- []util.Cell
	IoInput         <-chan uint8
	IoOutput        chan<- uint8
	CompletedTurns  int
	KeyPresses      <-chan rune
}

// capability to work simultaneously
func worker(startY, endY, startX, endX int, p util.Params, immutableWorld func(y, x int) byte, c DistributorChannels, tempWorld chan<- [][]byte) {
	calculatedPart := CalculateNextStage(startY, endY, startX, endX, p, immutableWorld, c)
	tempWorld <- calculatedPart
}

type Server struct {
	// gameStatus bool
	World [][]byte
}

// CalculationRunning implements basic calculation on aws
func (server *Server) CalculationRunning(p util.Params, world [][]byte, c DistributorChannels) {

	turns := p.Turns

	for turns > 0 {
		immutableWorld := MakeImmutableWorld(world)

		tempWorld := make([]chan [][]byte, p.Threads)
		for i := range tempWorld {
			tempWorld[i] = make(chan [][]byte)
		}

		heightPerThread := p.ImageHeight / p.Threads
		for i := 0; i < p.Threads-1; i++ {
			go worker(i*heightPerThread, (i+1)*heightPerThread, 0, p.ImageWidth, p, immutableWorld, c, tempWorld[i])
		}
		go worker((p.Threads-1)*heightPerThread, p.ImageHeight, 0, p.ImageWidth, p, immutableWorld, c, tempWorld[p.Threads-1])

		// merge calculated world in each threads
		mergedWorld := InitialisedWorld(0, 0)
		for i := 0; i < p.Threads; i++ {
			pieces := <-tempWorld[i]
			mergedWorld = append(mergedWorld, pieces...)
		}

		world = mergedWorld
		turns--

		c.Events <- TurnComplete{c.CompletedTurns}
		c.CompletedTurns = p.Turns - turns
	}

}

func handleConnection(conn *net.Conn) {
	for {
		fmt.Println("connection established")
	}
}

func main() {
	pAddr := flag.String("port", "8030", "port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	listener, _ := net.Listen("tcp", ":"+*pAddr)

	defer listener.Close()
	rpc.Accept(listener)
}
