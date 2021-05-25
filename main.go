package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/srowles/mazes/grid"
)

func main() {
	a := app.New()
	w := a.NewWindow("Mazes")

	maze := grid.New(10, 10)
	// go maze.BinaryTree()
	cells := createCells(maze)
	container := container.New(&scale{}, cells...)
	container.Resize(fyne.NewSize(800, 600))
	w.SetContent(container)
	w.Resize(fyne.NewSize(800, 600))
	go func() {
		container.Refresh()
		time.Sleep(time.Second)
	}()
	w.ShowAndRun()
}

type scale struct {
}

func (s *scale) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		pos := o.Position()
		size := o.MinSize()
		x := pos.X + size.Width
		if x > w {
			w = x
		}
		y := pos.Y + size.Height
		if y > h {
			h = y
		}
	}
	return fyne.NewSize(w, h)
}

func (s *scale) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	compSize := s.MinSize(objects)
	scaleX := containerSize.Width / compSize.Width
	scaleY := containerSize.Height / compSize.Height
	for _, o := range objects {
		size := o.Size()
		size.Height *= scaleY
		size.Width *= scaleX
		pos := o.Position()
		pos.X *= scaleX
		pos.Y *= scaleY
		o.Resize(size)
		o.Move(pos)
	}
}

func createCells(maze *grid.Grid) []fyne.CanvasObject {
	var result []fyne.CanvasObject
	for x := 0; x < maze.Width(); x++ {
		for y := 0; y < maze.Height(); y++ {
			cw := CellWidget{
				cell: maze.CellAt(grid.Point{X: x, Y: y}),
			}
			cw.Move(fyne.NewPos(float32(x*10), float32(y*10)))
			cw.Resize(fyne.NewSize(10, 10))
			result = append(result, &cw)
		}
	}

	return result
}

type CellWidget struct {
	widget.BaseWidget
	cell *grid.Cell
}

func (c *CellWidget) CreateRenderer() fyne.WidgetRenderer {
	return &CellWidgetRenderer{
		cell:  c.cell,
		north: line(0, 0, 10, 0),
		south: line(0, 10, 10, 0),
		east:  line(10, 0, 0, 10),
		west:  line(0, 0, 0, 10),
	}
}

// MinSize returns the size that this widget should not shrink below
func (c *CellWidget) MinSize() fyne.Size {
	c.ExtendBaseWidget(c)
	return fyne.NewSize(10, 10)
}

type CellWidgetRenderer struct {
	cell  *grid.Cell
	north *canvas.Line
	south *canvas.Line
	east  *canvas.Line
	west  *canvas.Line
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
	c.north.Hidden = c.cell.ExitNorth
	c.south.Hidden = c.cell.ExitSouth
	c.east.Hidden = c.cell.ExitEast
	c.west.Hidden = c.cell.ExitWest
}

var red = color.RGBA{R: 255, A: 255}

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
