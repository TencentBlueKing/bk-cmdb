package hashring

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func expectWeights(t *testing.T, ring *HashRing, expectedWeights map[string]int) {
	weightsEquality := reflect.DeepEqual(ring.weights, expectedWeights)
	if !weightsEquality {
		t.Error("Weights expected", expectedWeights, "but got", ring.weights)
	}
}

type testPair struct {
	key  string
	node string
}

type testNodes struct {
	key   string
	nodes []string
}

func assert2Nodes(t *testing.T, prefix string, ring *HashRing, data []testNodes) {
	t.Run(prefix, func(t *testing.T) {
		allActual := make([]string, 0)
		allExpected := make([]string, 0)
		for _, pair := range data {
			nodes, ok := ring.GetNodes(pair.key, 2)
			if assert.True(t, ok) {
				allActual = append(allActual, fmt.Sprintf("%s - %v", pair.key, nodes))
				allExpected = append(allExpected, fmt.Sprintf("%s - %v", pair.key, pair.nodes))
			}
		}
		assert.Equal(t, allExpected, allActual)
	})
}

func assertNodes(t *testing.T, prefix string, ring *HashRing, allExpected []testPair) {
	t.Run(prefix, func(t *testing.T) {
		allActual := make([]testPair, 0)
		for _, pair := range allExpected {
			node, ok := ring.GetNode(pair.key)
			if assert.True(t, ok) {
				allActual = append(allActual, testPair{key: pair.key, node: node})
			}
		}
	})
}

func expectNodesABC(t *testing.T, prefix string, ring *HashRing) {

	assertNodes(t, prefix, ring, []testPair{
		{"test", "a"},
		{"test", "a"},
		{"test1", "b"},
		{"test2", "b"},
		{"test3", "c"},
		{"test4", "a"},
		{"test5", "c"},
		{"aaaa", "c"},
		{"bbbb", "a"},
	})
}

func expectNodeRangesABC(t *testing.T, prefix string, ring *HashRing) {
	assert2Nodes(t, prefix, ring, []testNodes{
		{"test", []string{"a", "c"}},
		{"test", []string{"a", "c"}},
		{"test1", []string{"b", "a"}},
		{"test2", []string{"b", "a"}},
		{"test3", []string{"c", "b"}},
		{"test4", []string{"a", "c"}},
		{"test5", []string{"c", "b"}},
		{"aaaa", []string{"c", "b"}},
		{"bbbb", []string{"a", "c"}},
	})
}

func expectNodesABCD(t *testing.T, prefix string, ring *HashRing) {
	assertNodes(t, prefix, ring, []testPair{
		{"test", "d"},
		{"test", "d"},
		{"test1", "b"},
		{"test2", "b"},
		{"test3", "c"},
		{"test4", "d"},
		{"test5", "c"},
		{"aaaa", "c"},
		{"bbbb", "d"},
	})
}

func TestNew(t *testing.T) {
	nodes := []string{"a", "b", "c"}
	ring := New(nodes)

	expectNodesABC(t, "TestNew_1_", ring)
	expectNodeRangesABC(t, "", ring)
}

func TestNewEmpty(t *testing.T) {
	nodes := []string{}
	ring := New(nodes)

	node, ok := ring.GetNode("test")
	if ok || node != "" {
		t.Error("GetNode(test) expected (\"\", false) but got (", node, ",", ok, ")")
	}

	nodes, rok := ring.GetNodes("test", 2)
	if rok || !(len(nodes) == 0) {
		t.Error("GetNode(test) expected ( [], false ) but got (", nodes, ",", rok, ")")
	}
}

func TestForMoreNodes(t *testing.T) {
	nodes := []string{"a", "b", "c"}
	ring := New(nodes)

	nodes, ok := ring.GetNodes("test", 5)
	if ok || !(len(nodes) == 0) {
		t.Error("GetNode(test) expected ( [], false ) but got (", nodes, ",", ok, ")")
	}
}

func TestForEqualNodes(t *testing.T) {
	nodes := []string{"a", "b", "c"}
	ring := New(nodes)

	nodes, ok := ring.GetNodes("test", 3)
	if !ok && (len(nodes) == 3) {
		t.Error("GetNode(test) expected ( [a b c], true ) but got (", nodes, ",", ok, ")")
	}
}

func TestNewSingle(t *testing.T) {
	nodes := []string{"a"}
	ring := New(nodes)

	assertNodes(t, "", ring, []testPair{
		{"test", "a"},
		{"test", "a"},
		{"test1", "a"},
		{"test2", "a"},
		{"test3", "a"},

		{"test14", "a"},

		{"test15", "a"},
		{"test16", "a"},
		{"test17", "a"},
		{"test18", "a"},
		{"test19", "a"},
		{"test20", "a"},
	})
}

func TestNewWeighted(t *testing.T) {
	weights := make(map[string]int)
	weights["a"] = 1
	weights["b"] = 2
	weights["c"] = 1
	ring := NewWithWeights(weights)

	assertNodes(t, "", ring, []testPair{
		{"test", "b"},
		{"test", "b"},
		{"test1", "b"},
		{"test2", "b"},
		{"test3", "c"},
		{"test4", "b"},
		{"test5", "c"},
		{"aaaa", "c"},
		{"bbbb", "b"},
	})
	assert2Nodes(t, "", ring, []testNodes{
		{"test", []string{"b", "a"}},
	})
}

func TestRemoveNode(t *testing.T) {
	nodes := []string{"a", "b", "c"}
	ring := New(nodes)
	ring = ring.RemoveNode("b")

	assertNodes(t, "", ring, []testPair{
		{"test", "a"},
		{"test", "a"},
		{"test1", "a"},
		{"test2", "a"},
		{"test3", "c"},
		{"test4", "a"},
		{"test5", "c"},
		{"aaaa", "c"},
		{"bbbb", "a"},
	})

	assert2Nodes(t, "", ring, []testNodes{
		{"test", []string{"a", "c"}},
	})
}

func TestAddNode(t *testing.T) {
	nodes := []string{"a", "c"}
	ring := New(nodes)
	ring = ring.AddNode("b")

	expectNodesABC(t, "TestAddNode_1_", ring)

	defaultWeights := map[string]int{
		"a": 1,
		"b": 1,
		"c": 1,
	}
	expectWeights(t, ring, defaultWeights)
}

func TestAddNode2(t *testing.T) {
	nodes := []string{"a", "c"}
	ring := New(nodes)
	ring = ring.AddNode("b")
	ring = ring.AddNode("b")

	expectNodesABC(t, "TestAddNode2_", ring)
	expectNodeRangesABC(t, "", ring)
}

func TestAddNode3(t *testing.T) {
	nodes := []string{"a", "b", "c"}
	ring := New(nodes)
	ring = ring.AddNode("d")

	expectNodesABCD(t, "TestAddNode3_1_", ring)

	ring = ring.AddNode("e")

	assertNodes(t, "TestAddNode3_2_", ring, []testPair{
		{"test", "d"},
		{"test", "d"},
		{"test1", "b"},
		{"test2", "e"},
		{"test3", "c"},
		{"test4", "d"},
		{"test5", "c"},
		{"aaaa", "c"},
		{"bbbb", "d"},
	})

	assert2Nodes(t, "", ring, []testNodes{
		{"test", []string{"d", "a"}},
	})

	ring = ring.AddNode("f")

	assertNodes(t, "TestAddNode3_3_", ring, []testPair{
		{"test", "d"},
		{"test", "d"},
		{"test1", "b"},
		{"test2", "e"},
		{"test3", "c"},
		{"test4", "d"},
		{"test5", "c"},
		{"aaaa", "c"},
		{"bbbb", "d"},
	})

	assert2Nodes(t, "", ring, []testNodes{
		{"test", []string{"d", "a"}},
	})
}

func TestDuplicateNodes(t *testing.T) {
	nodes := []string{"a", "a", "a", "a", "b"}
	ring := New(nodes)

	assertNodes(t, "TestDuplicateNodes_", ring, []testPair{
		{"test", "a"},
		{"test", "a"},
		{"test1", "b"},
		{"test2", "b"},
		{"test3", "b"},
		{"test4", "a"},
		{"test5", "b"},
		{"aaaa", "b"},
		{"bbbb", "a"},
	})
}

func TestAddWeightedNode(t *testing.T) {
	nodes := []string{"a", "c"}
	ring := New(nodes)
	ring = ring.AddWeightedNode("b", 0)
	ring = ring.AddWeightedNode("b", 2)
	ring = ring.AddWeightedNode("b", 2)

	assertNodes(t, "TestAddWeightedNode_", ring, []testPair{
		{"test", "b"},
		{"test", "b"},
		{"test1", "b"},
		{"test2", "b"},
		{"test3", "c"},
		{"test4", "b"},
		{"test5", "c"},
		{"aaaa", "c"},
		{"bbbb", "b"},
	})

	assert2Nodes(t, "", ring, []testNodes{
		{"test", []string{"b", "a"}},
	})
}

func TestUpdateWeightedNode(t *testing.T) {
	nodes := []string{"a", "c"}
	ring := New(nodes)
	ring = ring.AddWeightedNode("b", 1)
	ring = ring.UpdateWeightedNode("b", 2)
	ring = ring.UpdateWeightedNode("b", 2)
	ring = ring.UpdateWeightedNode("b", 0)
	ring = ring.UpdateWeightedNode("d", 2)

	assertNodes(t, "TestUpdateWeightedNode_", ring, []testPair{
		{"test", "b"},
		{"test", "b"},
		{"test1", "b"},
		{"test2", "b"},
		{"test3", "c"},
		{"test4", "b"},
		{"test5", "c"},
		{"aaaa", "c"},
		{"bbbb", "b"},
	})

	assert2Nodes(t, "", ring, []testNodes{
		{"test", []string{"b", "a"}},
	})
}

func TestRemoveAddNode(t *testing.T) {
	nodes := []string{"a", "b", "c"}
	ring := New(nodes)

	expectNodesABC(t, "1_", ring)
	expectNodeRangesABC(t, "2_", ring)

	ring = ring.RemoveNode("b")

	assertNodes(t, "3_", ring, []testPair{
		{"test", "a"},
		{"test", "a"},
		{"test1", "a"},
		{"test2", "a"},
		{"test3", "c"},
		{"test4", "a"},
		{"test5", "c"},
		{"aaaa", "c"},
		{"bbbb", "a"},
	})

	assert2Nodes(t, "4_", ring, []testNodes{
		{"test", []string{"a", "c"}},
		{"test", []string{"a", "c"}},
		{"test1", []string{"a", "c"}},
		{"test2", []string{"a", "c"}},
		{"test3", []string{"c", "a"}},
		{"test4", []string{"a", "c"}},
		{"test5", []string{"c", "a"}},
		{"aaaa", []string{"c", "a"}},
		{"bbbb", []string{"a", "c"}},
	})

	ring = ring.AddNode("b")

	expectNodesABC(t, "5_", ring)
	expectNodeRangesABC(t, "6_", ring)
}

func TestRemoveAddWeightedNode(t *testing.T) {
	weights := make(map[string]int)
	weights["a"] = 1
	weights["b"] = 2
	weights["c"] = 1
	ring := NewWithWeights(weights)

	expectWeights(t, ring, weights)

	assertNodes(t, "1_", ring, []testPair{
		{"test", "b"},
		{"test", "b"},
		{"test1", "b"},
		{"test2", "b"},
		{"test3", "c"},
		{"test4", "b"},
		{"test5", "c"},
		{"aaaa", "c"},
		{"bbbb", "b"},
	})

	assert2Nodes(t, "2_", ring, []testNodes{
		{"test", []string{"b", "a"}},
		{"test", []string{"b", "a"}},
		{"test1", []string{"b", "a"}},
		{"test2", []string{"b", "a"}},
		{"test3", []string{"c", "b"}},
		{"test4", []string{"b", "a"}},
		{"test5", []string{"c", "b"}},
		{"aaaa", []string{"c", "b"}},
		{"bbbb", []string{"b", "a"}},
	})

	ring = ring.RemoveNode("c")

	delete(weights, "c")
	expectWeights(t, ring, weights)

	assertNodes(t, "3_", ring, []testPair{
		{"test", "b"},
		{"test", "b"},
		{"test1", "b"},
		{"test2", "b"},
		{"test3", "b"},
		{"test4", "b"},
		{"test5", "b"},
		{"aaaa", "b"},
		{"bbbb", "b"},
	})

	assert2Nodes(t, "4_", ring, []testNodes{
		{"test", []string{"b", "a"}},
		{"test", []string{"b", "a"}},
		{"test1", []string{"b", "a"}},
		{"test2", []string{"b", "a"}},
		{"test3", []string{"b", "a"}},
		{"test4", []string{"b", "a"}},
		{"test5", []string{"b", "a"}},
		{"aaaa", []string{"b", "a"}},
		{"bbbb", []string{"b", "a"}},
	})
}

func TestAddRemoveNode(t *testing.T) {
	nodes := []string{"a", "b", "c"}
	ring := New(nodes)
	ring = ring.AddNode("d")

	expectNodesABCD(t, "1_", ring)

	assert2Nodes(t, "2_", ring, []testNodes{
		{"test", []string{"d", "a"}},
		{"test", []string{"d", "a"}},
		{"test1", []string{"b", "d"}},
		{"test2", []string{"b", "d"}},
		{"test3", []string{"c", "b"}},
		{"test4", []string{"d", "a"}},
		{"test5", []string{"c", "b"}},
		{"aaaa", []string{"c", "b"}},
		{"bbbb", []string{"d", "a"}},
	})

	ring = ring.AddNode("e")

	assertNodes(t, "3_", ring, []testPair{
		{"test", "a"},
		{"test", "a"},
		{"test1", "b"},
		{"test2", "b"},
		{"test3", "c"},
		{"test4", "c"},
		{"test5", "a"},
		{"aaaa", "b"},
		{"bbbb", "e"},
	})

	assert2Nodes(t, "4_", ring, []testNodes{
		{"test", []string{"d", "a"}},
		{"test", []string{"d", "a"}},
		{"test1", []string{"b", "d"}},
		{"test2", []string{"e", "b"}},
		{"test3", []string{"c", "e"}},
		{"test4", []string{"d", "a"}},
		{"test5", []string{"c", "e"}},
		{"aaaa", []string{"c", "e"}},
		{"bbbb", []string{"d", "a"}},
	})

	ring = ring.AddNode("f")

	assertNodes(t, "5_", ring, []testPair{
		{"test", "a"},
		{"test", "a"},
		{"test1", "b"},
		{"test2", "f"},
		{"test3", "f"},
		{"test4", "c"},
		{"test5", "f"},
		{"aaaa", "b"},
		{"bbbb", "e"},
	})

	assert2Nodes(t, "6_", ring, []testNodes{
		{"test", []string{"d", "a"}},
		{"test", []string{"d", "a"}},
		{"test1", []string{"b", "d"}},
		{"test2", []string{"e", "f"}},
		{"test3", []string{"c", "e"}},
		{"test4", []string{"d", "a"}},
		{"test5", []string{"c", "e"}},
		{"aaaa", []string{"c", "e"}},
		{"bbbb", []string{"d", "a"}},
	})

	ring = ring.RemoveNode("e")

	assertNodes(t, "7_", ring, []testPair{
		{"test", "a"},
		{"test", "a"},
		{"test1", "b"},
		{"test2", "f"},
		{"test3", "f"},
		{"test4", "c"},
		{"test5", "f"},
		{"aaaa", "b"},
		{"bbbb", "f"},
	})

	assert2Nodes(t, "8_", ring, []testNodes{
		{"test", []string{"d", "a"}},
		{"test", []string{"d", "a"}},
		{"test1", []string{"b", "d"}},
		{"test2", []string{"f", "b"}},
		{"test3", []string{"c", "f"}},
		{"test4", []string{"d", "a"}},
		{"test5", []string{"c", "f"}},
		{"aaaa", []string{"c", "f"}},
		{"bbbb", []string{"d", "a"}},
	})

	ring = ring.RemoveNode("f")

	expectNodesABCD(t, "TestAddRemoveNode_5_", ring)

	assert2Nodes(t, "", ring, []testNodes{
		{"test", []string{"d", "a"}},
		{"test", []string{"d", "a"}},
		{"test1", []string{"b", "d"}},
		{"test2", []string{"b", "d"}},
		{"test3", []string{"c", "b"}},
		{"test4", []string{"d", "a"}},
		{"test5", []string{"c", "b"}},
		{"aaaa", []string{"c", "b"}},
		{"bbbb", []string{"d", "a"}},
	})

	ring = ring.RemoveNode("d")

	expectNodesABC(t, "TestAddRemoveNode_6_", ring)
	expectNodeRangesABC(t, "", ring)
}
