package thread

import "github.com/zengqiang96/threading/plotter"

type ThreadPlotter struct {
	segmentDrawnNumber int
	Plotter            plotter.Plotter
	ThreadComputer     *ThreadComputer
}

func NewThreadPlotter(p plotter.Plotter, computer *ThreadComputer) *ThreadPlotter {
	return &ThreadPlotter{
		Plotter:        p,
		ThreadComputer: computer,
	}
}

func (tp *ThreadPlotter) Plot() {
	if tp.segmentDrawnNumber == tp.ThreadComputer.GetSegmentsNumber() {
		return
	} else if tp.segmentDrawnNumber > tp.ThreadComputer.GetSegmentsNumber() {
		tp.segmentDrawnNumber = 0
	}

	drawFromScratch := tp.segmentDrawnNumber == 0
	if drawFromScratch {

	} else {
		tp.ThreadComputer.drawThread(tp.Plotter, tp.segmentDrawnNumber)
	}

	tp.segmentDrawnNumber = tp.ThreadComputer.GetSegmentsNumber()
}
