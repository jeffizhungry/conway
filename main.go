package main

import (
	"bufio"
	"fmt"
	"io"
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

// Conway encapsulates logic for Conway's Game of Life
type Conway struct {
	Living map[Point]bool
}

// NewConway creates a new Conway object.
func NewConway() *Conway {
	return &Conway{Living: map[Point]bool{}}
}

// Parse parses contents with Life 1.06 format
func (c *Conway) Parse(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
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
	sort.Slice(lines, func(i, j int) bool {
		return lines[i] < lines[j]
	})
	_, err := output.Write([]byte(strings.Join(lines, "\n")))
	return err
}

func main() {
	c := NewConway()
	if err := c.Parse(os.Stdin); err != nil {
		logrus.WithError(err).Fatal("Failed to read input")
	}
	if err := c.PrintLife106Format(os.Stdout); err != nil {
		logrus.WithError(err).Fatal("Failed to print output")
	}
}
