package main

import (
	"net"
	"net/rpc"
	"os"

	"uk.ac.bris.cs/gameoflife/comm"
)

// BrokerWorld is the initial world from the local machine
var BrokerWorld [][]byte

// Threads is the number of nodes
var Threads int

// PicHeight is the hight of picture of local machine
var PicHeight int

// PicWidth is the width of picture of local machine
var PicWidth int

// QStatus tells if the last version of this world is saved
var QStatus bool

// BrokerTurn is turn on broker
var BrokerTurn int

// world is the same as [][]byte, which is used to initialise array
type world [][]byte

// initialisedWorld is used to make 2-D world
func initialisedWorld(height, width int) [][]byte {
	world := make([][]byte, height)
	for i := range world {
		world[i] = make([]byte, width)
	}
	return world
}

// Broker is considered to be the remote manager
type Broker struct{}

// GetQStatus sends the status of q before any other behaviour
func (b *Broker) GetQStatus(req comm.CommonMsg, res *comm.QStatus) error {
	res.Status = QStatus
	return nil
}

// ModifyQStatus changes the q status when q pressed
func (b *Broker) ModifyQStatus(req comm.CommonMsg, res *comm.CommonMsg) error {
	res.Msg = "q status has been changed"
	QStatus = true
	return nil
}

// GetBrokerStatus gives the saving world on broker
func (b *Broker) GetBrokerStatus(req comm.CommonMsg, res *comm.BrokerSaved) error {
	res.World = BrokerWorld
	res.Threads = Threads
	res.ImageWidth = PicWidth
	res.ImageHeight = PicHeight
	res.Turns = BrokerTurn
	return nil
}

// WorldReceived is used to receive global variables
func (b *Broker) WorldReceived(initial comm.BrokerRequest, rep *comm.ResponseOnReceivedWorld) error {
	BrokerWorld = initial.World
	Threads = initial.Threads
	PicHeight = initial.ImageHeight
	PicWidth = initial.ImageWidth
	QStatus = false
	rep.Msg = "world received from local machine"

	return nil
}

// Calculate sent calculation request to remote factory
func (b *Broker) Calculate(req comm.Localsent, res *comm.BrokerReturn) error {
	BrokerTurn = req.Turns
	// heightPerThread := PicHeight / Threads
	flippedCell := []comm.Cell{}
	// PortList := [2]string{":8050", ":8060"}

	tempWorld := make([]world, Threads)

	clientone, _ := rpc.Dial("tcp", "127.0.0.1:8050")
	workerRequestone := comm.WorkerRequest{
		StartX: 0,
		EndX:   PicWidth,
		StartY: 0,
		EndY:   PicHeight,
		World:  BrokerWorld,
	}
	workerReplyone := new(comm.WorkerReply)

	clientone.Call("Remote.CalculateNextTurn", workerRequestone, workerReplyone)
	tempWorld[0] = workerReplyone.PartWorld
	flippedCell = append(flippedCell, workerReplyone.ChangedCells...)
	defer clientone.Close()

	// clienttwo, _ := rpc.Dial("tcp", "127.0.0.1:8060")
	// workerRequesttwo := comm.WorkerRequest{
	// 	StartX: 0,
	// 	EndX:   PicWidth,
	// 	StartY: 1 / 2 * PicHeight,
	// 	EndY:   PicHeight,
	// 	World:  BrokerWorld,
	// }
	// workerReplytwo := new(comm.WorkerReply)
	// clienttwo.Call("Remote.CalculateNextTurn", workerRequesttwo, workerReplytwo)
	// tempWorld[1] = workerReplytwo.PartWorld

	// flippedCell = append(flippedCell, workerReplytwo.ChangedCells...)
	// defer clienttwo.Close()

	mergedWorld := make(world, 0)
	// for i := 0; i < Threads; i++ {
	// 	mergedWorld = append(mergedWorld, tempWorld[i]...)
	// }
	// mergedWorld = append(append(mergedWorld, tempWorld[0]...), tempWorld[1]...)
	mergedWorld = append(mergedWorld, tempWorld[0]...)

	BrokerWorld = mergedWorld
	res.World = BrokerWorld
	res.ChangedCells = flippedCell

	return nil
}

// QuittingBroker is used to quit broker and sent command to quit factory
func (b *Broker) QuittingBroker(req comm.KStatus, res *comm.CommonMsg) error {
	if req.Status == true {
		// need to use the public address of aws
		// client, _ := rpc.Dial("tcp", "127.0.0.1:8050")
		client, _ := rpc.Dial("tcp", "3.238.200.91:8050")
		kq := new(comm.KQuitting)
		ks := comm.KStatus{
			Status: true,
		}
		client.Call("Remote.QuitingFactory", ks, kq)
		res.Msg = "Quitting command has sent to factory"
		os.Exit(0)
	}

	return nil
}

func main() {
	listener, _ := net.Listen("tcp", "127.0.0.1:8030")
	defer listener.Close()
	rpc.Register(&Broker{})
	rpc.Accept(listener)
}
