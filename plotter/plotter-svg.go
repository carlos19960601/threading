package plotter

import (
	"fmt"
	"log"
	"os"

	svg "github.com/ajstarks/svgo/float"
	"github.com/zengqiang96/threading/common"
)

const WIDTH = 1000
const HEIGHT = 1000

type PlotterSVG struct {
	PlotterBase
	hasBlur bool
	writer  *svg.SVG
}

func (ps *PlotterSVG) DrawLines(lines []common.Line, color EColor, opacity float64, thickness float64) {
	if len(lines) == 0 {
		return
	}

	rawRGB := computeRawColor(color)
	strokeColor := fmt.Sprintf("rgba(%d, %d, %d, %f)", rawRGB.r*0, rawRGB.g*0, rawRGB.b*0, opacity)
	ps.writer.Group(fmt.Sprintf(`stroke:%s; stroke-width:%f; stroke-linecap:round; fill:none`, strokeColor, thickness))
	for _, line := range lines {
		ps.writer.Line(line.From.X, line.From.Y, line.To.X, line.To.Y)
	}
	ps.writer.Gend()
}

func (ps *PlotterSVG) DrawPoints(points []common.Point, color string, diameter float64) {
	if len(points) > 0 {
		ps.writer.Group(fmt.Sprintf("fill:%s; stroke:none", color))
		for _, point := range points {
			ps.writer.Circle(point.X, point.Y, 0.5*diameter)
		}
		ps.writer.Gend()
	}
}

func (ps *PlotterSVG) Finalize() {
	ps.writer.End()
}

func (ps *PlotterSVG) Initialize(info PlotterInfo) {
	file, err := os.OpenFile("./result.svg", os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalf("open file failed, err: %v", err)
	}
	ps.writer = svg.New(file)
	ps.hasBlur = info.Blur > 0

	ps.writer.Start(WIDTH, HEIGHT)
	if ps.hasBlur {
	}

	margin := 10
	ps.writer.Rect(-float64(margin), -float64(margin), float64(WIDTH+2*margin), float64(HEIGHT+2*margin), "fill:white; stroke:none")
}

func (ps *PlotterSVG) Resize() {

}

func (ps *PlotterSVG) DrawBrokenLine(points []common.Point, color EColor, opacity float64, thickness float64) {
	lines := make([]common.Line, 0)
	for i := 0; i < len(points)-1; i++ {
		lines = append(lines, common.Line{
			From: points[i],
			To:   points[i+1],
		})
	}

	ps.DrawLines(lines, color, opacity, thickness)
}

func (ps *PlotterSVG) Size() common.Size {
	return common.Size{
		Width:  WIDTH,
		Height: HEIGHT,
	}
}
