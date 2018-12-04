// +build int

package main

// Hash keys.

import (
	"testing"
)

func TestHash(t *testing.T) {
	testCommands(t,
		succ("HSET", "aap", "noot", "mies"),
		succ("HGET", "aap", "noot"),
		succ("HMGET", "aap", "noot"),
		succ("HLEN", "aap"),
		succ("HKEYS", "aap"),
		succ("HVALS", "aap"),

		succ("HDEL", "aap", "noot"),
		succ("HGET", "aap", "noot"),
		succ("EXISTS", "aap"), // key is gone

		// failure cases
		fail("HSET", "aap", "noot"),
		fail("HGET", "aap"),
		fail("HMGET", "aap"),
		fail("HLEN"),
		fail("HKEYS"),
		fail("HVALS"),
		succ("SET", "str", "I am a string"),
		fail("HSET", "str", "noot", "mies"),
		fail("HGET", "str", "noot"),
		fail("HMGET", "str", "noot"),
		fail("HLEN", "str"),
		fail("HKEYS", "str"),
		fail("HVALS", "str"),
	)
}

func TestHashSetnx(t *testing.T) {
	testCommands(t,
		succ("HSETNX", "aap", "noot", "mies"),
		succ("EXISTS", "aap"),
		succ("HEXISTS", "aap", "noot"),

		succ("HSETNX", "aap", "noot", "mies2"),
		succ("HGET", "aap", "noot"),

		// failure cases
		fail("HSETNX", "aap"),
		fail("HSETNX", "aap", "noot"),
		fail("HSETNX", "aap", "noot", "too", "many"),
	)
}

func TestHashDelExists(t *testing.T) {
	testCommands(t,
		succ("HSET", "aap", "noot", "mies"),
		succ("HSET", "aap", "vuur", "wim"),
		succ("HEXISTS", "aap", "noot"),
		succ("HEXISTS", "aap", "vuur"),
		succ("HDEL", "aap", "noot"),
		succ("HEXISTS", "aap", "noot"),
		succ("HEXISTS", "aap", "vuur"),

		succ("HEXISTS", "nosuch", "vuur"),

		// failure cases
		fail("HDEL"),
		fail("HDEL", "aap"),
		succ("SET", "str", "I am a string"),
		fail("HDEL", "str", "key"),

		fail("HEXISTS"),
		fail("HEXISTS", "aap"),
		fail("HEXISTS", "aap", "too", "many"),
		fail("HEXISTS", "str", "field"),
	)
}

func TestHashGetall(t *testing.T) {
	testCommands(t,
		succ("HSET", "aap", "noot", "mies"),
		succ("HSET", "aap", "vuur", "wim"),
		succSorted("HGETALL", "aap"),

		succ("HGETALL", "nosuch"),

		// failure cases
		fail("HGETALL"),
		fail("HGETALL", "too", "many"),
		succ("SET", "str", "I am a string"),
		fail("HGETALL", "str"),
	)
}

func TestHmset(t *testing.T) {
	testCommands(t,
		succ("HMSET", "aap", "noot", "mies", "vuur", "zus"),
		succ("HGET", "aap", "noot"),
		succ("HGET", "aap", "vuur"),
		succ("HLEN", "aap"),

		// failure cases
		fail("HMSET", "aap"),
		fail("HMSET", "aap", "key"),
		fail("HMSET", "aap", "key", "value", "odd"),
		succ("SET", "str", "I am a string"),
		fail("HMSET", "str", "key", "value"),
	)
}

func TestHashIncr(t *testing.T) {
	testCommands(t,
		succ("HINCRBY", "aap", "noot", 12),
		succ("HINCRBY", "aap", "noot", -13),
		succ("HINCRBY", "aap", "noot", 2123),
		succ("HGET", "aap", "noot"),

		// Simple failure cases.
		fail("HINCRBY"),
		fail("HINCRBY", "aap"),
		fail("HINCRBY", "aap", "noot"),
		fail("HINCRBY", "aap", "noot", "noint"),
		fail("HINCRBY", "aap", "noot", 12, "toomany"),
		succ("SET", "str", "value"),
		fail("HINCRBY", "str", "value", 12),
		succ("HINCRBY", "aap", "noot", 12),
	)

	testCommands(t,
		succ("HINCRBYFLOAT", "aap", "noot", 12.3),
		succ("HINCRBYFLOAT", "aap", "noot", -13.1),
		succ("HINCRBYFLOAT", "aap", "noot", 200),
		succ("HGET", "aap", "noot"),

		// Simple failure cases.
		fail("HINCRBYFLOAT"),
		fail("HINCRBYFLOAT", "aap"),
		fail("HINCRBYFLOAT", "aap", "noot"),
		fail("HINCRBYFLOAT", "aap", "noot", "noint"),
		fail("HINCRBYFLOAT", "aap", "noot", 12, "toomany"),
		succ("SET", "str", "value"),
		fail("HINCRBYFLOAT", "str", "value", 12),
		succ("HINCRBYFLOAT", "aap", "noot", 12),
	)
}

func TestHscan(t *testing.T) {
	testCommands(t,
		// No set yet
		succ("HSCAN", "h", 0),

		succ("HSET", "h", "key1", "value1"),
		succ("HSCAN", "h", 0),
		succ("HSCAN", "h", 0, "COUNT", 12),
		succ("HSCAN", "h", 0, "cOuNt", 12),

		succ("HSET", "h", "anotherkey", "value2"),
		succ("HSCAN", "h", 0, "MATCH", "anoth*"),
		succ("HSCAN", "h", 0, "MATCH", "anoth*", "COUNT", 100),
		succ("HSCAN", "h", 0, "COUNT", 100, "MATCH", "anoth*"),

		// Can't really test multiple keys.
		// succ("SET", "key2", "value2"),
		// succ("SCAN", 0),

		// Error cases
		fail("HSCAN"),
		fail("HSCAN", "noint"),
		fail("HSCAN", "h", 0, "COUNT", "noint"),
		fail("HSCAN", "h", 0, "COUNT"),
		fail("HSCAN", "h", 0, "MATCH"),
		fail("HSCAN", "h", 0, "garbage"),
		fail("HSCAN", "h", 0, "COUNT", 12, "MATCH", "foo", "garbage"),
		// fail("HSCAN", "nosuch", 0, "COUNT", "garbage"),
		succ("SET", "str", "1"),
		fail("HSCAN", "str", 0),
	)
}
