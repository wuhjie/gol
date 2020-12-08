package main

import (
	"net"
	"net/rpc"

	"uk.ac.bris.cs/gameoflife/remoteutil"
	"uk.ac.bris.cs/gameoflife/server"
)

// capability to work simultaneously
func worker(startY, endY, startX, endX int, p remoteutil.Params, immutableWorld func(y, x int) byte, tempWorld chan<- [][]byte) {
	calculatedPart := server.CalculateNextStage(startY, endY, startX, endX, p, immutableWorld)
	tempWorld <- calculatedPart
}

// Remote structure
type Remote struct{}

// CalculateNextTurn calculates the world after changing, called every turn
func (r *Remote) CalculateNextTurn(localSent remoteutil.LocalSent, remoteResponse *remoteutil.RemoteReply) error {

	p := remoteutil.Params{
		localSent.Turns,
		localSent.Threads,
		localSent.ImageHeight,
		localSent.ImageWidth,
	}
	immutableWorld := server.MakeImmutableWorld(localSent.World)

	tempWorld := make([]chan [][]byte, p.Threads)
	for i := range tempWorld {
		tempWorld[i] = make(chan [][]byte)
	}

	heightPerThread := p.ImageHeight / p.Threads
	for i := 0; i < p.Threads-1; i++ {
		go worker(i*heightPerThread, (i+1)*heightPerThread, 0, p.ImageWidth, p, immutableWorld, tempWorld[i])
	}
	go worker((p.Threads-1)*heightPerThread, p.ImageHeight, 0, p.ImageWidth, p, immutableWorld, tempWorld[p.Threads-1])

	mergedWorld := server.InitialisedWorld(0, 0)
	for i := 0; i < p.Threads; i++ {
		pieces := <-tempWorld[i]
		mergedWorld = append(mergedWorld, pieces...)
	}

	// global variable
	remoteResponse.AliveCells = server.AliveCells
	remoteResponse.World = mergedWorld

	return nil
}

func main() {
	listener, _ := net.Listen("tcp", ":"+"8030")
	rpc.Register(&Remote{})
	rpc.Accept(listener)
}
