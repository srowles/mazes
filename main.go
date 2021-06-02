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
		for {
			container.Refresh()
			time.Sleep(100 * time.Millisecond)
		}
	}()
	w.ShowAndRun()
}

type scale struct {
}

func (s *scale) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(100, 100)
}

func (s *scale) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	for _, o := range objects {
		size := o.MinSize()
		size.Height *= containerSize.Height
		size.Width *= containerSize.Width
		pos := o.Position()
		pos.X *= containerSize.Width
		pos.Y *= containerSize.Height
		o.Resize(size)
		o.Move(pos)
	}
}

func createCells(maze *grid.Grid) []fyne.CanvasObject {
	var result []fyne.CanvasObject
	incx := 1.0 / float32(maze.Width())
	incy := 1.0 / float32(maze.Height())
	var X, Y float32
	for x := 0; x < maze.Width(); x++ {
		for y := 0; y < maze.Height(); y++ {
			cw := CellWidget{
				cell: maze.CellAt(grid.Point{X: x, Y: y}),
			}
			cw.ExtendBaseWidget(&cw)
			cw.Move(fyne.NewPos(X, Y))
			cw.Resize(fyne.NewSize(incx, incy))
			result = append(result, &cw)
			Y += incy
		}
		X += incx
		Y = 0
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
		north: line(0, 0, 0.1, 0),
		south: line(0, 0.1, 0.1, 0),
		east:  line(0.1, 0, 0, 0.1),
		west:  line(0, 0, 0, 0.1),
	}
}

// MinSize returns the size that this widget should not shrink below
func (c *CellWidget) MinSize() fyne.Size {
	c.ExtendBaseWidget(c)
	return fyne.NewSize(0.1, 0.1)
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
	return fyne.NewSize(0.1, 0.1)
}

func (c *CellWidgetRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{c.north, c.south, c.east, c.west}
}

func (c *CellWidgetRenderer) Destroy() {

}

func (c *CellWidgetRenderer) Refresh() {
	c.north.Hidden = c.cell.ExitNorth
	c.north.Refresh()
	c.south.Hidden = c.cell.ExitSouth
	c.south.Refresh()
	c.east.Hidden = c.cell.ExitEast
	c.east.Refresh()
	c.west.Hidden = c.cell.ExitWest
	c.west.Refresh()
}

var red = color.RGBA{R: 255, A: 255}

func line(x, y, w, h float32) *canvas.Line {
	sx := x
	sy := y
	ex := x + w
	ey := y + h
	line := canvas.NewLine(red)
	line.Show()
	line.Position1 = fyne.NewPos(sx, sy)
	line.Position2 = fyne.NewPos(ex, ey)
	return line
}
