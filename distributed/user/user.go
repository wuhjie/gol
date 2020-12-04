package main

import (
	"flag"
	"fmt"
	"net/rpc"
	"uk.ac.bris.cs/gameoflife/gol"
)

// user structure; act as a client
type UserChannels struct{
	userKeyPresses <-chan rune
	userTempWorld chan<- uint8
	userEvent chan<- gol.Event
	userParams chan gol.Params
	distributedChannels chan<- gol.DistributorChannels
}

func control(u UserChannels, keypress chan rune) {
	// only send the command to distributor when command is correct
	select {
	case command := <-u.userKeyPresses:
		if command == 's' || command == 'p' || command == 'q' {
			keypress <- command
		}else {
			fmt.Println("please enter the correct command")
		}

	}
}

// to start the game as a user
func launchGameEngine(u UserChannels) {
	gol.Run(<-u.userParams, u.userEvent, u.userKeyPresses )
}

func main() {
	server := flag.String("server", "127.0.0.1:8030", "IP: port string to connect to as server")
	flag.Parse()
	client, _ := rpc.Dial("tcp", *server)
	fmt.Println("GAME STARTS...")
	defer client.Close()

}