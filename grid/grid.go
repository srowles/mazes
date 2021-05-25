package grid

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// Cell represents a single cell in the grid
type Cell struct {
	North, South, East, West                 *Cell
	ExitNorth, ExitSouth, ExitEast, ExitWest bool
}

// Point is an x/y coordinate
type Point struct {
	X int
	Y int
}

// Grid represents a selection of cells making up a maze
type Grid struct {
	sync.Mutex
	width, height int
	cells         map[Point]*Cell
}

// New creates and returns a pointer to a new grid
func New(width, height int) *Grid {
	grid := &Grid{
		width:  width,
		height: height,
		cells:  make(map[Point]*Cell, width*height),
	}
	grid.Empty()

	return grid
}

func (g *Grid) BinaryTree() {
	// Create an empty grid
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			time.Sleep(2 * time.Second)
			g.Lock()
			// adjust 1 cell
			g.Unlock()
		}
	}
}

func (g *Grid) Empty() {
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			g.cells[Point{X: x, Y: y}] = &Cell{}
		}
	}
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			if x-1 >= 0 {
				g.cells[Point{X: x, Y: y}].West = g.cells[Point{X: x - 1, Y: y}]
			}
			if x+1 <= g.width {
				g.cells[Point{X: x, Y: y}].East = g.cells[Point{X: x + 1, Y: y}]
			}
			if y-1 >= 0 {
				g.cells[Point{X: x, Y: y}].North = g.cells[Point{X: x, Y: y - 1}]
			}
			if y+1 <= g.height {
				g.cells[Point{X: x, Y: y}].South = g.cells[Point{X: x, Y: y + 1}]
			}
		}
	}
}

func (g *Grid) Width() int {
	return g.width
}

func (g *Grid) Height() int {
	return g.height
}

func (g *Grid) CellAt(p Point) *Cell {
	g.Lock()
	defer g.Unlock()
	return g.cells[p]
}

var red = color.RGBA{R: 255, A: 255}

func (g *Grid) Lines() []fyne.CanvasObject {
	var result []fyne.CanvasObject
	// use single "pixel" lines for the corners
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			cell := g.CellAt(Point{X: x, Y: y})
			sx := x * 10
			sy := y * 10

			if !cell.ExitNorth {
				result = append(result, line(sx, sy, 9, 0))
			}
			if !cell.ExitWest {
				result = append(result, line(sx, sy, 0, 9))
			}
			if !cell.ExitEast {
				result = append(result, line(sx+9, sy, 0, 9))
			}
			if !cell.ExitSouth {
				result = append(result, line(sx, sy+9, 9, 0))
			}
		}
	}
	return result
}

func line(x, y, w, h int) *canvas.Line {
	sx := float32(x)
	sy := float32(y)
	ex := float32(x + w)
	ey := float32(y + h)
	line := canvas.NewLine(red)
	line.Show()
	line.Position1 = fyne.NewPos(sx, sy)
	line.Position2 = fyne.NewPos(ex, ey)
	return line
}
