package main

import "github.com/zengqiang96/threading/thread"

const MAX_COMPUTING_TIME_PER_FRAME = 20 // ms

func main() {
	threadPlotter := thread.NewThreadPlotter()
	threadComputer := thread.NewThreadComputer()

	for threadComputer.ComputeNextSegments(MAX_COMPUTING_TIME_PER_FRAME) {
		threadPlotter.Plot()
	}
}
