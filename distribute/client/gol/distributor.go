package gol

import (
	"uk.ac.bris.cs/gameoflife/util"

	"strconv"
	"strings"
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

func LocalFilesReading(u UserChannels, p Params, c DistributorChannels) {
	world := InitialisedWorld(p.ImageHeight, p.ImageWidth)

	c.IoCommand <- ioInput
	c.IoFilename <- strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight)}, "x")

	//flipped the initial alive cells
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			val := <-c.IoInput
			world[y][x] = val
		}
	}

	u.InitialWorld <- world
}
