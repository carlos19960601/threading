package common

import "math"

func Distance(a Point, b Point) float64 {
	dx := a.x - b.x
	dy := a.y - b.y
	return math.Sqrt(dx*dx + dy*dy)
}

func Rotate(p Point, angle float64) Point {
	cosAngle := math.Cos(angle)
	sinAngle := math.Sin(angle)
	return Point{
		x: p.x*cosAngle - p.y*sinAngle,
		y: p.x*sinAngle + p.y*cosAngle,
	}
}

func Mix(a, b, x float64) float64 {
	return a*(1-x) + b*x
}
