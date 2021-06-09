package main

import (
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/srowles/mazes/grid"
)

func main() {
	rand.Seed(time.Now().UTC().Unix())
	maze := grid.New(20, 20)

	a := app.New()
	w := a.NewWindow("Mazes")
	var cont *fyne.Container
	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Refresh", func() {
				cont.Refresh()
			}),
			fyne.NewMenuItem("Quit", func() {
				w.Close()
			}),
		),
		fyne.NewMenu("Generate",
			fyne.NewMenuItem("Reset", func() {
				maze.Reset()
				cont.Refresh()
			}),
			fyne.NewMenuItem("BinaryTree", func() {
				maze.Reset()
				cont.Refresh()
				maze.BinaryTree()
				cont.Refresh()
			}),
			fyne.NewMenuItem("Sidewinder", func() {
				maze.Reset()
				cont.Refresh()
				maze.Sidewinder()
				cont.Refresh()
			}),
		),
	)
	w.SetMainMenu(menu)

	cellSize := float32(10)
	cellMap := createCells(maze, cellSize)
	var cells []fyne.CanvasObject
	for _, cell := range cellMap {
		cells = append(cells, cell)
	}
	cont = container.New(&scale{cellsWide: float32(maze.Width()), cellsHigh: float32(maze.Height()), size: cellSize}, cells...)
	// cont = container.NewWithoutLayout(cells...)
	cont.Resize(fyne.NewSize(800, 600))
	w.SetContent(cont)
	w.Resize(fyne.NewSize(800, 600))
	go func() {
		for {
			p := <-maze.RequiresRefresh()
			// refresh the cell that changed, and adjacent nsew
			cellMap[p].Refresh()
			if c := cellMap[grid.Point{X: p.X + 1, Y: p.Y}]; c != nil {
				c.Refresh()
			}
			if c := cellMap[grid.Point{X: p.X - 1, Y: p.Y}]; c != nil {
				c.Refresh()
			}
			if c := cellMap[grid.Point{X: p.X, Y: p.Y + 1}]; c != nil {
				c.Refresh()
			}
			if c := cellMap[grid.Point{X: p.X, Y: p.Y - 1}]; c != nil {
				c.Refresh()
			}
		}
	}()
	w.ShowAndRun()
}

type scale struct {
	size       float32
	cellsWide  float32
	cellsHigh  float32
	lastWidth  float32
	lastHeight float32
}

func (s *scale) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(s.size*s.cellsWide, s.size*s.cellsHigh)
}

func (s *scale) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if s.lastHeight == containerSize.Height && s.lastWidth == containerSize.Width {
		return
	}
	s.lastHeight = containerSize.Height
	s.lastWidth = containerSize.Width
	xscale := containerSize.Width / s.cellsWide
	yscale := containerSize.Height / s.cellsHigh
	for _, o := range objects {
		if obj, ok := o.(*CellWidget); ok {
			pos := obj.cell.Pos
			newX := float32(pos.X) * xscale
			newY := float32(pos.Y) * yscale
			o.Resize(fyne.NewSize(xscale, yscale))
			o.Move(fyne.NewPos(newX, newY))
		}
	}
}

func createCells(maze *grid.Grid, size float32) map[grid.Point]fyne.CanvasObject {
	result := make(map[grid.Point]fyne.CanvasObject, maze.Width()*maze.Height())
	for x := 0; x < maze.Width(); x++ {
		for y := 0; y < maze.Height(); y++ {
			p := grid.Point{X: x, Y: y}
			cw := CellWidget{
				cell: maze.CellAt(p),
				size: size,
			}
			cw.ExtendBaseWidget(&cw)
			cw.Move(fyne.NewPos(float32(x)*size, float32(y)*size))
			cw.Resize(fyne.NewSize(size, size))
			result[p] = &cw
		}
	}

	return result
}

type CellWidget struct {
	widget.BaseWidget
	cell *grid.Cell
	size float32
}

func (c *CellWidget) CreateRenderer() fyne.WidgetRenderer {
	return &CellWidgetRenderer{
		cell:  c.cell,
		north: line(0, 0, c.size, 0, blue),
		south: line(0, c.size, c.size, 0, blue),
		east:  line(c.size, 0, 0, c.size, blue),
		west:  line(0, 0, 0, c.size, blue),
		text:  canvas.NewText("", grayBlue),
	}
}

// MinSize returns the size that this widget should not shrink below
func (c *CellWidget) MinSize() fyne.Size {
	c.ExtendBaseWidget(c)
	return fyne.NewSize(c.size, c.size)
}

type CellWidgetRenderer struct {
	cell      *grid.Cell
	north     *canvas.Line
	south     *canvas.Line
	east      *canvas.Line
	west      *canvas.Line
	text      *canvas.Text
	textValue string
}

func (c *CellWidgetRenderer) Layout(containerSize fyne.Size) {
	c.north.Position1 = fyne.NewPos(0, 0)
	c.north.Position2 = fyne.NewPos(containerSize.Width, 0)
	c.south.Position1 = fyne.NewPos(0, containerSize.Height)
	c.south.Position2 = fyne.NewPos(containerSize.Width, containerSize.Height)
	c.east.Position1 = fyne.NewPos(containerSize.Width, 0)
	c.east.Position2 = fyne.NewPos(containerSize.Width, containerSize.Height)
	c.west.Position1 = fyne.NewPos(0, 0)
	c.west.Position2 = fyne.NewPos(0, containerSize.Height)
	c.text.Move(fyne.NewPos(0, 0))
}

func (c *CellWidgetRenderer) MinSize() fyne.Size {
	return fyne.NewSize(10, 10)
}

func (c *CellWidgetRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{c.north, c.south, c.east, c.west}
}

func (c *CellWidgetRenderer) Destroy() {

}

func (c *CellWidgetRenderer) Refresh() {
	if c.north.Hidden != c.cell.ExitNorth {
		c.north.Hidden = c.cell.ExitNorth
		c.north.Refresh()
	}
	if c.south.Hidden != c.cell.ExitSouth {
		c.south.Hidden = c.cell.ExitSouth
		c.south.Refresh()
	}
	if c.east.Hidden != c.cell.ExitEast {
		c.east.Hidden = c.cell.ExitEast
		c.east.Refresh()
	}
	if c.west.Hidden != c.cell.ExitWest {
		c.west.Hidden = c.cell.ExitWest
		c.west.Refresh()
	}
	if c.text.Text != c.textValue {
		c.text.Text = c.cell.Text
		c.text.Refresh()
	}
}

var blue = color.RGBA{R: 0, G: 64, B: 254, A: 255}
var grayBlue = color.RGBA{R: 169, G: 180, B: 212, A: 128}

func line(x, y, w, h float32, colour color.RGBA) *canvas.Line {
	sx := x
	sy := y
	ex := x + w
	ey := y + h
	line := canvas.NewLine(colour)
	line.Show()
	line.Position1 = fyne.NewPos(sx, sy)
	line.Position2 = fyne.NewPos(ex, ey)
	return line
}
