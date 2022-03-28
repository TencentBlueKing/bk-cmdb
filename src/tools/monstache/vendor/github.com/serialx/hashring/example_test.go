package hashring

import (
	"crypto/sha1"
	"fmt"
)

func ExampleGetAllNodes() {
	hashRing := New([]string{"node1", "node2", "node3"})
	nodes, _ := hashRing.GetNodes("key", hashRing.Size())
	fmt.Printf("%v", nodes)
	// Output: [node3 node2 node1]
}

func ExampleCustomHashError() {
	_, err := NewHash(sha1.New).Use(NewInt64PairHashKey)
	fmt.Printf("%s", err.Error())
	// Output: can't use given hash.Hash with given hashKeyFunc: expected 16 bytes, got 20 bytes
}

func ExampleCustomHash() {
	hashFunc, _ := NewHash(sha1.New).FirstBytes(16).Use(NewInt64PairHashKey)
	hashRing := NewWithHash([]string{"node1", "node2", "node3"}, hashFunc)
	nodes, _ := hashRing.GetNodes("key", hashRing.Size())
	fmt.Printf("%v", nodes)
	// Output: [node2 node1 node3]
}

func ExampleNewHashFunc() {
	hashFunc, _ := NewHash(sha1.New).FirstBytes(16).Use(NewInt64PairHashKey)
	fmt.Printf("%v\n", hashFunc([]byte("test")))
	// Output: &{-6441359348440544599 -8653224871661646820}
}
