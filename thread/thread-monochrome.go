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

func (tm *ThreadMonochrome) adjustCanvasData(data []uint8) {
	computeAdjustedValue := func(rawValue uint8) uint8 {
		return rawValue / 2
	}

	nbPixels := len(data) / 4
	for i := 0; i < nbPixels; i++ {
		averageSourceValue := (data[4*i+0] + data[4*i+1] + data[4*i+2]) / 3
		adjustedValue := computeAdjustedValue(averageSourceValue)
		data[4*i+0] = adjustedValue
		data[4*i+1] = adjustedValue
		data[4*i+2] = adjustedValue
	}
}

func (tm *ThreadMonochrome) EnableSamplingFor(color plotter.EColor) {

}

func (tm *ThreadMonochrome) SampleCanvas(data []uint8, index int) uint8 {
	return data[index+0] // only check the red channel because the hidden canvas is in black and white
}
