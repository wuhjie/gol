package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"time"

	"uk.ac.bris.cs/gameoflife/server"
)

func handleConnection(conn *net.Conn) {
	fmt.Println("connection established")

}

//Server
type Server struct{}

// CalculationRunning implements basic calculation on aws
func (server *Server) CalculationRunning(p server.Params) {

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

		// different conditions
		select {
		case <-ticker.C:
			c.events <- AliveCellsCount{c.completedTurns, len(calculateAliveCells(p, world))}
		case command := <-c.keyPresses:
			switch command {
			case 's':
				c.events <- StateChange{c.completedTurns, Executing}
				outputWorldImage(c, p, world)
			case 'q':
				c.events <- StateChange{c.completedTurns, Quitting}
				qStatus = true
			case 'p':
				c.events <- StateChange{c.completedTurns, Paused}
				outputWorldImage(c, p, world)
				pStatus := 0

				for {
					command := <-c.keyPresses
					switch command {
					case 'p':
						fmt.Println("Continuing")
						c.events <- StateChange{c.completedTurns, Executing}
						c.events <- TurnComplete{c.completedTurns}
						pStatus = 1
					}
					if pStatus == 1 {
						break
					}
				}
			}
		default:
		}
		// for quiting the programme: q
		if qStatus == true {
			break
		}
	}

}

func main() {
	pAddr := flag.String("port", "8030", "port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	rpc.Register(&Server{})
	listener, _ := net.Listen("tcp", ":"+*pAddr)
	defer listener.Close()
	rpc.Accept(listener)
}
