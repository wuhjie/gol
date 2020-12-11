package main

import (
	"fmt"
	"os"
	"testing"

	"uk.ac.bris.cs/gameoflife/gol"
)

func Benchmark(b *testing.B) {
	os.Stdout = nil

	tests := []gol.Params{
		{ImageWidth: 16, ImageHeight: 16},
		{ImageWidth: 64, ImageHeight: 64},
		{ImageWidth: 512, ImageHeight: 512},
	}

	var threads = [16]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	for _, p := range tests {
		// for _, turns := range []int{100, 400, 700, 1000, 5000} {
		for _, turns := range []int{7000} {
			p.Turns = turns
			for _, thread := range threads {
				p.Threads = thread
				testName := fmt.Sprintf("%dx%dx%d-%d", p.ImageWidth, p.ImageHeight, p.Turns, p.Threads)
				b.Run(testName, func(t *testing.B) {
					for i := 0; i < b.N; i++ {
						events := make(chan gol.Event)
						gol.Run(p, events, nil)
					}
				})
			}
		}
	}
}
