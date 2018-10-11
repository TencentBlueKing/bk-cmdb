// +build int

package main

// List keys.

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
)

func TestLPush(t *testing.T) {
	testCommands(t,
		succ("LPUSH", "l", "aap", "noot", "mies"),
		succ("TYPE", "l"),
		succ("LPUSH", "l", "more", "keys"),
		succ("LRANGE", "l", 0, -1),
		succ("LRANGE", "l", 0, 6),
		succ("LRANGE", "l", 2, 6),
		succ("LRANGE", "l", -100, -100),
		succ("LRANGE", "nosuch", 2, 6),
		succ("LPOP", "l"),
		succ("LPOP", "l"),
		succ("LPOP", "l"),
		succ("LPOP", "l"),
		succ("LPOP", "l"),
		succ("LPOP", "l"),
		succ("EXISTS", "l"),
		succ("LPOP", "nosuch"),

		// failure cases
		fail("LPUSH"),
		fail("LPUSH", "l"),
		succ("SET", "str", "I am a string"),
		fail("LPUSH", "str", "noot", "mies"),
		fail("LRANGE"),
		fail("LRANGE", "key"),
		fail("LRANGE", "key", 2),
		fail("LRANGE", "key", 2, 6, "toomany"),
		fail("LRANGE", "key", "noint", 6),
		fail("LRANGE", "key", 2, "noint"),
		fail("LPOP"),
		fail("LPOP", "key", "args"),
	)
}

func TestLPushx(t *testing.T) {
	testCommands(t,
		succ("LPUSHX", "l", "aap"),
		succ("EXISTS", "l"),
		succ("LRANGE", "l", 0, -1),
		succ("LPUSH", "l", "noot"),
		succ("LPUSHX", "l", "mies"),
		succ("EXISTS", "l"),
		succ("LRANGE", "l", 0, -1),
		succ("LPUSHX", "l", "even", "more", "arguments"),

		// failure cases
		fail("LPUSHX"),
		fail("LPUSHX", "l"),
		succ("SET", "str", "I am a string"),
		fail("LPUSHX", "str", "mies"),
	)
}

func TestRPush(t *testing.T) {
	testCommands(t,
		succ("RPUSH", "l", "aap", "noot", "mies"),
		succ("TYPE", "l"),
		succ("RPUSH", "l", "more", "keys"),
		succ("LRANGE", "l", 0, -1),
		succ("LRANGE", "l", 0, 6),
		succ("LRANGE", "l", 2, 6),
		succ("RPOP", "l"),
		succ("RPOP", "l"),
		succ("RPOP", "l"),
		succ("RPOP", "l"),
		succ("RPOP", "l"),
		succ("RPOP", "l"),
		succ("EXISTS", "l"),
		succ("RPOP", "nosuch"),

		// failure cases
		fail("RPUSH"),
		fail("RPUSH", "l"),
		succ("SET", "str", "I am a string"),
		fail("RPUSH", "str", "noot", "mies"),
		fail("RPOP"),
		fail("RPOP", "key", "args"),
	)
}

func TestLinxed(t *testing.T) {
	testCommands(t,
		succ("RPUSH", "l", "aap", "noot", "mies"),
		succ("LINDEX", "l", 0),
		succ("LINDEX", "l", 1),
		succ("LINDEX", "l", 2),
		succ("LINDEX", "l", 3),
		succ("LINDEX", "l", 4),
		succ("LINDEX", "l", 44444),
		succ("LINDEX", "l", -0),
		succ("LINDEX", "l", -1),
		succ("LINDEX", "l", -2),
		succ("LINDEX", "l", -3),
		succ("LINDEX", "l", -4),
		succ("LINDEX", "l", -4000),

		// failure cases
		fail("LINDEX"),
		fail("LINDEX", "l"),
		succ("SET", "str", "I am a string"),
		fail("LINDEX", "str", 1),
		fail("LINDEX", "l", "noint"),
		fail("LINDEX", "l", 1, "too many"),
	)
}

func TestLlen(t *testing.T) {
	testCommands(t,
		succ("RPUSH", "l", "aap", "noot", "mies"),
		succ("LLEN", "l"),
		succ("LLEN", "nosuch"),

		// failure cases
		succ("SET", "str", "I am a string"),
		fail("LLEN", "str"),
		fail("LLEN"),
		fail("LLEN", "l", "too many"),
	)
}

func TestLtrim(t *testing.T) {
	testCommands(t,
		succ("RPUSH", "l", "aap", "noot", "mies"),
		succ("LTRIM", "l", 0, 1),
		succ("LRANGE", "l", 0, -1),
		succ("RPUSH", "l2", "aap", "noot", "mies", "vuur"),
		succ("LTRIM", "l2", -2, -1),
		succ("LRANGE", "l2", 0, -1),
		succ("RPUSH", "l3", "aap", "noot", "mies", "vuur"),
		succ("LTRIM", "l3", -2, -1000),
		succ("LRANGE", "l3", 0, -1),

		// remove the list
		succ("RPUSH", "l4", "aap"),
		succ("LTRIM", "l4", 0, -999),
		succ("EXISTS", "l4"),

		// failure cases
		succ("SET", "str", "I am a string"),
		fail("LTRIM", "str", 0, 1),
		fail("LTRIM", "l", 0, 1, "toomany"),
		fail("LTRIM", "l", "noint", 1),
		fail("LTRIM", "l", 0, "noint"),
		fail("LTRIM", "l", 0),
		fail("LTRIM", "l"),
		fail("LTRIM"),
	)
}

func TestLrem(t *testing.T) {
	testCommands(t,
		succ("RPUSH", "l", "aap", "noot", "mies", "mies", "mies"),
		succ("LREM", "l", 1, "mies"),
		succ("LRANGE", "l", 0, -1),
		succ("RPUSH", "l2", "aap", "noot", "mies", "mies", "mies"),
		succ("LREM", "l2", -2, "mies"),
		succ("LRANGE", "l2", 0, -1),
		succ("RPUSH", "l3", "aap", "noot", "mies", "mies", "mies"),
		succ("LREM", "l3", 0, "mies"),
		succ("LRANGE", "l3", 0, -1),

		// remove the list
		succ("RPUSH", "l4", "aap"),
		succ("LREM", "l4", 999, "aap"),
		succ("EXISTS", "l4"),

		// failure cases
		succ("SET", "str", "I am a string"),
		fail("LREM", "str", 0, "aap"),
		fail("LREM", "l", 0, "aap", "toomany"),
		fail("LREM", "l", "noint", "aap"),
		fail("LREM", "l", 0),
		fail("LREM", "l"),
		fail("LREM"),
	)
}

func TestLset(t *testing.T) {
	testCommands(t,
		succ("RPUSH", "l", "aap", "noot", "mies", "mies", "mies"),
		succ("LSET", "l", 1, "[cencored]"),
		succ("LRANGE", "l", 0, -1),
		succ("LSET", "l", -1, "[cencored]"),
		succ("LRANGE", "l", 0, -1),
		fail("LSET", "l", 1000, "new"),
		fail("LSET", "l", -7000, "new"),
		fail("LSET", "nosuch", 1, "new"),

		// failure cases
		fail("LSET"),
		fail("LSET", "l"),
		fail("LSET", "l", 0),
		fail("LSET", "l", "noint", "aap"),
		fail("LSET", "l", 0, "aap", "toomany"),
		succ("SET", "str", "I am a string"),
		fail("LSET", "str", 0, "aap"),
	)
}

func TestLinsert(t *testing.T) {
	testCommands(t,
		succ("RPUSH", "l", "aap", "noot", "mies", "mies", "mies!"),
		succ("LINSERT", "l", "before", "aap", "1"),
		succ("LINSERT", "l", "before", "noot", "2"),
		succ("LINSERT", "l", "after", "mies!", "3"),
		succ("LINSERT", "l", "after", "mies", "4"),
		succ("LINSERT", "l", "after", "nosuch", "0"),
		succ("LINSERT", "nosuch", "after", "nosuch", "0"),
		succ("LRANGE", "l", 0, -1),
		succ("LINSERT", "l", "AfTeR", "mies", "4"),
		succ("LRANGE", "l", 0, -1),

		// failure cases
		fail("LINSERT"),
		fail("LINSERT", "l"),
		fail("LINSERT", "l", "before"),
		fail("LINSERT", "l", "before", "aap"),
		fail("LINSERT", "l", "before", "aap", "too", "many"),
		fail("LINSERT", "l", "What?", "aap", "noot"),
		succ("SET", "str", "I am a string"),
		fail("LINSERT", "str", "before", "aap", "noot"),
	)
}

func TestRpoplpush(t *testing.T) {
	testCommands(t,
		succ("RPUSH", "l", "aap", "noot", "mies"),
		succ("RPOPLPUSH", "l", "l2"),
		succ("LRANGE", "l", 0, -1),
		succ("LRANGE", "2l", 0, -1),
		succ("RPOPLPUSH", "l", "l2"),
		succ("RPOPLPUSH", "l", "l2"),
		succ("RPOPLPUSH", "l", "l2"), // now empty
		succ("EXISTS", "l"),
		succ("LRANGE", "2l", 0, -1),

		succ("RPUSH", "round", "aap", "noot", "mies"),
		succ("RPOPLPUSH", "round", "round"),
		succ("LRANGE", "round", 0, -1),
		succ("RPOPLPUSH", "round", "round"),
		succ("RPOPLPUSH", "round", "round"),
		succ("RPOPLPUSH", "round", "round"),
		succ("RPOPLPUSH", "round", "round"),
		succ("LRANGE", "round", 0, -1),

		// failure cases
		succ("RPUSH", "chk", "aap", "noot", "mies"),
		fail("RPOPLPUSH"),
		fail("RPOPLPUSH", "chk"),
		fail("RPOPLPUSH", "chk", "too", "many"),
		succ("SET", "str", "I am a string"),
		fail("RPOPLPUSH", "chk", "str"),
		fail("RPOPLPUSH", "str", "chk"),
		succ("LRANGE", "chk", 0, -1),
	)
}

func TestRpushx(t *testing.T) {
	testCommands(t,
		succ("RPUSHX", "l", "aap"),
		succ("EXISTS", "l"),
		succ("RPUSH", "l", "noot", "mies"),
		succ("RPUSHX", "l", "vuur"),
		succ("EXISTS", "l"),
		succ("LRANGE", "l", 0, -1),
		succ("RPUSHX", "l", "more", "arguments"),

		// failure cases
		succ("RPUSH", "chk", "noot", "mies"),
		fail("RPUSHX"),
		fail("RPUSHX", "chk"),
		succ("LRANGE", "chk", 0, -1),
		succ("SET", "str", "I am a string"),
		fail("RPUSHX", "str", "value"),
	)
}

func TestBrpop(t *testing.T) {
	testCommands(t,
		succ("LPUSH", "l", "one"),
		succ("BRPOP", "l", 1),
		succ("EXISTS", "l"),

		// transaction
		succ("MULTI"),
		succ("BRPOP", "nosuch", 10),
		succ("EXEC"),

		// failure cases
		fail("BRPOP"),
		fail("BRPOP", "l"),
		fail("BRPOP", "l", "X"),
		fail("BRPOP", "l", ""),
		fail("BRPOP", 1),
		fail("BRPOP", "key", -1),
	)
}

func TestBrpopMulti(t *testing.T) {
	testMultiCommands(t,
		func(r chan<- command, _ *miniredis.Miniredis) {
			r <- succ("BRPOP", "key", 1)
			r <- succ("BRPOP", "key", 1)
			r <- succ("BRPOP", "key", 1)
			r <- succ("BRPOP", "key", 1)
			r <- succ("BRPOP", "key", 1) // will timeout
		},
		func(r chan<- command, _ *miniredis.Miniredis) {
			r <- succ("LPUSH", "key", "aap", "noot", "mies")
			time.Sleep(50 * time.Millisecond)
			r <- succ("LPUSH", "key", "toon")
		},
	)
}

func TestBrpopTrans(t *testing.T) {
	testMultiCommands(t,
		func(r chan<- command, _ *miniredis.Miniredis) {
			r <- succ("BRPOP", "key", 1)
		},
		func(r chan<- command, _ *miniredis.Miniredis) {
			r <- succ("MULTI")
			r <- succ("LPUSH", "key", "toon")
			r <- succ("EXEC")
		},
	)
}

func TestBlpop(t *testing.T) {
	testCommands(t,
		succ("LPUSH", "l", "one"),
		succ("BLPOP", "l", 1),
		succ("EXISTS", "l"),

		// failure cases
		fail("BLPOP"),
		fail("BLPOP", "l"),
		fail("BLPOP", "l", "X"),
		fail("BLPOP", "l", ""),
		fail("BLPOP", 1),
		fail("BLPOP", "key", -1),
	)

	testMultiCommands(t,
		func(r chan<- command, _ *miniredis.Miniredis) {
			r <- succ("BLPOP", "key", 1)
			r <- succ("BLPOP", "key", 1)
			r <- succ("BLPOP", "key", 1)
			r <- succ("BLPOP", "key", 1)
			r <- succ("BLPOP", "key", 1) // will timeout
		},
		func(r chan<- command, _ *miniredis.Miniredis) {
			r <- succ("LPUSH", "key", "aap", "noot", "mies")
			time.Sleep(10 * time.Millisecond)
			r <- succ("LPUSH", "key", "toon")
		},
	)
}

func TestBrpoplpush(t *testing.T) {
	testCommands(t,
		succ("LPUSH", "l", "one"),
		succ("BRPOPLPUSH", "l", "l2", 1),
		succ("EXISTS", "l"),
		succ("EXISTS", "l2"),
		succ("LRANGE", "l", 0, -1),
		succ("LRANGE", "l2", 0, -1),

		// failure cases
		fail("BRPOPLPUSH"),
		fail("BRPOPLPUSH", "l"),
		fail("BRPOPLPUSH", "l", "x"),
		fail("BRPOPLPUSH", 1),
		fail("BRPOPLPUSH", "from", "to", -1),
		fail("BRPOPLPUSH", "from", "to", -1, "xxx"),
	)
	testMultiCommands(t,
		func(r chan<- command, _ *miniredis.Miniredis) {
			r <- succ("BRPOPLPUSH", "from", "to", 1)
			r <- succ("BRPOPLPUSH", "from", "to", 1)
			r <- succ("BRPOPLPUSH", "from", "to", 1)
			r <- succ("BRPOPLPUSH", "from", "to", 1)
			r <- succ("BRPOPLPUSH", "from", "to", 1) // will timeout
		},
		func(r chan<- command, _ *miniredis.Miniredis) {
			r <- succ("LPUSH", "from", "aap", "noot", "mies")
			time.Sleep(10 * time.Millisecond)
			r <- succ("LPUSH", "from", "toon")
			r <- succ("LRANGE", "from", 0, -1)
			r <- succ("LRANGE", "to", 0, -1)
		},
	)
}
