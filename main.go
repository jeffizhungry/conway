package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// Point represents a coordinate
type Point struct {
	X int64
	Y int64
}

// ------------------------------------
// Game of Life
// ------------------------------------

// Conway encapsulates logic for Conway's Game of Life
type Conway struct {
	// Using this data structure for quick lookups and to handle large integers
	// Assuming a reasonable amount of living datapoints
	Living map[Point]bool
}

// NewConway creates a new Conway object.
func NewConway() *Conway {
	return &Conway{Living: map[Point]bool{}}
}

// Parse parses contents with Life 1.06 format and stores
// coordinate data in Conway.Living
//
// TODO(jehwang): Ask about how we want to handle invalid data,
// fail hard or handle gracefully?
func (c *Conway) Parse(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// TODO(jehwang): Should consider a file invalid if it doesn't have this header?
		if line == "#Life 1.06" {
			continue
		}
		spl := strings.Split(line, " ")
		if len(spl) != 2 {
			// Ignore invalid lines
			continue
		}
		x, err := strconv.ParseInt(spl[0], 10, 64)
		if err != nil {
			continue
		}
		y, err := strconv.ParseInt(spl[1], 10, 64)
		if err != nil {
			continue
		}
		c.Living[Point{X: x, Y: y}] = true
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// PrintLife106Format parses contents with Life 1.06 format
func (c *Conway) PrintLife106Format(output io.Writer) error {
	if _, err := output.Write([]byte("#Life 1.06\n")); err != nil {
		return err
	}

	var lines []string
	for p := range c.Living {
		// TODO(jehwang): Think about compat with 32-bit machines
		x := strconv.Itoa(int(p.X))
		y := strconv.Itoa(int(p.Y))
		lines = append(lines, fmt.Sprintf("%v %v", x, y))
	}
	if len(lines) == 0 {
		return nil
	}
	// Sort for stable output
	sort.Slice(lines, func(i, j int) bool {
		return lines[i] < lines[j]
	})
	_, err := output.Write([]byte(strings.Join(lines, "\n")))
	return err
}

// Simulate runs simulations for conway's game of life
func (c *Conway) Simulate(numGenerations int) {
	for i := 0; i < numGenerations; i++ {
		c.simulateOneGeneration()
	}
}

// simulateOneGeneration runs simulations for conway's game of life
func (c *Conway) simulateOneGeneration() {
	nextGeneration := map[Point]bool{}
	for k, v := range c.Living {
		nextGeneration[k] = v
	}

	// NOTE(jehwang): Could parallelize by chopping up living dataset
	// and parallelize w.r.t # of cores. Would need to introduce lock
	// on c.Living and collect data on deaths / births

	// --------------------------------
	// Compute deaths
	// --------------------------------
	for p := range c.Living {
		n := c.findNumberOfLivingNeighbors(p)
		if shouldDie(n) {
			delete(nextGeneration, p)
		}
	}

	// --------------------------------
	// Compute births
	// --------------------------------

	// Get list of empty cells that could possibly be reborn next round
	var possibleBabies []Point
	for p := range c.Living {
		neighbors := computeValidNeighbors(p)
		for _, p := range neighbors {
			if !c.Living[p] {
				possibleBabies = append(possibleBabies, p)
			}
		}
	}
	for _, p := range possibleBabies {
		n := c.findNumberOfLivingNeighbors(p)
		if shouldBeReborn(n) {
			nextGeneration[p] = true
		}
	}
	c.Living = nextGeneration
}

func (c *Conway) isLiving(p Point) bool {
	return c.Living[p]
}

// findNumberOfLivingNeighbors computes number of living neighbors
func (c *Conway) findNumberOfLivingNeighbors(p Point) int {
	// Old Implementation
	// neighbors := computeValidNeighbors(p)
	// var result int
	// for _, p := range neighbors {
	// 	if c.Living[p] {
	// 		result++
	// 	}
	// }

	x, y := p.X, p.Y
	var result int
	if x != math.MinInt64 && y != math.MinInt64 && c.isLiving(Point{x - 1, y - 1}) {
		result++
	}
	if x != math.MinInt64 && c.isLiving(Point{x - 1, y}) {
		result++
	}
	if x != math.MinInt64 && y != math.MaxInt64 && c.isLiving(Point{x - 1, y + 1}) {
		result++
	}
	if y != math.MinInt64 && c.isLiving(Point{x, y - 1}) {
		result++
	}
	if y != math.MaxInt64 && c.isLiving(Point{x, y + 1}) {
		result++
	}
	if x != math.MaxInt64 && y != math.MinInt64 && c.isLiving(Point{x + 1, y - 1}) {
		result++
	}
	if x != math.MaxInt64 && c.isLiving(Point{x + 1, y}) {
		result++
	}
	if x != math.MaxInt64 && y != math.MaxInt64 && c.isLiving(Point{x + 1, y + 1}) {
		result++
	}
	return result
}

// ------------------------------------
// Helper Functions
// ------------------------------------

// computeValidNeighbors returns a slice of all neighbors that are within a int64 grid
func computeValidNeighbors(p Point) []Point {
	x, y := p.X, p.Y
	var neighbors []Point
	if x != math.MinInt64 && y != math.MinInt64 {
		neighbors = append(neighbors, Point{x - 1, y - 1})
	}
	if x != math.MinInt64 {
		neighbors = append(neighbors, Point{x - 1, y})
	}
	if x != math.MinInt64 && y != math.MaxInt64 {
		neighbors = append(neighbors, Point{x - 1, y + 1})
	}
	if y != math.MinInt64 {
		neighbors = append(neighbors, Point{x, y - 1})
	}
	if y != math.MaxInt64 {
		neighbors = append(neighbors, Point{x, y + 1})
	}
	if x != math.MaxInt64 && y != math.MinInt64 {
		neighbors = append(neighbors, Point{x + 1, y - 1})
	}
	if x != math.MaxInt64 {
		neighbors = append(neighbors, Point{x + 1, y})
	}
	if x != math.MaxInt64 && y != math.MaxInt64 {
		neighbors = append(neighbors, Point{x + 1, y + 1})
	}
	return neighbors
}

// shouldDie encasulates rule about whether an "alive" cell should die
func shouldDie(numLivingNeighbors int) bool {
	return numLivingNeighbors < 2 || numLivingNeighbors > 3
}

// shouldBeReborn encasulates rule about whether a "dead" cell should be reborn
func shouldBeReborn(numLivingNeighbors int) bool {
	return numLivingNeighbors == 3
}

func main() {
	// Init
	c := NewConway()

	// Read data
	// NOTE(jehwang): This will block if nothing is provided in stdin
	if err := c.Parse(os.Stdin); err != nil {
		logrus.WithError(err).Fatal("Failed to read input")
	}

	// Run simulations
	c.Simulate(10)

	// Print out output
	if err := c.PrintLife106Format(os.Stdout); err != nil {
		logrus.WithError(err).Fatal("Failed to print output")
	}
}
