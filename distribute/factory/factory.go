package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"

	"uk.ac.bris.cs/gameoflife/client/gol"
	"uk.ac.bris.cs/gameoflife/factory/server"

	"uk.ac.bris.cs/gameoflife/commstruct"
)

// Remote structure
type Remote struct{}

// CalculateNextTurn calculates the world after changing, called every turn
func (r *Remote) CalculateNextTurn(req commstruct.WorkerRequest, res *commstruct.WorkerReply) error {

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

// func (r *Remote) WorkerCalculation(req commstruct.WorkerRequest, res commstruct.WorkerReply) error {

// }

// QuitingServer is used to quit remote server
func (r *Remote) QuitingServer(kStatus commstruct.KStatus, kQuitting *commstruct.KQuitting) error {
	if kStatus.Status == true {
		os.Exit(0)
	}
	return nil
}

func main() {
	listener, err := net.Listen("tcp", ":8050")
	rpc.Register(&Remote{})
	rpc.Accept(listener)
	defer listener.Close()

	if err != nil {
		fmt.Println(err)
	}

}
