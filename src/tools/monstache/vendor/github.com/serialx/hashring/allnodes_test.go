package hashring

import (
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateWeights(n int) map[string]int {
	result := make(map[string]int)
	for i := 0; i < n; i++ {
		result[fmt.Sprintf("%03d", i)] = i + 1
	}
	return result
}

func generateNodes(n int) []string {
	result := make([]string, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, fmt.Sprintf("%03d", i))
	}
	return result
}

func TestListOf1000Nodes(t *testing.T) {
	testData := map[string]struct {
		ring *HashRing
	}{
		"nodes":   {ring: New(generateNodes(1000))},
		"weights": {ring: NewWithWeights(generateWeights(1000))},
	}

	for testName, data := range testData {
		ring := data.ring
		t.Run(testName, func(t *testing.T) {
			nodes, ok := ring.GetNodes("key", ring.Size())
			assert.True(t, ok)
			if !assert.Equal(t, ring.Size(), len(nodes)) {
				// print debug info on failure
				sort.Strings(nodes)
				fmt.Printf("%v\n", nodes)
				return
			}

			// assert that each node shows up exatly once
			sort.Strings(nodes)
			for i, node := range nodes {
				actual, err := strconv.ParseInt(node, 10, 64)
				if !assert.NoError(t, err) {
					return
				}
				if !assert.Equal(t, int64(i), actual) {
					return
				}
			}
		})
	}
}
