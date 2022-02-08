package thread

import (
	"github.com/zengqiang96/threading/plotter"
)

type Thread2Grow struct {
	thread *[]Peg
	color  plotter.EColor
}

type Thread interface {
	adjustCanvasData(data []uint8)
	GetTotalSegmentNumber() int
	GetThread2Grow() Thread2Grow
	EnableSamplingFor(color plotter.EColor)
	SampleCanvas(data []uint8, index int) uint8
}

type ThreadBase struct {
}

func (tb *ThreadBase) ComputeSegmentNumber(pegs []Peg) int {
	if len(pegs) > 1 {
		return len(pegs) - 1
	}

	return 0
}

func SliceContains(pegs []Peg, peg Peg) bool {
	for _, p := range pegs {
		if p == peg {
			return true
		}
	}

	return false
}
