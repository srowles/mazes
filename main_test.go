package main

import (
	"testing"

	"fyne.io/fyne/v2"
	"github.com/srowles/mazes/grid"
	"github.com/stretchr/testify/assert"
)

func TestLayout1(t *testing.T) {
	maze := grid.New(1, 1)
	cells := createCells(maze, 10)
	s := scale{}
	s.Layout(cells, fyne.NewSize(1000, 1000))
	for _, c := range cells {
		size := c.Size()
		assert.Equal(t, fyne.NewSize(1000, 1000), size)
		assert.Equal(t, fyne.NewPos(0, 0), c.Position())
	}
}

func TestLayout2(t *testing.T) {
	maze := grid.New(2, 2)
	cells := createCells(maze, 10)
	s := scale{}
	s.Layout(cells, fyne.NewSize(1000, 1000))
	assert.Equal(t, fyne.NewSize(500, 500), cells[0].Size())
	assert.Equal(t, fyne.NewPos(0, 0), cells[0].Position())
	assert.Equal(t, fyne.NewSize(500, 500), cells[1].Size())
	assert.Equal(t, fyne.NewPos(0, 500), cells[1].Position())
	assert.Equal(t, fyne.NewSize(500, 500), cells[2].Size())
	assert.Equal(t, fyne.NewPos(500, 0), cells[2].Position())
	assert.Equal(t, fyne.NewSize(500, 500), cells[3].Size())
	assert.Equal(t, fyne.NewPos(500, 500), cells[3].Position())
}
