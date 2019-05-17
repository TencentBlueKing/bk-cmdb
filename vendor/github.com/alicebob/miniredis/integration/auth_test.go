// +build int

package main

import (
	"testing"
)

func TestAuth(t *testing.T) {
	testAuthCommands(t,
		"supersecret",
		fail("PING"),
		fail("SET", "foo", "bar"),
		fail("SET"),
		fail("SET", "foo", "bar", "baz"),
		fail("GET", "foo"),
		fail("AUTH", "nosecret"),
		succ("AUTH", "supersecret"),
		succ("SET", "foo", "bar"),
		succ("GET", "foo"),
	)
}

func TestNoAuth(t *testing.T) {
	testCommands(t,
		fail("AUTH", "foo"),
	)
}
