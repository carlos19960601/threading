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
