package thread

import "github.com/zengqiang96/threading/plotter"

type ThreadMonochrome struct {
	ThreadBase
	threadPegs []Peg
}

func NewThreadMonochrome() *ThreadMonochrome {
	return &ThreadMonochrome{
		threadPegs: make([]Peg, 0),
	}
}

func (tm *ThreadMonochrome) GetTotalSegmentNumber() int {
	return tm.ComputeSegmentNumber(tm.threadPegs)
}

func (tm *ThreadMonochrome) GetThread2Grow() Thread2Grow {
	return Thread2Grow{
		thread: &tm.threadPegs,
		color:  plotter.MONOCHROME,
	}
}

func (tm *ThreadMonochrome) AdjustCanvasData(data []uint8) {
	computeAdjustedValue := func(rawValue float64) float64 {
		return rawValue / 2
	}

	nbPixels := len(data) / 4
	for i := 0; i < nbPixels; i++ {
		averageSourceValue := float64(int(data[4*i+0])+int(data[4*i+1])+int(data[4*i+2])) / 3
		adjustedValue := computeAdjustedValue(averageSourceValue)
		data[4*i+0] = uint8(adjustedValue)
		data[4*i+1] = uint8(adjustedValue)
		data[4*i+2] = uint8(adjustedValue)
	}
}

func (tm *ThreadMonochrome) EnableSamplingFor(color plotter.EColor) {

}

func (tm *ThreadMonochrome) SampleCanvas(data []uint8, index int) uint8 {
	return data[index+0] // only check the red channel because the hidden canvas is in black and white
}

func (tm *ThreadMonochrome) IterateOnThreads(nbSegmentsToIgnore int, callback ThreadsIterator) {
	tm.IterateOnThread(tm.threadPegs, plotter.MONOCHROME, nbSegmentsToIgnore, callback)
}
