package gol

import (
	"strconv"
	"strings"

	"uk.ac.bris.cs/gameoflife/util"
)

// DistributorChannels contains things that need for parallel calculation
type DistributorChannels struct {
	Events          chan<- Event //Events is what communicate with SDL
	IoCommand       chan<- ioCommand
	IoIdle          <-chan bool
	IoFilename      chan<- string
	AliveCellsCount chan<- []util.Cell
	IoInput         <-chan uint8
	IoOutput        chan<- uint8
	CompletedTurns  int
	KeyPresses      <-chan rune
}

// Distributor is calculated related
func Distributor(u UserChannels, p Params, c DistributorChannels) {
	world := initialisedWorld(p.ImageHeight, p.ImageWidth)

	c.IoCommand <- ioInput
	c.IoFilename <- strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight)}, "x")

	//adding the values in IoInput channel to initialised world inside distributor
	//flipped the initial alive cells
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			val := <-c.IoInput
			world[y][x] = val
		}
	}

	u.initialWorld <- world

	// senging 2-D array to aws

}
