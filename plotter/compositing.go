package plotter

type EColor int

const (
	MONOCHROME EColor = iota
	RED
	GREEN
	BLUE
)

type Color struct {
	r int
	g int
	b int
}
