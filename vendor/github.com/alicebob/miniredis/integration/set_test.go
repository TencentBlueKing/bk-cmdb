// +build int

package main

// Set keys.

import (
	"testing"
)

func TestSet(t *testing.T) {
	testCommands(t,
		succ("SADD", "s", "aap", "noot", "mies"),
		succ("SADD", "s", "vuur", "noot"),
		succ("TYPE", "s"),
		succ("EXISTS", "s"),
		succ("SCARD", "s"),
		succSorted("SMEMBERS", "s"),
		succSorted("SMEMBERS", "nosuch"),
		succ("SISMEMBER", "s", "aap"),
		succ("SISMEMBER", "s", "nosuch"),

		succ("SCARD", "nosuch"),
		succ("SISMEMBER", "nosuch", "nosuch"),

		// failure cases
		fail("SADD"),
		fail("SADD", "s"),
		fail("SMEMBERS"),
		fail("SMEMBERS", "too", "many"),
		fail("SCARD"),
		fail("SCARD", "too", "many"),
		fail("SISMEMBER"),
		fail("SISMEMBER", "few"),
		fail("SISMEMBER", "too", "many", "arguments"),
		// Wrong type
		succ("SET", "str", "I am a string"),
		fail("SADD", "str", "noot", "mies"),
		fail("SMEMBERS", "str"),
		fail("SISMEMBER", "str", "noot"),
		fail("SCARD", "str"),
	)
}

func TestSetMove(t *testing.T) {
	// Move a set around
	testCommands(t,
		succ("SADD", "s", "aap", "noot", "mies"),
		succ("RENAME", "s", "others"),
		succSorted("SMEMBERS", "s"),
		succSorted("SMEMBERS", "others"),
		succ("MOVE", "others", 2),
		succSorted("SMEMBERS", "others"),
		succ("SELECT", 2),
		succSorted("SMEMBERS", "others"),
	)
}

func TestSetDel(t *testing.T) {
	testCommands(t,
		succ("SADD", "s", "aap", "noot", "mies"),
		succ("SREM", "s", "noot", "nosuch"),
		succ("SCARD", "s"),
		succSorted("SMEMBERS", "s"),

		// failure cases
		fail("SREM"),
		fail("SREM", "s"),
		// Wrong type
		succ("SET", "str", "I am a string"),
		fail("SREM", "str", "noot"),
	)
}

func TestSetSMove(t *testing.T) {
	testCommands(t,
		succ("SADD", "s", "aap", "noot", "mies"),
		succ("SMOVE", "s", "s2", "aap"),
		succ("SCARD", "s"),
		succ("SCARD", "s2"),
		succ("SMOVE", "s", "s2", "nosuch"),
		succ("SCARD", "s"),
		succ("SCARD", "s2"),
		succ("SMOVE", "s", "nosuch", "noot"),
		succ("SCARD", "s"),
		succ("SCARD", "s2"),

		succ("SMOVE", "s", "s2", "mies"),
		succ("SCARD", "s"),
		succ("EXISTS", "s"),
		succ("SCARD", "s2"),
		succ("EXISTS", "s2"),

		succ("SMOVE", "s2", "s2", "mies"),

		succ("SADD", "s5", "aap"),
		succ("SADD", "s6", "aap"),
		succ("SMOVE", "s5", "s6", "aap"),

		// failure cases
		fail("SMOVE"),
		fail("SMOVE", "s"),
		fail("SMOVE", "s", "s2"),
		fail("SMOVE", "s", "s2", "too", "many"),
		// Wrong type
		succ("SET", "str", "I am a string"),
		fail("SMOVE", "str", "s2", "noot"),
		fail("SMOVE", "s2", "str", "noot"),
	)
}

func TestSetSpop(t *testing.T) {
	testCommands(t,
		// Without count argument
		succ("SADD", "s", "aap"),
		succ("SPOP", "s"),
		succ("EXISTS", "s"),

		succ("SPOP", "nosuch"),

		succ("SADD", "s", "aap"),
		succ("SADD", "s", "noot"),
		succ("SADD", "s", "mies"),
		succ("SADD", "s", "noot"),
		succ("SCARD", "s"),
		succLoosely("SMEMBERS", "s"),

		// failure cases
		fail("SPOP"),
		succ("SADD", "s", "aap"),
		fail("SPOP", "s", "s2"),
		fail("SPOP", "nosuch", "s2"),
		// Wrong type
		succ("SET", "str", "I am a string"),
		fail("SPOP", "str"),
	)

	testCommands(t,
		// With count argument
		succ("SADD", "s", "aap"),
		succ("SADD", "s", "noot"),
		succ("SADD", "s", "mies"),
		succ("SADD", "s", "vuur"),
		succLoosely("SPOP", "s", 2),
		succ("EXISTS", "s"),
		succ("SCARD", "s"),

		succLoosely("SPOP", "s", 200),
		succ("SPOP", "s", 1),
		succ("SCARD", "s"),

		// failure cases
		fail("SPOP", "foo", "one"),
	)
}

func TestSetSrandmember(t *testing.T) {
	testCommands(t,
		// Set with a single member...
		succ("SADD", "s", "aap"),
		succ("SRANDMEMBER", "s"),
		succ("SRANDMEMBER", "s", 1),
		succ("SRANDMEMBER", "s", 5),
		succ("SRANDMEMBER", "s", -1),
		succ("SRANDMEMBER", "s", -5),

		succ("SRANDMEMBER", "s", 0),
		succ("SPOP", "nosuch"),

		// failure cases
		fail("SRANDMEMBER"),
		fail("SRANDMEMBER", "s", "noint"),
		fail("SRANDMEMBER", "s", 1, "toomany"),
		// Wrong type
		succ("SET", "str", "I am a string"),
		fail("SRANDMEMBER", "str"),
	)
}

func TestSetSdiff(t *testing.T) {
	testCommands(t,
		succ("SDIFF", "s1", "aap", "noot", "mies"),
		succ("SDIFF", "s2", "noot", "mies", "vuur"),
		succ("SDIFF", "s3", "mies", "wim"),
		succ("SDIFF", "s1"),
		succ("SDIFF", "s1", "s2"),
		succ("SDIFF", "s1", "s2", "s3"),
		succ("SDIFF", "nosuch"),
		succ("SDIFF", "s1", "nosuch", "s2", "nosuch", "s3"),
		succ("SDIFF", "s1", "s1"),

		succ("SDIFFSTORE", "res", "s3", "nosuch", "s1"),
		succ("SMEMBERS", "res"),

		// failure cases
		fail("SDIFF"),
		fail("SDIFFSTORE"),
		fail("SDIFFSTORE", "key"),
		// Wrong type
		succ("SET", "str", "I am a string"),
		fail("SDIFF", "s1", "str"),
		fail("SDIFF", "nosuch", "str"),
		fail("SDIFF", "str", "s1"),
		fail("SDIFFSTORE", "res", "str", "s1"),
		fail("SDIFFSTORE", "res", "s1", "str"),
	)
}

func TestSetSinter(t *testing.T) {
	testCommands(t,
		succ("SINTER", "s1", "aap", "noot", "mies"),
		succ("SINTER", "s2", "noot", "mies", "vuur"),
		succ("SINTER", "s3", "mies", "wim"),
		succ("SINTER", "s1"),
		succ("SINTER", "s1", "s2"),
		succ("SINTER", "s1", "s2", "s3"),
		succ("SINTER", "nosuch"),
		succ("SINTER", "s1", "nosuch", "s2", "nosuch", "s3"),
		succ("SINTER", "s1", "s1"),

		succ("SINTERSTORE", "res", "s3", "nosuch", "s1"),
		succ("SMEMBERS", "res"),

		// failure cases
		fail("SINTER"),
		fail("SINTERSTORE"),
		fail("SINTERSTORE", "key"),
		// Wrong type
		succ("SET", "str", "I am a string"),
		succ("SINTER", "s1", "str"),     // !
		succ("SINTER", "nosuch", "str"), // !
		fail("SINTER", "str", "s1"),
		fail("SINTERSTORE", "res", "str", "s1"),
		succ("SINTERSTORE", "res", "s1", "str"), // !
	)
}

func TestSetSunion(t *testing.T) {
	testCommands(t,
		succ("SUNION", "s1", "aap", "noot", "mies"),
		succ("SUNION", "s2", "noot", "mies", "vuur"),
		succ("SUNION", "s3", "mies", "wim"),
		succ("SUNION", "s1"),
		succ("SUNION", "s1", "s2"),
		succ("SUNION", "s1", "s2", "s3"),
		succ("SUNION", "nosuch"),
		succ("SUNION", "s1", "nosuch", "s2", "nosuch", "s3"),
		succ("SUNION", "s1", "s1"),

		succ("SUNIONSTORE", "res", "s3", "nosuch", "s1"),
		succ("SMEMBERS", "res"),

		// failure cases
		fail("SUNION"),
		fail("SUNIONSTORE"),
		fail("SUNIONSTORE", "key"),
		// Wrong type
		succ("SET", "str", "I am a string"),
		fail("SUNION", "s1", "str"),
		fail("SUNION", "nosuch", "str"),
		fail("SUNION", "str", "s1"),
		fail("SUNIONSTORE", "res", "str", "s1"),
		fail("SUNIONSTORE", "res", "s1", "str"),
	)
}

func TestSscan(t *testing.T) {
	testCommands(t,
		// No set yet
		succ("SSCAN", "set", 0),

		succ("SADD", "set", "key1"),
		succ("SSCAN", "set", 0),
		succ("SSCAN", "set", 0, "COUNT", 12),
		succ("SSCAN", "set", 0, "cOuNt", 12),

		succ("SADD", "set", "anotherkey"),
		succ("SSCAN", "set", 0, "MATCH", "anoth*"),
		succ("SSCAN", "set", 0, "MATCH", "anoth*", "COUNT", 100),
		succ("SSCAN", "set", 0, "COUNT", 100, "MATCH", "anoth*"),

		// Can't really test multiple keys.
		// succ("SET", "key2", "value2"),
		// succ("SCAN", 0),

		// Error cases
		fail("SSCAN"),
		fail("SSCAN", "noint"),
		fail("SSCAN", "set", 0, "COUNT", "noint"),
		fail("SSCAN", "set", 0, "COUNT"),
		fail("SSCAN", "set", 0, "MATCH"),
		fail("SSCAN", "set", 0, "garbage"),
		fail("SSCAN", "set", 0, "COUNT", 12, "MATCH", "foo", "garbage"),
		succ("SET", "str", "1"),
		fail("SSCAN", "str", 0),
	)
}

func TestSetNoAuth(t *testing.T) {
	testAuthCommands(t,
		"supersecret",
		failWith(
			"NOAUTH Authentication required.",
			"SET", "foo", "bar",
		),
		succ("AUTH", "supersecret"),
		succ(
			"SET", "foo", "bar",
		),
	)
}
