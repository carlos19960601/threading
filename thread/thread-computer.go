package thread

import (
	"image"
	"time"

	"github.com/zengqiang96/threading/common"
	"github.com/zengqiang96/threading/plotter"
)

type Peg struct {
	x float64
	y float64
}

type Segment struct {
	peg1 Peg
	peg2 Peg
}

type ThreadComputer struct {
	Config      common.Config
	SourceImage image.Image
	pegs        []Peg
	thread      Thread
}

func NewThreadComputer(sourceImage image.Image, config common.Config) *ThreadComputer {
	tc := &ThreadComputer{
		SourceImage: sourceImage,
	}

	return tc
}

func (tc *ThreadComputer) ComputeNextSegments(maxMsTaken int64) bool {

	start := time.Now()
	targetSegmentNumber := tc.Config.LineNumber
	if tc.GetSegmentsNumber() == targetSegmentNumber {
		return false
	} else if tc.GetSegmentsNumber() > targetSegmentNumber {
		// TODO
	}

	for tc.GetSegmentsNumber() < targetSegmentNumber && time.Since(start).Milliseconds() < maxMsTaken {
		thread2Grow := tc.thread.GetThread2Grow()
		tc.computeSegment(thread2Grow.thread)
	}
	return true
}

func (tc *ThreadComputer) GetSegmentsNumber() int {
	return tc.thread.GetTotalSegmentNumber()
}

func (tc *ThreadComputer) computeSegment(thread []Peg) {
	var lastPeg, nextPeg Peg
	if len(thread) == 0 {

	}
}

func (tc *ThreadComputer) drawThread(p plotter.Plotter, segmentsToIgnoreNumber int) {

}
