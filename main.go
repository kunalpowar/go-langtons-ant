package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

type direction int

const (
	dirUp direction = iota
	dirDown
	dirRight
	dirLeft
)

func (d direction) String() string {
	switch d {
	case dirUp:
		return "Up"
	case dirDown:
		return "Down"
	case dirRight:
		return "Right"
	case dirLeft:
		return "Left"
	}

	panic("invalid direction")
}

type gridAction int

const (
	decRow gridAction = iota
	incRow
	decCol
	incCol
)

type directionChange struct {
	curr, next direction
}

type nextAction struct {
	gridAction gridAction
	dir        direction
}

var gridActionReference = map[directionChange]nextAction{
	{dirUp, dirRight}:    {incCol, dirRight},
	{dirUp, dirLeft}:     {decCol, dirLeft},
	{dirDown, dirRight}:  {decCol, dirLeft},
	{dirDown, dirLeft}:   {incCol, dirRight},
	{dirRight, dirRight}: {incRow, dirDown},
	{dirRight, dirLeft}:  {decRow, dirUp},
	{dirLeft, dirRight}:  {decRow, dirUp},
	{dirLeft, dirLeft}:   {incRow, dirDown},
}

type ant struct {
	r, c int
	dir  direction
}

func (a *ant) String() string {
	return fmt.Sprintf("row: %d, column: %d and pointed %s", a.r, a.c, a.dir)
}

func (a *ant) move(nextDir direction, steps int) {
	log.Debugf("moving ant at %s to direction %s and %d steps", a, nextDir, steps)

	action, present := gridActionReference[directionChange{a.dir, nextDir}]
	if !present {
		log.Fatalf("no grid action reference found for current direction %s and next direction %s", a.dir, nextDir)
	}

	switch action.gridAction {
	case incRow:
		a.r += steps
	case decRow:
		a.r -= steps
	case incCol:
		a.c += steps
	case decCol:
		a.c -= steps
	default:
		log.Fatalf("got invalid grid action reference for ant %s when moving towards %s", a, dirLeft)
	}

	a.dir = action.dir
}

type grid [][]bool

func (g grid) String() string {
	var buf bytes.Buffer
	for _, r := range g {
		for _, c := range r {
			if c {
				if _, err := buf.WriteString("1 "); err != nil {
					panic(err)
				}

				continue
			}

			if _, err := buf.WriteString("0 "); err != nil {
				panic(err)
			}
		}
		if _, err := buf.WriteString("\n"); err != nil {
			panic(err)
		}
	}

	return buf.String()
}

var outOfBounds = false

func (g grid) addColumnOnLeft() grid {
	gg := make(grid, len(g))
	for r, _ := range gg {
		gg[r] = append([]bool{false}, g[r]...)
	}

	return gg
}

func (g grid) addRowOnTop() grid {
	gg := make(grid, len(g)+1)
	gg[0] = make([]bool, len(g[0]))
	for r, _ := range g {
		gg[r+1] = g[r]
	}

	return gg
}

func (g grid) addColumnOnRight() grid {
	gg := make(grid, len(g))
	for r, _ := range gg {
		gg[r] = append(g[r], false)
	}

	return gg
}

func (g grid) addRowOnBottom() grid {
	gg := make(grid, len(g)+1)
	for r, _ := range g {
		gg[r] = g[r]
	}
	gg[len(gg)-1] = make([]bool, len(gg[0]))

	return gg
}

func (g grid) iterate(a *ant) grid {
	currIteration++

	log.Debugf("grid size: rows: %d, cols: %d", len(g), len(g[0]))
	log.Debugf("running iteration %d", currIteration)
	if currIteration == *iterations || outOfBounds {
		return g
	}

	oldR, oldC := a.r, a.c

	// At a black square, turn 90° left, flip the color of the square, move forward one unit
	if g[a.r][a.c] {
		a.move(dirLeft, 1)
	} else {
		// At a white square, turn 90° right, flip the color of the square, move forward one unit
		a.move(dirRight, 1)
	}
	log.Debugf("ant new position is at %s", a)

	switch {
	case a.c < 0:
		log.Debugf("adding column on left")
		a.c = 0
		g = g.addColumnOnLeft()

	case a.r < 0:
		a.r = 0
		log.Debugf("adding row on top")
		g = g.addRowOnTop()

	case a.c >= len(g[0]):
		log.Debugf("adding column on right")
		g = g.addColumnOnRight()

	case a.r >= len(g):
		log.Debugf("adding row on bottom")
		g = g.addRowOnBottom()
	}

	log.Debugf("flipping cell at row %d col %d", oldR, oldC)
	g[oldR][oldC] = !g[oldR][oldC]

	log.Debugf("ant at %s\n", a)
	return g.iterate(a)
}

var (
	gridSize   = flag.Int("s", 10, "initial size of grid")
	iterations = flag.Int("n", 10, "number of iterations to run")
	logLevel   = flag.String("v", "info", "debug level <info|debug|fatal>")

	currIteration = -1
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	flag.Parse()

	switch *logLevel {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		log.Fatalf("invalid debug level %s. can be one of info, debug and fatal", *logLevel)
	}

	log.Infof("Grid size: %d", *gridSize)

	g := make(grid, *gridSize)
	for i, _ := range g {
		g[i] = make([]bool, *gridSize)
	}

	a := ant{dir: dirUp, r: *gridSize / 2, c: *gridSize / 2}
	log.Infof("ant starting at %s", &a)

	g = g.iterate(&a)

	fmt.Printf("End grid: \n%s", g)
}
