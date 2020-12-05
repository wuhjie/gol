package server

import (
	"flag"
	"fmt"
	"gol/distribute/gameLogic/gol"
	"math/rand"
	"net"
	"net/rpc"
	"time"
)

type serverstruct {
	c gol.DistributorChannels
}

func (server *server) withCommand() {
	
}

// after all round been executed
func (server *server) roundEnds() {

}

// implementing basic calculation on aws
func (server *server) calculationRunning() {
	
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
