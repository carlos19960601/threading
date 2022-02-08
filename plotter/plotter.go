package plotter

import "github.com/zengqiang96/threading/common"

type PlotterInfo struct {
	BackgroundColor string
	Blur            int
}

type Plotter interface {
	Size() common.Size
	Resize()
	Initialize(info PlotterInfo)
	Finalize()
	DrawLines(lines []common.Line, color EColor, opacity float64, thickness float64)
	DrawBrokenLine(points []common.Point, color EColor, opacity float64, thickness float64)
	DrawPoints(points []common.Point, color string, diameter float64)
}

type PlotterBase struct {
}

func (pb *PlotterBase) DrawBrokenLine(points []common.Point, color EColor, opacity float64, thickness int) {
	lines := make([]common.Line, 0)
	for i := 0; i < len(points)-1; i++ {
		lines = append(lines, common.Line{
			From: points[i],
			To:   points[i+1],
		})
	}
}
