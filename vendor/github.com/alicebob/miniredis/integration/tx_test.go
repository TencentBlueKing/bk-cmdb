// +build int

package main

// Transaction things.

import (
	"testing"
)

func TestTx(t *testing.T) {
	testCommands(t,
		succ("MULTI"),
		succ("SET", "AAP", 1),
		succ("GET", "AAP"),
		succ("EXEC"),
		succ("GET", "AAP"),
	)

	// err: Double MULTI
	testCommands(t,
		succ("MULTI"),
		fail("MULTI"),
	)

	// err: No MULTI
	testCommands(t,
		fail("EXEC"),
	)

	// Errors in the MULTI sequence
	testCommands(t,
		succ("MULTI"),
		succ("SET", "foo", "bar"),
		fail("SET", "foo"),
		succ("SET", "foo", "bar"),
		fail("EXEC"),
	)

	// Simple WATCH
	testCommands(t,
		succ("SET", "foo", "bar"),
		succ("WATCH", "foo"),
		succ("MULTI"),
		succ("GET", "foo"),
		succ("EXEC"),
	)

	// Simple UNWATCH
	testCommands(t,
		succ("SET", "foo", "bar"),
		succ("WATCH", "foo"),
		succ("UNWATCH"),
		succ("MULTI"),
		succ("GET", "foo"),
		succ("EXEC"),
	)

	// UNWATCH in a MULTI. Yep. Weird.
	testCommands(t,
		succ("WATCH", "foo"),
		succ("MULTI"),
		succ("UNWATCH"), // Valid. Somehow.
		succ("EXEC"),
	)

	// Test whether all these commands support transactions.
	testCommands(t,
		succ("MULTI"),
		succ("GET", "str"),
		succ("SET", "str", "bar"),
		succ("SETNX", "str", "bar"),
		succ("GETSET", "str", "bar"),
		succ("MGET", "str", "bar"),
		succ("MSET", "str", "bar"),
		succ("MSETNX", "str", "bar"),
		succ("SETEX", "str", 12, "newv"),
		succ("PSETEX", "str", 12, "newv"),
		succ("STRLEN", "str"),
		succ("APPEND", "str", "more"),
		succ("GETRANGE", "str", 0, 2),
		succ("SETRANGE", "str", 0, "B"),
		succ("EXEC"),
		succ("GET", "str"),
	)

	testCommands(t,
		succ("MULTI"),
		succ("SET", "bits", "\xff\x00"),
		succ("BITCOUNT", "bits"),
		succ("BITOP", "OR", "bits", "bits", "nosuch"),
		succ("BITPOS", "bits", 1),
		succ("GETBIT", "bits", 12),
		succ("SETBIT", "bits", 12, 1),
		succ("EXEC"),
		succ("GET", "bits"),
	)

	testCommands(t,
		succ("MULTI"),
		succ("INCR", "number"),
		succ("INCRBY", "number", 12),
		succ("INCRBYFLOAT", "number", 12.2),
		succ("DECR", "number"),
		succ("GET", "number"),
		succ("DECRBY", "number", 2),
		succ("GET", "number"),
	)

	testCommands(t,
		succ("MULTI"),
		succ("HSET", "hash", "foo", "bar"),
		succ("HDEL", "hash", "foo"),
		succ("HEXISTS", "hash", "foo"),
		succ("HSET", "hash", "foo", "bar22"),
		succ("HSETNX", "hash", "foo", "bar22"),
		succ("HGET", "hash", "foo"),
		succ("HMGET", "hash", "foo", "baz"),
		succ("HLEN", "hash"),
		succ("HGETALL", "hash"),
		succ("HKEYS", "hash"),
		succ("HVALS", "hash"),
	)

	testCommands(t,
		succ("MULTI"),
		succ("SET", "key", "foo"),
		succ("TYPE", "key"),
		succ("EXPIRE", "key", 12),
		succ("TTL", "key"),
		succ("PEXPIRE", "key", 12),
		succ("PTTL", "key"),
		succ("PERSIST", "key"),
		succ("DEL", "key"),
		succ("TYPE", "key"),
		succ("EXEC"),
	)

	// BITOP OPs are checked after the transaction.
	testCommands(t,
		succ("MULTI"),
		succ("BITOP", "BROKEN", "str", ""),
		succ("EXEC"),
	)
}
