package main

import (
	"image"
	"log"
	"os"

	"github.com/zengqiang96/threading/common"
	"github.com/zengqiang96/threading/plotter"
	"github.com/zengqiang96/threading/thread"
)

const MAX_COMPUTING_TIME_PER_FRAME = 20 // ms

func main() {

	file, err := os.OpenFile("", 777, os.ModePerm)
	if err != nil {
		log.Fatalf("open file failed, err: %v", err)
	}
	sourceImage, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("image decode failed, err: %v", err)
	}

	svgPlotter := plotter.PlotterSVG{}
	threadComputer := thread.NewThreadComputer(sourceImage, common.Config{})
	threadPlotter := thread.NewThreadPlotter(svgPlotter, threadComputer)

	for threadComputer.ComputeNextSegments(MAX_COMPUTING_TIME_PER_FRAME) {
		threadPlotter.Plot()
	}
}
