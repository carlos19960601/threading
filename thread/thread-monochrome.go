package thread

import "github.com/zengqiang96/threading/plotter"

type ThreadMonochrome struct {
	ThreadBase
	threadPegs []Peg
}

func (tm *ThreadMonochrome) GetTotalSegmentNumber() int {
	return tm.ComputeSegmentNumber(tm.threadPegs)
}

func (tm *ThreadMonochrome) GetThread2Grow() Thread2Grow {
	return Thread2Grow{
		thread: tm.threadPegs,
		color:  plotter.MONOCHROME,
	}
}
