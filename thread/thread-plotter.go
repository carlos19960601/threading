package thread

import (
	"github.com/zengqiang96/threading/common"
	"github.com/zengqiang96/threading/plotter"
)

type ThreadPlotter struct {
	segmentDrawnNumber int
	config             common.Config
	Plotter            plotter.Plotter
	ThreadComputer     *ThreadComputer
}

func NewThreadPlotter(p plotter.Plotter, computer *ThreadComputer) *ThreadPlotter {
	return &ThreadPlotter{
		Plotter:        p,
		ThreadComputer: computer,
		config:         computer.Config,
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
		plotterInfos := plotter.PlotterInfo{
			BackgroundColor: "white",
			Blur:            0,
		}
		tp.Plotter.Resize()
		tp.Plotter.Initialize(plotterInfos)

		if tp.config.DisplayPegs {
			tp.ThreadComputer.DrawPegs(tp.Plotter)
		}

		tp.ThreadComputer.drawThread(tp.Plotter, 0)
		// tp.Plotter.Finalize()
	} else {
		tp.ThreadComputer.drawThread(tp.Plotter, tp.segmentDrawnNumber)
	}

	tp.segmentDrawnNumber = tp.ThreadComputer.GetSegmentsNumber()
}
