package plotter

import "github.com/zengqiang96/threading/common"

type PlotterSVG struct {
	PlotterBase
	hasBlur bool
}

func (ps *PlotterSVG) DrawLines(lines []common.Line, color EColor, opacity float64, thickness int) {

}

func (ps *PlotterSVG) Finalize() {}

func (ps *PlotterSVG) Initialize(info PlotterInfo) {}

func (ps *PlotterSVG) Resize() {

}
