package gol

import (
	"strconv"
	"strings"

	"uk.ac.bris.cs/gameoflife/client/util"
)

// initialisedWorld is used to make 2-D world
func initialisedWorld(height, width int) [][]byte {
	world := make([][]byte, height)
	for i := range world {
		world[i] = make([]byte, width)
	}
	return world
}

// InputWorldImage is related to loading images from the io channel
func InputWorldImage(p Params, c DistributorChannels) [][]byte {
	world := initialisedWorld(p.ImageHeight, p.ImageWidth)
	c.IoCommand <- ioInput
	c.IoFilename <- strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight)}, "x")

	//adding the values in ioInput channel to initialised world inside distributor
	//flipped the initial alive cells
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			val := <-c.IoInput
			world[y][x] = val
			if val == alive {
				c.Events <- CellFlipped{CompletedTurns: 0, Cell: struct{ X, Y int }{X: x, Y: y}}
			}
		}
	}
	return world
}

// OutputWorldImage sends the world into the IoOutput channel
func OutputWorldImage(c DistributorChannels, p Params, world [][]byte) {
	c.IoCommand <- ioOutput
	filename := strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight), strconv.Itoa(c.CompletedTurns)}, "x")
	c.IoFilename <- filename

	for m := 0; m < p.ImageHeight; m++ {
		for n := 0; n < p.ImageWidth; n++ {
			c.IoOutput <- world[m][n]
		}
	}
	c.Events <- ImageOutputComplete{c.CompletedTurns, filename}
}

// CalculateAliveCells the alive cells in current round
func CalculateAliveCells(p Params, world [][]byte) []util.Cell {
	var aliveCells []util.Cell

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			if world[y][x] == alive {
				aliveCells = append(aliveCells, util.Cell{X: x, Y: y})
			}
		}
	}
	return aliveCells
}
