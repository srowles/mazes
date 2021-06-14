package grid

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	algoWait = 20 * time.Millisecond
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
		c.grid.flagRefreshRequired(c.Pos)
		return
	}
	if c.grid.finish == nil {
		c.grid.finish = c
		c.Text = "F"
		c.grid.flagRefreshRequired(c.Pos)
		go c.grid.Route()
		return
	}
	// reset if clicked 3rd time
	for _, c := range c.grid.cells {
		c.Text = ""
		c.grid.flagRefreshRequired(c.Pos)
	}
	c.grid.finish = nil
	c.grid.start = c
	c.grid.start.Text = "S"
	c.grid.flagRefreshRequired(c.grid.start.Pos)
}

func (c *Cell) accessible() []*Cell {
	var cells []*Cell
	if c.ExitNorth {
		cells = append(cells, c.North)
	}
	if c.ExitSouth {
		cells = append(cells, c.South)
	}
	if c.ExitEast {
		cells = append(cells, c.East)
	}
	if c.ExitWest {
		cells = append(cells, c.West)
	}
	return cells
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
	distances := make(map[Point]int)
	var frontier []*Cell
	dist := 0
	distances[g.start.Pos] = dist
	frontier = append(frontier, g.start.accessible()...)
	for {
		dist++
		var newFrontier []*Cell
		for _, c := range frontier {
			distances[c.Pos] = dist
			if c.Text != "F" {
				c.Text = fmt.Sprintf("%d", dist)
			}
			g.flagRefreshRequired(c.Pos)
			for _, next := range c.accessible() {
				if _, visited := distances[next.Pos]; !visited {
					newFrontier = append(newFrontier, next)
				}
			}
		}
		if len(newFrontier) == 0 {
			break
		}
		frontier = newFrontier
		time.Sleep(algoWait)
	}

	time.Sleep(4 * algoWait)

	cell := g.finish
	for {
		g.flagRefreshRequired(cell.Pos)
		nextPoint := g.findNext(distances, cell.Pos)
		cell = g.CellAt(nextPoint)
		if cell == g.start {
			break
		}
		cell.Text = "█"
		time.Sleep(algoWait)
	}

	for _, c := range g.cells {
		if c.Text != "█" && c.Text != "S" && c.Text != "F" {
			c.Text = ""
			g.flagRefreshRequired(c.Pos)
		}
	}
}

func (g *Grid) findNext(distances map[Point]int, pos Point) Point {
	current := distances[pos]
	next := current - 1
	for _, c := range g.CellAt(pos).accessible() {
		if distances[c.Pos] == next {
			return c.Pos
		}
	}
	panic("uhoh")
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
		c.Text = ""
	}
}

func (g *Grid) BinaryTree() {
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			time.Sleep(algoWait)
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
			time.Sleep(algoWait)
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
