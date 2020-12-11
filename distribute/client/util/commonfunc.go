package util

// InitialisedWorld makes new 2-Dimension array of byte
func InitialisedWorld(height, width int) [][]byte {
	world := make([][]byte, height)
	for i := range world {
		world[i] = make([]byte, width)
	}
	return world
}

// RemoteCalculation is calling methods, only the method is needed
var RemoteCalculation = "Remote.CalculateNextTurn"

// QuitingServer is calling remote localKPresses function
var QuitingServer = "Remote.QuitingServer"

// RemoteWorldRecieved is calling remote methods
var RemoteWorldRecieved = "Broker.WorldReceived"
