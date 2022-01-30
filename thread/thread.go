package thread

import (
	"github.com/zengqiang96/threading/plotter"
)

type Thread2Grow struct {
	thread []Peg
	color  plotter.EColor
}

type Thread interface {
	GetTotalSegmentNumber() int
	GetThread2Grow() Thread2Grow
}

type ThreadBase struct {
}

func (tb *ThreadBase) ComputeSegmentNumber(pegs []Peg) int {
	if len(pegs) > 1 {
		return len(pegs) - 1
	}

	return 0
}
