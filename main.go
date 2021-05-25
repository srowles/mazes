package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/srowles/mazes/grid"
)

func main() {
	a := app.New()
	w := a.NewWindow("Mazes")

	maze := grid.New(10, 10)
	go maze.BinaryTree()
	container := container.New(&scale{}, maze.Lines()...)
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
