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
	return fmt.Sprintf("row: %d, column: %d and pointed: %s", a.r, a.c, a.dir)
}

func (a *ant) applyGridAction(nextDir direction, steps int) {
	log.Debugf("applying grid action to ant %s with next direction %s and steps %d", a, nextDir, steps)

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

func (a *ant) move(currCell bool, steps int) {
	// At a black square, turn 90° left, flip the color of the square, move forward one unit
	if currCell {
		a.applyGridAction(dirLeft, steps)
		return
	}

	// At a white square, turn 90° right, flip the color of the square, move forward one unit
	a.applyGridAction(dirRight, steps)
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

func (g grid) iterate(a *ant) {
	currIteration++

	log.Debugf("running iteration %d", currIteration)
	if currIteration == *iterations || outOfBounds {
		return
	}

	oldR, oldC := a.r, a.c

	a.move(g[a.r][a.c], 1)
	if a.c >= len(g[0]) || a.r >= len(g) || a.c < 0 || a.r < 0 {
		outOfBounds = true
		log.Debugln("ant going out of bounds")
		return
	}

	log.Debugf("Flipping cell at row %d col %d", oldR, oldC)
	g[oldR][oldC] = !g[oldR][oldC]

	log.Debugf("%s\n", a)
	g.iterate(a)
}

var (
	gridSize   = flag.Int("s", 10, "size of the grid")
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
	log.Infof("Ant starting at %s", a)

	g.iterate(&a)

	fmt.Printf("End grid: \n%s", g)
}
