package server

import (
	"flag"
	"fmt"
	gol2 "gol/distribute/client/gol"
	"math/rand"
	"net"
	"net/rpc"
	"time"
)

// capability to work simultaneously
func worker(startY, endY, startX, endX int, p gol2.Params, immutableWorld func(y, x int) byte, c gol2.DistributorChannels, tempWorld chan<- [][]byte) {
	calculatedPart := gol2.CalculateNextStage(startY, endY, startX, endX, p, immutableWorld, c)
	tempWorld <- calculatedPart
}

type Server struct {
	// gameStatus bool
	World [][]byte
}

// implementing basic calculation on aws
func (server *Server) CalculationRunning(p gol2.Params, world [][]byte, c gol2.DistributorChannels) {

	turns := p.Turns

	for turns > 0 {
		immutableWorld := gol2.MakeImmutableWorld(world)

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
		mergedWorld := gol2.InitialisedWorld(0, 0)
		for i := 0; i < p.Threads; i++ {
			pieces := <-tempWorld[i]
			mergedWorld = append(mergedWorld, pieces...)
		}

		world = mergedWorld
		turns--

		c.Events <- gol2.TurnComplete{c.CompletedTurns}
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
