package main

import (
	"fmt"
	"net"
	"net/rpc"

	"uk.ac.bris.cs/gameoflife/commstruct"
)

// Below are global variables stored in remote server
var BrokerWorld [][]byte
var Threads int
var PicHeight int
var PicWidth int

type world [][]byte

// var FlippedCell []commstruct.Cell
func initialisedWorld(height, width int) [][]byte {
	world := make([][]byte, height)
	for i := range world {
		world[i] = make([]byte, width)
	}
	return world
}

// Broker is considered to be the remote manager
type Broker struct{}

// WorldReceived is used to receive global variables
func (b *Broker) WorldReceived(initial commstruct.BrokerRequest, rep *commstruct.ResponseOnReceivedWorld) error {
	BrokerWorld = initial.World
	Threads = initial.Threads
	PicHeight = initial.ImageHeight
	PicWidth = initial.ImageWidth

	rep.Msg = "world received from local machine"

	return nil
}

// Calculate establishes connections
func (b *Broker) Calculate(req commstruct.BrokerConnection, res *commstruct.BrokerReturn) error {

	heightPerThread := PicHeight / Threads
	flippedCell := []commstruct.Cell{}
	PortList := [2]string{"8050", "8060"}

	tempWorld := make([]world, Threads)

	for i := 0; i < Threads-1; i++ {
		client, _ := rpc.Dial("tcp", "127.0.0.1:"+PortList[i])
		workerRequest := commstruct.WorkerRequest{
			StartX: 0,
			EndX:   PicWidth,
			StartY: i * heightPerThread,
			EndY:   (i + 1) * heightPerThread,
			World:  BrokerWorld,
		}
		workerReply := new(commstruct.WorkerReply)

		client.Call("Remote.CalculateNextTurn", workerRequest, workerReply)
		tempWorld[i] <- workerReply.PartWorld
		flippedCell = append(flippedCell, workerReply.ChangedCells...)
		defer client.Close()
	}
	client, _ := rpc.Dial("tcp", "127.0.0.1:"+PortList[Threads-1])
	workerRequest := commstruct.WorkerRequest{
		StartX: 0,
		EndX:   PicWidth,
		StartY: (Threads - 1) * heightPerThread,
		EndY:   PicHeight,
		World:  BrokerWorld,
	}
	workerReply := new(commstruct.WorkerReply)
	client.Call("Remote.CalculateNextTurn", workerRequest, workerReply)
	tempWorld[Threads-1] = workerReply.PartWorld

	fmt.Printf("tempWorld[Threads-1] length: %d", len(tempWorld[Threads-1]))

	flippedCell = append(flippedCell, workerReply.ChangedCells...)
	defer client.Close()

	mergedWorld := make(world, 0)
	for i := 0; i < Threads; i++ {
		mergedWorld = append(mergedWorld, tempWorld[i]...)
	}

	fmt.Printf("mergedWorld length: %d", len(mergedWorld))

	BrokerWorld = mergedWorld
	res.World = BrokerWorld
	res.ChangedCells = flippedCell

	return nil
}

func main() {
	listener, _ := net.Listen("tcp", ":8030")
	defer listener.Close()
	rpc.Register(&Broker{})
	rpc.Accept(listener)
}
