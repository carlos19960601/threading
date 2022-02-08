package thread

import (
	"image"
	"image/png"
	"math"
	"os"
	"time"

	"github.com/zengqiang96/threading/common"
	"github.com/zengqiang96/threading/plotter"

	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/softwarebackend"
	"github.com/valyala/fastrand"
)

const (
	MIN_SAFE_NUMBER = -9007199254740991
	TWO_PI          = 2 * math.Pi
)

type Peg interface {
	GetX() float64
	GetY() float64
	Index() int
}

type peg struct {
	X float64
	Y float64
}

func (p *peg) GetX() float64 {
	return p.X
}

func (p *peg) GetY() float64 {
	return p.Y
}

type Segment struct {
	peg1 Peg
	peg2 Peg
}

type pegCircle struct {
	peg
	angle float64
	index int
}

func (pc *pegCircle) GetX() float64 {
	return pc.peg.X
}

func (pc *pegCircle) GetY() float64 {
	return pc.peg.Y
}

func (pc *pegCircle) Index() int {
	return pc.index
}

type ThreadComputer struct {
	Config              common.Config
	SourceImage         image.Image
	pegs                []Peg
	thread              Thread
	lineOpacity         float64
	lineOpacityInternal float64
	lineThickness       float64
	arePegsTooClose     func(peg1, peg2 interface{}) bool
	canvasBackend       *softwarebackend.SoftwareBackend
	hiddenCanvasData    *image.RGBA
	hiddenCanvasContext *canvas.Canvas
	hiddenCanvasScale   int
}

func NewThreadComputer(sourceImage image.Image, config common.Config) *ThreadComputer {
	tc := &ThreadComputer{
		SourceImage:       sourceImage,
		Config:            config,
		hiddenCanvasScale: config.Quality,
	}

	tc.Reset(math.Pow(2, float64(config.LinesOpacity-7)), config.LinesThickness)
	return tc
}

func (tc *ThreadComputer) Reset(opacity float64, lineThickness float64) {
	tc.lineOpacity = opacity
	tc.lineThickness = lineThickness

	tc.thread = &ThreadMonochrome{}
	tc.resetHiddenCanvas()
	tc.pegs = tc.computePegs()
}

func (tc *ThreadComputer) resetHiddenCanvas() {
	wantedSize := computeBestSize(tc.SourceImage, 100*tc.hiddenCanvasScale)
	backend := softwarebackend.New(int(wantedSize.Width), int(wantedSize.Height))
	tc.canvasBackend = backend
	tc.hiddenCanvasContext = canvas.New(backend)

	tc.hiddenCanvasContext.DrawImage(tc.SourceImage, 0, 0, wantedSize.Width, wantedSize.Height)
	imageData := tc.hiddenCanvasContext.GetImageData(0, 0, int(wantedSize.Width), int(wantedSize.Height))
	tc.thread.AdjustCanvasData(imageData.Pix)
	tc.hiddenCanvasContext.PutImageData(imageData, 0, 0)
	tc.computeError()
	tc.initializeHiddenCanvasLineProperties()
}

func (tc *ThreadComputer) computeError() {
	tc.uploadCanvasDataToCPU()
}

func (tc *ThreadComputer) uploadCanvasDataToCPU() {
	if tc.hiddenCanvasData == nil {
		width := tc.hiddenCanvasContext.Width()
		height := tc.hiddenCanvasContext.Height()
		tc.hiddenCanvasData = tc.hiddenCanvasContext.GetImageData(0, 0, width, height)
	}
}

func (tc *ThreadComputer) initializeHiddenCanvasLineProperties() {
	theoricalThickness := tc.lineThickness * float64(tc.hiddenCanvasScale)
	if theoricalThickness <= 1 {
		tc.lineOpacityInternal = 0.5 * tc.lineOpacity * theoricalThickness
		tc.hiddenCanvasContext.SetLineWidth(1)
	} else {
		tc.lineOpacityInternal = 0.5 * tc.lineOpacity
		tc.hiddenCanvasContext.SetLineWidth(theoricalThickness)
	}
}

func (tc *ThreadComputer) ComputeNextSegments(maxMsTaken int64) bool {
	start := time.Now()
	targetSegmentNumber := tc.Config.LineNumber
	if tc.GetSegmentsNumber() == targetSegmentNumber {
		return false
	}

	var lastColor plotter.EColor = -1
	for tc.GetSegmentsNumber() < targetSegmentNumber && time.Since(start).Milliseconds() < maxMsTaken {
		thread2Grow := tc.thread.GetThread2Grow()
		if lastColor != thread2Grow.color {
			plotter.ApplyCanvasCompositing(tc.hiddenCanvasContext, thread2Grow.color, tc.lineOpacityInternal)
			tc.thread.EnableSamplingFor(thread2Grow.color)
			lastColor = thread2Grow.color
		}
		tc.computeSegment(thread2Grow.thread)
	}
	return true
}

func (tc *ThreadComputer) GetSegmentsNumber() int {
	return tc.thread.GetTotalSegmentNumber()
}

func (tc *ThreadComputer) computeSegment(thread *[]Peg) {
	var lastPeg, nextPeg Peg
	if len(*thread) == 0 {
		startingSegment := tc.computeBestStartingSegment()
		*thread = append(*thread, startingSegment.peg1)
		lastPeg = startingSegment.peg1
		nextPeg = startingSegment.peg2
	} else {
		lastPeg = (*thread)[len(*thread)-1]
		prevousPegs := *thread
		if len(*thread) > 20 {
			prevousPegs = (*thread)[len(*thread)-20:]
		}
		nextPeg = tc.computeBestNextPeg(lastPeg, prevousPegs)
	}

	*thread = append(*thread, nextPeg)
	tc.drawSegmentOnHiddenCanvas(lastPeg, nextPeg)
}

func (tc *ThreadComputer) drawSegmentOnHiddenCanvas(peg1, peg2 Peg) {
	tc.hiddenCanvasContext.BeginPath()
	tc.hiddenCanvasContext.MoveTo(peg1.GetX(), peg1.GetY())
	tc.hiddenCanvasContext.LineTo(peg2.GetX(), peg2.GetY())
	tc.hiddenCanvasContext.Stroke()
	tc.hiddenCanvasContext.ClosePath()

	tc.OutputCanvas()
	// invalidate CPU data
	tc.hiddenCanvasData = nil
}

func (tc *ThreadComputer) computeBestNextPeg(currentPeg Peg, pegsToAvoid []Peg) Peg {
	candidates := make([]Peg, 0)
	var bestScore float64 = MIN_SAFE_NUMBER

	for _, peg := range tc.pegs {
		if !tc.arePegsTooClose(currentPeg, peg) && !SliceContains(pegsToAvoid, peg) {
			condidateScore := tc.computeSegmentPotential(currentPeg, peg)
			if condidateScore > float64(bestScore) {
				bestScore = condidateScore
				candidates = []Peg{peg}
			} else if condidateScore == float64(bestScore) {
				candidates = append(candidates, peg)
			}
		}
	}

	return randomPeg(candidates)
}

func (tc *ThreadComputer) computeBestStartingSegment() Segment {
	candidates := make([]Segment, 0)
	var bestScore float64 = MIN_SAFE_NUMBER

	step := 1 + math.Floor(float64(len(tc.pegs))/100)

	for iPegId1 := 0; iPegId1 < len(tc.pegs); iPegId1 += int(step) {
		for iPegId2 := iPegId1 + 1; iPegId2 < len(tc.pegs); iPegId2 += int(step) {
			peg1, peg2 := tc.pegs[iPegId1], tc.pegs[iPegId2]
			if !tc.arePegsTooClose(peg1, peg2) {
				candidateScore := tc.computeSegmentPotential(peg1, peg2)
				if candidateScore > bestScore {
					bestScore = candidateScore
					candidates = []Segment{{peg1: peg1, peg2: peg2}}
				} else if candidateScore == bestScore {
					candidates = append(candidates, Segment{peg1: peg1, peg2: peg2})
				}
			}
		}
	}

	return randomSegment(candidates)
}

func (tc *ThreadComputer) drawThread(p plotter.Plotter, segmentsToIgnoreNumber int) {
	transformation := tc.computeTransformation(p.Size())
	lineWidth := (transformation.Scaling * float64(tc.hiddenCanvasScale)) * float64(tc.lineThickness)
	tc.thread.IterateOnThreads(segmentsToIgnoreNumber, func(thread []Peg, color plotter.EColor) {
		points := make([]common.Point, 0)
		for _, p := range thread {
			points = append(points, transformation.Transform(common.Point{
				X: p.GetX(),
				Y: p.GetY(),
			}))
		}

		p.DrawBrokenLine(points, color, tc.lineOpacity, lineWidth)
	})
}

func (tc *ThreadComputer) computeTransformation(targetSize common.Size) *Transformation {
	return NewTransformation(targetSize, tc.hiddenCanvasContext)
}

func (tc *ThreadComputer) computePegs() []Peg {
	width := tc.hiddenCanvasContext.Width()
	height := tc.hiddenCanvasContext.Height()

	defaultSize := 1000

	var domainSize common.Size
	aspectRatio := float64(width) / float64(height)
	if aspectRatio > 1 {
		domainSize = common.Size{
			Width:  float64(defaultSize),
			Height: math.Round(float64(defaultSize) / aspectRatio),
		}
	} else {
		domainSize = common.Size{
			Width:  math.Round(float64(defaultSize) * aspectRatio),
			Height: float64(defaultSize),
		}
	}

	tc.arePegsTooClose = func(p1, p2 interface{}) bool {
		peg1, peg2 := p1.(*pegCircle), p2.(*pegCircle)
		absDeltaAngle := math.Abs(peg1.angle - peg2.angle)
		minAngle := math.Min(absDeltaAngle, TWO_PI-absDeltaAngle)
		return minAngle <= TWO_PI/16
	}
	halfWidth := 0.5 * domainSize.Width
	halfHeight := 0.5 * domainSize.Height
	// 拉马努金 椭圆周长估计
	circumference := math.Pi * (3*(halfWidth+halfHeight) - math.Sqrt((3*halfWidth+halfHeight)*(3*halfHeight+halfWidth)))

	distanceBetweenPegs := circumference / float64(tc.Config.PegsCount)
	var angle float64 = 0
	for len(tc.pegs) < tc.Config.PegsCount {
		cosAngle := math.Cos(angle)
		sinAngle := math.Sin(angle)

		peg := &pegCircle{
			peg: peg{
				X: halfWidth * (1 + cosAngle),
				Y: halfHeight * (1 + sinAngle),
			},
			angle: angle,
		}

		tc.pegs = append(tc.pegs, peg)

		deltaAngle := distanceBetweenPegs / math.Sqrt(halfWidth*halfWidth*sinAngle*sinAngle+halfHeight*halfHeight*cosAngle*cosAngle)
		angle += deltaAngle
	}

	for index, p := range tc.pegs {
		peg := p.(*pegCircle)
		peg.X *= float64(tc.hiddenCanvasContext.Width()) / domainSize.Width
		peg.Y *= float64(tc.hiddenCanvasContext.Height()) / domainSize.Height
		peg.index = index
	}

	return tc.pegs
}

func (tc *ThreadComputer) computeSegmentPotential(peg1, peg2 Peg) float64 {
	tc.uploadCanvasDataToCPU()
	var potential float64 = 0

	segmentLength := common.Distance(common.Point{
		X: peg1.GetX(),
		Y: peg1.GetY(),
	}, common.Point{
		X: peg2.GetX(),
		Y: peg2.GetY(),
	})

	samplesNumber := math.Ceil(segmentLength)
	for iSample := 0; iSample < int(samplesNumber); iSample++ {
		r := (float64(iSample) + 1) / (float64(samplesNumber) + 1)
		sample := common.Point{
			X: common.Mix(peg1.GetX(), peg2.GetX(), r),
			Y: common.Mix(peg1.GetY(), peg2.GetY(), r),
		}

		imageValue := tc.sampleCanvasData(sample)
		finalValue := imageValue + tc.lineOpacityInternal*255
		contribution := 127 - finalValue

		potential += contribution
	}

	return potential / samplesNumber
}

func (tc *ThreadComputer) sampleCanvasData(coords common.Point) float64 {
	width := tc.hiddenCanvasData.Bounds().Size().X
	height := tc.hiddenCanvasData.Bounds().Size().Y

	minX := clamp(int(math.Floor(coords.X)), 0, width-1)
	maxX := clamp(int(math.Ceil(coords.X)), 0, width-1)
	minY := clamp(int(math.Floor(coords.Y)), 0, height-1)
	maxY := clamp(int(math.Ceil(coords.Y)), 0, height-1)

	topLeft := tc.sampleCanvasPixel(minX, minY)
	topRight := tc.sampleCanvasPixel(maxX, minY)
	bottomLeft := tc.sampleCanvasPixel(minX, maxY)
	bottomRight := tc.sampleCanvasPixel(maxX, maxY)

	fractX := math.Mod(coords.X, 1)
	top := common.Mix(float64(topLeft), float64(topRight), fractX)
	bottom := common.Mix(float64(bottomLeft), float64(bottomRight), fractX)

	fractY := math.Mod(coords.Y, 1)
	return common.Mix(top, bottom, fractY)
}

func (tc *ThreadComputer) sampleCanvasPixel(pixelX, pixelY int) uint8 {
	index := 4 * (pixelX + pixelY*tc.hiddenCanvasData.Rect.Dx())
	return tc.thread.SampleCanvas(tc.hiddenCanvasData.Pix, index)
}

func randomSegment(candidates []Segment) Segment {
	if len(candidates) == 0 {
		return Segment{}
	}

	randomIndex := fastrand.Uint32n(uint32(len(candidates)))
	return candidates[int(randomIndex)]
}

func randomPeg(candidates []Peg) Peg {
	if len(candidates) == 0 {
		return nil
	}

	randomIndex := fastrand.Uint32n(uint32(len(candidates)))
	return candidates[int(randomIndex)]
}

// 通过Quality，确定canvas的size
func computeBestSize(sourceImage image.Image, maxSize int) common.Size {
	maxSourceSide := sourceImage.Bounds().Size().X
	if sourceImage.Bounds().Size().Y > maxSourceSide {
		maxSourceSide = sourceImage.Bounds().Size().Y
	}
	sizingFactor := float64(maxSize) / float64(maxSourceSide)
	return common.Size{
		Width:  math.Ceil(float64(sourceImage.Bounds().Size().X) * sizingFactor),
		Height: math.Ceil(float64(sourceImage.Bounds().Size().Y) * sizingFactor),
	}
}

func clamp(x, min, max int) int {
	if x < min {
		return min
	} else if x > max {
		return max
	}

	return x
}

func (tc *ThreadComputer) OutputCanvas() {
	f, err := os.OpenFile("./result.png", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, tc.canvasBackend.Image)
	if err != nil {
		panic(err)
	}
}

func (tc *ThreadComputer) DrawPegs(p plotter.Plotter) {
	transformation := tc.computeTransformation(p.Size())
	pointSize := 0.5 * (transformation.Scaling * float64(tc.hiddenCanvasScale))

	points := make([]common.Point, 0)
	for _, p := range tc.pegs {
		point := transformation.Transform(common.Point{
			X: p.GetX(),
			Y: p.GetY(),
		})
		points = append(points, point)
	}

	p.DrawPoints(points, "red", pointSize)
}
