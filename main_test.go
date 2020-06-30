package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConawy(t *testing.T) {
	testcases := map[string]struct {
		input      string
		iterations int
		expected   string
	}{
		"basic": {
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
	}
	for testName, tc := range testcases {
		t.Run(testName, func(t *testing.T) {
			// Create
			c := NewConway()

			// Parse
			input := bytes.NewBufferString(tc.input)
			assert.NoError(t, c.Parse(input))

			// Run simulations

			// Compare
			output := bytes.Buffer{}
			assert.NoError(t, c.PrintLife106Format(&output))
			assert.Equal(t, tc.expected, string(output.Bytes()))
		})
	}
}
