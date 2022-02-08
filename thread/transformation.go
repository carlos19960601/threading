package thread

import (
	"math"

	"github.com/tfriedel6/canvas"
	"github.com/zengqiang96/threading/common"
)

type Transformation struct {
	Scaling float64
	Origin  common.Point
}

func NewTransformation(frameSize common.Size, elementSize *canvas.Canvas) *Transformation {
	scaleToFitWidth := frameSize.Width / float64(elementSize.Width())
	scaleToFitHeight := frameSize.Height / float64(elementSize.Height())

	scaling := math.Min(scaleToFitWidth, scaleToFitHeight)
	return &Transformation{
		Scaling: scaling,
		Origin: common.Point{
			X: 0.5 * (frameSize.Width - scaling*float64(elementSize.Width())),
			Y: 0.5 * (frameSize.Height - scaling*float64(elementSize.Height())),
		},
	}
}

func (tr *Transformation) Transform(point common.Point) common.Point {
	return common.Point{
		X: tr.Origin.X + point.X*tr.Scaling,
		Y: tr.Origin.Y + point.Y*tr.Scaling,
	}
}
