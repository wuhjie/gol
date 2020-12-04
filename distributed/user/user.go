package user

import (
	"flag"
	"net/rpc"
	"uk.ac.bris.cs/gameoflife/gol"
)

// user structure; act as a client
type userChannels struct{
	keyPresses <-chan rune
	tempWorld chan<- uint8
	event chan<- gol.Event
}


func main() {
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	flag.Parse()
	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()


}