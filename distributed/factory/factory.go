package main

import (
	"net"
	"net/rpc"
	"os"

	"uk.ac.bris.cs/gameoflife/client/gol"
	"uk.ac.bris.cs/gameoflife/comm"
	"uk.ac.bris.cs/gameoflife/factory/server"
)

// Remote structure used as factory
type Remote struct{}

// CalculateNextTurn calculates the world after changing, called every turn
func (r *Remote) CalculateNextTurn(req comm.WorkerRequest, res *comm.WorkerReply) error {
	immutableWorld := server.MakeImmutableWorld(req.World)

	p := gol.Params{
		Turns:       0,
		Threads:     0,
		ImageHeight: req.EndY,
		ImageWidth:  req.EndX,
	}
	calculatedItems := server.CalculateNextStage(0, req.EndX, req.StartY, req.EndY, p, immutableWorld)

	res.PartWorld = calculatedItems.CalculatedWorld
	res.ChangedCells = calculatedItems.AliveCellsCount

	return nil
}

// QuitingFactory is used to quit factory
func (r *Remote) QuitingFactory(req comm.KStatus, res *comm.KQuitting) error {
	if req.Status == true {
		os.Exit(0)
	}
	return nil
}

func main() {
	listenerone, _ := net.Listen("tcp", ":8050")
	// listenertwo, _ := net.Listen("tcp", ":8060")
	rpc.Register(&Remote{})
	rpc.Accept(listenerone)
	// rpc.Accept(listenertwo)
	defer listenerone.Close()
	// defer listenertwo.Close()

}
