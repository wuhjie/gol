package controllerrelated

import (
	"fmt"
	gol2 "uk.ac.bris.cs/gameoflife/gol"
)

func SdlController(c gol2.DistributorChannels, world [][]byte, qstatus bool, p Params) {
	select {
	case command := <-c.KeyPresses:

		if command == 's' {
			c.Events <- gol2.StateChange{CompletedTurns: c.CompletedTurns, NewState: gol2.Executing}
			gol2.OutputWorldImage(c, p, world)
		}
		if command == 'q' {
			c.Events <- gol2.StateChange{CompletedTurns: c.CompletedTurns, NewState: gol2.Quitting}
			qStatus = true
		}
		if command == 'p' {
			c.Events <- gol2.StateChange{CompletedTurns: c.CompletedTurns, NewState: gol2.Paused}

			gol2.OutputWorldImage(c, p, world)
			for {
				command := <-c.KeyPresses
				if command == 'p' {
					fmt.Println("Continuing")
					c.Events <- gol2.StateChange{CompletedTurns: c.CompletedTurns, NewState: gol2.Executing}
					c.Events <- gol2.TurnComplete{CompletedTurns: c.CompletedTurns}
				}
				break
			}
		}
	default:
	}

}
