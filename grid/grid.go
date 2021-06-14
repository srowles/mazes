package grid

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Cell represents a single cell in the grid
type Cell struct {
	Pos                                      Point
	North, South, East, West                 *Cell
	ExitNorth, ExitSouth, ExitEast, ExitWest bool
	Text                                     string
	grid                                     *Grid
}

func (c *Cell) Clicked() {
	if c.grid.start == nil {
		c.grid.start = c
		c.Text = "S"
		fmt.Println("start = ", c)
		c.grid.flagRefreshRequired(c.Pos)
		return
	}
	if c.grid.finish == nil {
		c.grid.finish = c
		c.Text = "F"
		fmt.Println("finish = ", c)
		c.grid.flagRefreshRequired(c.Pos)
		go c.grid.Route()
		return
	}
	// reset if clicked 3rd time
	// TODO clear all cells!
	c.grid.start.Text = ""
	c.grid.flagRefreshRequired(c.grid.start.Pos)
	c.grid.finish.Text = ""
	c.grid.flagRefreshRequired(c.grid.finish.Pos)
	c.grid.finish = nil
	c.grid.start = c
	c.grid.start.Text = "S"
	c.grid.flagRefreshRequired(c.grid.start.Pos)
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
	// refreshChan will have an emptry struct on it if refresh is required
	refreshChan   chan Point
	start, finish *Cell
}

// New creates and returns a pointer to a new grid
func New(width, height int) *Grid {
	grid := &Grid{
		width:       width,
		height:      height,
		cells:       make(map[Point]*Cell, width*height),
		refreshChan: make(chan Point, width*height+1),
	}
	grid.Empty()

	return grid
}

func (g *Grid) Route() {

}

func (g *Grid) RequiresRefresh() chan Point {
	return g.refreshChan
}

func (g *Grid) flagRefreshRequired(p Point) {
	g.refreshChan <- p
}

func (g *Grid) Reset() {
	for _, c := range g.cells {
		c.ExitNorth = false
		c.ExitSouth = false
		c.ExitEast = false
		c.ExitWest = false
	}
}

func (g *Grid) BinaryTree() {
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			time.Sleep(20 * time.Millisecond)
			cell := g.CellAt(Point{X: x, Y: y})
			if cell.North == nil && cell.East == nil {
				// in top right with nothing to carve, so skip
				continue
			}

			if cell.North == nil {
				// cannot carve north, must carve east
				g.carveEast(cell)
				continue
			}

			if cell.East == nil {
				// cannot carve east, must carve north
				g.carveNorth(cell)
				continue
			}

			if heads() {
				g.carveEast(cell)
			} else {
				g.carveNorth(cell)
			}
		}
	}
}

func heads() bool {
	r := rand.Intn(2)
	return r == 1
}

func (g *Grid) Sidewinder() {
	for row := 0; row < g.height; row++ {
		var cellRun []*Cell
		for x := 0; x < g.width; x++ {
			time.Sleep(20 * time.Millisecond)
			cell := g.CellAt(Point{X: x, Y: row})
			cellRun = append(cellRun, cell)

			if cell.East == nil {
				g.carveNorth(cellRun[rand.Intn(len(cellRun))])
				cellRun = nil
				continue
			}
			if cell.North == nil && cell.East == nil {
				// top right, do nothing
				continue
			}
			if cell.North == nil {
				g.carveEast(cell)
				continue
			}

			if heads() {
				g.carveNorth(cellRun[rand.Intn(len(cellRun))])
				cellRun = nil
				continue
			}
			g.carveEast(cell)
		}
	}
}

func (g *Grid) carveEast(cell *Cell) {
	if cell.East == nil {
		return
	}
	cell.ExitEast = true
	cell.East.ExitWest = true
	g.flagRefreshRequired(cell.Pos)
}

func (g *Grid) carveNorth(cell *Cell) {
	if cell.North == nil {
		return
	}
	cell.ExitNorth = true
	cell.North.ExitSouth = true
	g.flagRefreshRequired(cell.Pos)
}

func (g *Grid) Empty() {
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			pos := Point{X: x, Y: y}
			g.cells[pos] = &Cell{Pos: pos, grid: g}
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
			g.flagRefreshRequired(Point{X: x, Y: y})
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
