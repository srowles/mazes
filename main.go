package main

import (
	"fmt"
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
	maze := grid.New(10, 10)

	a := app.New()
	w := a.NewWindow("Mazes")
	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Quit", func() {
				w.Close()
			}),
		),
		fyne.NewMenu("Generate",
			fyne.NewMenuItem("BinaryTree", func() {
				maze.Reset()
				maze.BinaryTree()
			}),
		),
	)
	w.SetMainMenu(menu)

	cells := createCells(maze, 10)
	container := container.New(&scale{cellsWide: float32(maze.Width()), cellsHigh: float32(maze.Height())}, cells...)
	// container := container.NewWithoutLayout(cells...)
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
	cellsWide float32
	cellsHigh float32
}

func (s *scale) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(100, 100)
}

func (s *scale) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	xscale := containerSize.Width / (s.cellsWide)
	yscale := containerSize.Height / (s.cellsHigh)
	for _, o := range objects {
		if obj, ok := o.(*CellWidget); ok {
			pos := obj.cell.Pos
			newX := float32(pos.X) * xscale
			newY := float32(pos.Y) * yscale
			o.Resize(fyne.NewSize(xscale, yscale))
			o.Move(fyne.NewPos(newX, newY))
		} else {
			fmt.Printf("%t\n", o)
		}
	}
}

func createCells(maze *grid.Grid, size float32) []fyne.CanvasObject {
	var result []fyne.CanvasObject
	for x := 0; x < maze.Width(); x++ {
		for y := 0; y < maze.Height(); y++ {
			cw := CellWidget{
				cell: maze.CellAt(grid.Point{X: x, Y: y}),
				size: size,
			}
			cw.ExtendBaseWidget(&cw)
			cw.Move(fyne.NewPos(float32(x)*size, float32(y)*size))
			cw.Resize(fyne.NewSize(size, size))
			result = append(result, &cw)
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
		north: line(0, 0, c.size, 0),
		south: line(0, c.size, c.size, 0),
		east:  line(c.size, 0, 0, c.size),
		west:  line(0, 0, 0, c.size),
	}
}

// MinSize returns the size that this widget should not shrink below
func (c *CellWidget) MinSize() fyne.Size {
	c.ExtendBaseWidget(c)
	return fyne.NewSize(c.size, c.size)
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
