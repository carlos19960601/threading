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

func computeRawColor(color EColor) Color {
	if color == MONOCHROME {
		return Color{r: 1, g: 1, b: 1}
	}

	r, g, b := 0, 0, 0
	if color == RED {
		r = 1
	}
	if color == GREEN {
		g = 1
	}
	if color == BLUE {
		b = 1
	}
	return Color{
		r: r,
		g: g,
		b: b,
	}
}
