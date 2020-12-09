package main

import (
	"net"
	"net/rpc"

	"uk.ac.bris.cs/gameoflife/gol"
	"uk.ac.bris.cs/gameoflife/server/server"
	"uk.ac.bris.cs/gameoflife/util"
)

// capability to work simultaneously
// func worker(startY, endY, startX, endX int, p gol.Params, immutableWorld func(y, x int) byte, tempWorld chan<- [][]byte) {
// 	calculatedPart := server.CalculateNextStage(startY, endY, startX, endX, p, immutableWorld)
// 	tempWorld <- calculatedPart
// }

// Remote structure
type Remote struct{}

// CalculateNextTurn calculates the world after changing, called every turn
func (r *Remote) CalculateNextTurn(localSent util.Localsent, remoteResponse *util.RemoteReply) error {

	p := gol.Params{
		localSent.Turns,
		localSent.Threads,
		localSent.ImageHeight,
		localSent.ImageWidth,
	}
	immutableWorld := server.MakeImmutableWorld(localSent.World)

	calculatedWorld := server.CalculateNextStage(0, p.ImageHeight, 0, p.ImageWidth, p, immutableWorld)

	remoteResponse.AliveCells = server.AliveCells
	remoteResponse.World = calculatedWorld

	return nil
}

func main() {
	listener, _ := net.Listen("tcp", ":8030")
	rpc.Register(&Remote{})
	rpc.Accept(listener)
}
