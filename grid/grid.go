package grid

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// Cell represents a single cell in the grid
type Cell struct {
	ExitNorth, ExitSouth, ExitEast, ExitWest bool
}

// Point is an x/y coordinate
type Point struct {
	X int
	Y int
}

// Grid represents a selection of cells making up a maze
type Grid struct {
	width, height int
	cells         map[Point]*Cell
}

// New creates and returns a pointer to a new grid
func New(width, height int) *Grid {
	return &Grid{
		width:  width,
		height: height,
		cells:  make(map[Point]*Cell, width*height),
	}
}

func (g *Grid) BinaryTree() {

}

func (g *Grid) Full() {
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			g.cells[Point{X: x, Y: y}] = &Cell{
				ExitNorth: true,
				ExitSouth: true,
				ExitEast:  true,
				ExitWest:  true,
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
	return g.cells[p]
}

func (g *Grid) GenerateImage(w, h int) {
	img := image.NewRGBA(image.Rect(0, 0, g.width*10, g.height*10))
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			cell := g.CellAt(Point{X: x, Y: y})
			drawCell(img, x, y, cell)
		}
	}

	f, err := os.Create("outimage.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Encode to `PNG` with `DefaultCompression` level
	// then save to file
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}

var red = color.RGBA{R: 255, A: 255}

func drawCell(img *image.RGBA, x, y int, cell *Cell) {
	x *= 10
	y *= 10
	// draw corners
	img.Set(x, 0, red)
	for a := x; a < x+10; a++ {
		img.Set(a, y, red)
		img.Set(a, y+9, red)
	}
	for a := y; a < y+9; a++ {
		img.Set(x, a, red)
		img.Set(x+9, a, red)
	}
}

func (g *Grid) Lines() []fyne.CanvasObject {
	var result []fyne.CanvasObject
	// use single "pixel" lines for the corners
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			cell := g.CellAt(Point{X: x, Y: y})
			sx := x * 10
			sy := y * 10

			if cell.ExitNorth {
				result = append(result, line(sx, sy, 9, 0))
			}
			if cell.ExitWest {
				result = append(result, line(sx, sy, 0, 9))
			}
			if cell.ExitEast {
				result = append(result, line(sx+9, sy, 0, 9))
			}
			if cell.ExitSouth {
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
