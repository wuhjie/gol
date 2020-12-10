package main

import (
	"flag"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"time"

	"uk.ac.bris.cs/gameoflife/client/gol"
	"uk.ac.bris.cs/gameoflife/commstruct"
	"uk.ac.bris.cs/gameoflife/remote/server"
)

// Below are global variables stored in remote server
var RemoteWorld [][]byte
var FixedThreads int
var PicHeight int
var PicWidth int

// Remote structure
type Remote struct{}

// WorldReceived is used to receive global variables
func (r *Remote) WorldReceived(initial commstruct.InitialToRemote, rep *commstruct.ResponseOnReceivedWorld) error {

	RemoteWorld = initial.World
	FixedThreads = initial.Threads
	PicHeight = initial.ImageHeight
	PicWidth = initial.ImageWidth

	rep.Msg = "world received from local machine"

	return nil
}

// CalculateNextTurn calculates the world after changing, called every turn
func (r *Remote) CalculateNextTurn(localSent commstruct.Localsent, remoteResponse *commstruct.RemoteReply) error {

	p := gol.Params{
		Turns:       localSent.Turns,
		Threads:     FixedThreads,
		ImageHeight: PicHeight,
		ImageWidth:  PicWidth,
	}
	immutableWorld := server.MakeImmutableWorld(RemoteWorld)

	calculatedItems := server.CalculateNextStage(0, p.ImageHeight, 0, p.ImageWidth, p, immutableWorld)

	RemoteWorld = calculatedItems.CalculatedWorld
	remoteResponse.AliveCells = calculatedItems.AliveCellsCount
	remoteResponse.World = calculatedItems.CalculatedWorld

	return nil
}

// QuitingServer is used to quit remote server
func (r *Remote) QuitingServer(kStatus commstruct.KStatus, kQuitting *commstruct.KQuitting) error {
	if kStatus.Status == true {
		os.Exit(0)
	}
	return nil
}

func main() {
	// listener, _ := net.Listen("tcp", ":8030")
	// rpc.Register(&Remote{})
	// rpc.Accept(listener)

	pAddr := flag.String("port", "8030", "port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	rpc.Register(&Remote{})
	listener, _ := net.Listen("tcp", ":"+*pAddr)
	rpc.Accept(listener)
}
