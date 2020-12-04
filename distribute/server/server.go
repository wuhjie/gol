package main

import (
	"flag"
	"math/rand"
	"net"
	"net/rpc"
	"time"
)

type service struct {

}

func (service *service) withCommand() {

}

// after all round been executed
func (service *service) roundEnds() {

}

func main () {
	pAddr := flag.String("port", "8030", "port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	//rpc.Register()
	listener, _ := net.Listen("tcp", ":"+ *pAddr)

	defer listener.Close()
	rpc.Accept(listener)
}
