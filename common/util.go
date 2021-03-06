package common

import "math"

func Distance(a Point, b Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func Rotate(p Point, angle float64) Point {
	cosAngle := math.Cos(angle)
	sinAngle := math.Sin(angle)
	return Point{
		X: p.X*cosAngle - p.Y*sinAngle,
		Y: p.X*sinAngle + p.Y*cosAngle,
	}
}

func Mix(a, b, x float64) float64 {
	return a*(1-x) + b*x
}

