// +build int

package main

import (
	"testing"
)

func TestServer(t *testing.T) {
	testCommands(t,
		succ("SET", "foo", "bar"),
		succ("SET", "baz", "bak"),
		succ("DBSIZE"),
		succ("SELECT", 2),
		succ("DBSIZE"),
		succ("SET", "baz", "bak"),

		succ("SELECT", 0),
		succ("FLUSHDB"),
		succ("DBSIZE"),

		succ("SELECT", 2),
		succ("DBSIZE"),
		succ("FLUSHALL"),
		succ("DBSIZE"),

		succ("FLUSHDB", "aSyNc"),
		succ("FLUSHALL", "AsYnC"),

		// Failure cases
		fail("DBSIZE", "foo"),
		fail("FLUSHDB", "foo"),
		fail("FLUSHALL", "foo"),
		fail("FLUSHDB", "ASYNC", "foo"),
		fail("FLUSHDB", "ASYNC", "ASYNC"),
		fail("FLUSHALL", "ASYNC", "foo"),
	)
}
