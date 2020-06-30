package main

import (
	"bytes"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// End to end test
func TestConawy(t *testing.T) {
	testcases := map[string]struct {
		numGenerations int
		input          string
		expected       string
	}{
		"no generations": {
			numGenerations: 0,
			input: `#Life 1.06
-1 1
0 -1
0 1
1 0
1 1`,
			expected: `#Life 1.06
-1 1
0 -1
0 1
1 0
1 1`,
		},
		"single cell, immediate death": {
			numGenerations: 1,
			input: `#Life 1.06
				1 1`,
			expected: "#Life 1.06\n",
		},
		"still life": {
			numGenerations: 20,
			input: `#Life 1.06
0 0
0 1
1 0
1 1`,
			expected: `#Life 1.06
0 0
0 1
1 0
1 1`,
		},
		"oscillator (even number of gens)": {
			numGenerations: 2000,
			input: `#Life 1.06
-1 0
0 0
1 0`,
			expected: `#Life 1.06
-1 0
0 0
1 0`,
		},
		"oscillator (odd number of gens)": {
			numGenerations: 2001,
			input: `#Life 1.06
-1 0
0 0
1 0`,
			expected: `#Life 1.06
0 -1
0 0
0 1`,
		},
	}
	for testName, tc := range testcases {
		t.Run(testName, func(t *testing.T) {
			// Create
			c := NewConway()

			// Parse
			input := bytes.NewBufferString(tc.input)
			assert.NoError(t, c.Parse(input))

			// Run simulations
			c.Simulate(tc.numGenerations)

			// Compare
			output := bytes.Buffer{}
			assert.NoError(t, c.PrintLife106Format(&output))
			assert.Equal(t, tc.expected, string(output.Bytes()))
		})
	}
}

func Test_computeValidNeighbors(t *testing.T) {
	testcases := map[string]struct {
		p        Point
		expected []Point
	}{
		"center": {
			p: Point{0, 0},
			expected: []Point{
				{X: -1, Y: -1},
				{X: -1, Y: 0},
				{X: -1, Y: 1},
				{X: 0, Y: -1},
				{X: 0, Y: 1},
				{X: 1, Y: -1},
				{X: 1, Y: 0},
				{X: 1, Y: 1},
			},
		},
		"top left corner": {
			p: Point{math.MinInt64, math.MaxInt64},
			expected: []Point{
				{X: -9223372036854775808, Y: 9223372036854775806},
				{X: -9223372036854775807, Y: 9223372036854775806},
				{X: -9223372036854775807, Y: 9223372036854775807},
			},
		},
		"top right corner": {
			p: Point{math.MaxInt64, math.MaxInt64},
			expected: []Point{
				{X: 9223372036854775806, Y: 9223372036854775806},
				{X: 9223372036854775806, Y: 9223372036854775807},
				{X: 9223372036854775807, Y: 9223372036854775806},
			},
		},
		"bottom left corner": {
			p: Point{math.MinInt64, math.MinInt64},
			expected: []Point{
				{X: -9223372036854775808, Y: -9223372036854775807},
				{X: -9223372036854775807, Y: -9223372036854775808},
				{X: -9223372036854775807, Y: -9223372036854775807},
			},
		},
		"bottom right corner": {
			p: Point{math.MaxInt64, math.MinInt64},
			expected: []Point{
				{X: 9223372036854775806, Y: -9223372036854775808},
				{X: 9223372036854775806, Y: -9223372036854775807},
				{X: 9223372036854775807, Y: -9223372036854775807},
			},
		},
	}
	for testName, tc := range testcases {
		t.Run(testName, func(t *testing.T) {
			actual := computeValidNeighbors(tc.p)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
