package client

import (
	"flag"
	"net/rpc"
)



func main() {
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	flag.Parse()
	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()

	//file, _ := os.Open()

}