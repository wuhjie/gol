package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"time"
)

type server struct {
	gameStatus bool
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
