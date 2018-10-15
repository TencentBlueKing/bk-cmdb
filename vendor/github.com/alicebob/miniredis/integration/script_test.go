// +build int

package main

// Script commands

import (
	"testing"
)

func TestEval(t *testing.T) {
	testCommands(t,
		succ("EVAL", "return 42", 0),
		succ("EVAL", "", 0),
		succ("EVAL", "return 42", 1, "foo"),
		succ("EVAL", "return {KEYS[1],KEYS[2],ARGV[1],ARGV[2]}", 2, "key1", "key2", "first", "second"),
		succ("EVAL", "return {ARGV[1]}", 0, "first"),
		succ("EVAL", "return {ARGV[1]}", 0, "first\nwith\nnewlines!\r\r\n\t!"),
		succ("EVAL", "return redis.call('GET', 'nosuch')==false", 0),
		succ("EVAL", "return redis.call('GET', 'nosuch')==nil", 0),
		succ("EVAL", "local a = redis.call('MGET', 'bar'); return a[1] == false", 0),
		succ("EVAL", "local a = redis.call('MGET', 'bar'); return a[1] == nil", 0),

		// failure cases
		fail("EVAL"),
		fail("EVAL", "return 42"),
		fail("EVAL", "["),
		fail("EVAL", "return 42", "return 43"),
		fail("EVAL", "return 42", 1),
		fail("EVAL", "return 42", -1),
		fail("EVAL", 42),
	)
}

func TestScript(t *testing.T) {
	testCommands(t,
		succ("SCRIPT", "LOAD", "return 42"),
		succ("SCRIPT", "LOAD", "return 42"),
		succ("SCRIPT", "LOAD", "return 43"),

		succ("SCRIPT", "EXISTS", "1fa00e76656cc152ad327c13fe365858fd7be306"),
		succ("SCRIPT", "EXISTS", "0", "1fa00e76656cc152ad327c13fe365858fd7be306"),
		succ("SCRIPT", "EXISTS", 0),
		succ("SCRIPT", "EXISTS"),

		succ("SCRIPT", "FLUSH"),
		succ("SCRIPT", "EXISTS", "1fa00e76656cc152ad327c13fe365858fd7be306"),

		fail("SCRIPT"),
		fail("SCRIPT", "LOAD", "return 42", "return 42"),
		failLoosely("SCRIPT", "LOAD", "]"),
		fail("SCRIPT", "LOAD", "]", "foo"),
		fail("SCRIPT", "LOAD"),
		fail("SCRIPT", "FLUSH", "foo"),
		fail("SCRIPT", "FOO"),
	)
}

func TestEvalsha(t *testing.T) {
	sha1 := "1fa00e76656cc152ad327c13fe365858fd7be306" // "return 42"
	sha2 := "bfbf458525d6a0b19200bfd6db3af481156b367b" // keys[1], argv[1]

	testCommands(t,
		succ("SCRIPT", "LOAD", "return 42"),
		succ("SCRIPT", "LOAD", "return {KEYS[1],ARGV[1]}"),
		succ("EVALSHA", sha1, "0"),
		succ("EVALSHA", sha2, "0"),
		succ("EVALSHA", sha2, "0", "foo"),
		succ("EVALSHA", sha2, "1", "foo"),
		succ("EVALSHA", sha2, "1", "foo", "bar"),
		succ("EVALSHA", sha2, "1", "foo", "bar", "baz"),

		succ("SCRIPT", "FLUSH"),
		fail("EVALSHA", sha1, "0"),

		succ("SCRIPT", "LOAD", "return 42"),
		fail("EVALSHA", sha1),
		fail("EVALSHA"),
		fail("EVALSHA", "nosuch"),
		fail("EVALSHA", "nosuch", 0),
	)
}

func TestLua(t *testing.T) {
	// basic datatype things
	testCommands(t,
		succ("EVAL", "", 0),
		succ("EVAL", "return 42", 0),
		succ("EVAL", "return 42, 43", 0),
		succ("EVAL", "return true", 0),
		succ("EVAL", "return 'foo'", 0),
		succ("EVAL", "return 3.1415", 0),
		succ("EVAL", "return 3.9999", 0),
		succ("EVAL", "return {1,'foo'}", 0),
		succ("EVAL", "return {1,'foo',nil,'foo'}", 0),
		succ("EVAL", "return 3.9999+3", 0),
		succ("EVAL", "return 3.99+0.0001", 0),
		succ("EVAL", "return 3.9999+0.201", 0),
		succ("EVAL", "return {{1}}", 0),
		succ("EVAL", "return {1,{1,{1,'bar'}}}", 0),
	)

	// special returns
	testCommands(t,
		fail("EVAL", "return {err = 'oops'}", 0),
		succ("EVAL", "return {1,{err = 'oops'}}", 0),
		fail("EVAL", "return redis.error_reply('oops')", 0),
		succ("EVAL", "return {1,redis.error_reply('oops')}", 0),
		fail("EVAL", "return {err = 'oops', noerr = true}", 0), // doc error?
		fail("EVAL", "return {1, 2, err = 'oops'}", 0),         // doc error?

		succ("EVAL", "return {ok = 'great'}", 0),
		succ("EVAL", "return {1,{ok = 'great'}}", 0),
		succ("EVAL", "return redis.status_reply('great')", 0),
		succ("EVAL", "return {1,redis.status_reply('great')}", 0),
		succ("EVAL", "return {ok = 'great', notok = 'yes'}", 0),       // doc error?
		succ("EVAL", "return {1, 2, ok = 'great', notok = 'yes'}", 0), // doc error?

		failLoosely("EVAL", "return redis.error_reply(1)", 0),
		failLoosely("EVAL", "return redis.error_reply()", 0),
		failLoosely("EVAL", "return redis.error_reply(redis.error_reply('foo'))", 0),
		failLoosely("EVAL", "return redis.status_reply(1)", 0),
		failLoosely("EVAL", "return redis.status_reply()", 0),
		failLoosely("EVAL", "return redis.status_reply(redis.status_reply('foo'))", 0),
	)

	// state inside lua
	testCommands(t,
		succ("EVAL", "redis.call('SELECT', 3); redis.call('SET', 'foo', 'bar')", 0),
		succ("GET", "foo"),
		succ("SELECT", 3),
		succ("GET", "foo"),
	)

	// lua env
	testCommands(t,
		// succ("EVAL", "print(1)", 0),
		succ("EVAL", `return string.format('%q', "pretty string")`, 0),
		failLoosely("EVAL", "os.clock()", 0),
		failLoosely("EVAL", "os.exit(42)", 0),
		succ("EVAL", "return table.concat({1,2,3})", 0),
		succ("EVAL", "return math.abs(-42)", 0),
		failLoosely("EVAL", `return utf8.len("hello world")`, 0),
		failLoosely("EVAL", `require("utf8")`, 0),
		succ("EVAL", `return coroutine.running()`, 0),
	)

	// sha1hex
	testCommands(t,
		succ("EVAL", `return redis.sha1hex("foo")`, 0),
		succ("SET", "bar", "32"),
		succ("EVAL", `return redis.sha1hex(KEYS["bar"])`, 0),
		succ("EVAL", `return redis.sha1hex(KEYS[1])`, 1, "bar"),
		succ("EVAL", `return redis.sha1hex(nil)`, 0),
		succ("EVAL", `return redis.sha1hex(42)`, 0),
		succ("EVAL", `return redis.sha1hex({})`, 0),
		succ("EVAL", `return redis.sha1hex(KEYS[1])`, 0),
		failWith(
			"wrong number of arguments",
			"EVAL", `return redis.sha1hex()`, 0,
		),
		failWith(
			"wrong number of arguments",
			"EVAL", `return redis.sha1hex(1, 2)`, 0,
		),
	)

	// cjson module
	testCommands(t,
		succ("EVAL", `return cjson.decode('{"id":"foo"}')['id']`, 0),
		// succ("SET", "foo", `{"value":42}`),
		// succ("EVAL", `return KEYS[1]`, 1, "foo"),
		// succ("EVAL", `return cjson.decode(KEYS[1])['value']`, 1, "foo"),
		succ("EVAL", `return cjson.decode(ARGV[1])['value']`, 0, `{"value":"42"}`),
		succ("EVAL", `return redis.call("SET", "enc", cjson.encode({["foo"]="bar"}))`, 0),
		succ("EVAL", `return redis.call("SET", "enc", cjson.encode({["foo"]={["foo"]=42}}))`, 0),
		succ("GET", "enc"),

		failWith(
			"bad argument #1 to ",
			"EVAL", `return cjson.encode()`, 0,
		),
		failWith(
			"bad argument #1 to ",
			"EVAL", `return cjson.encode(1, 2)`, 0,
		),
		failWith(
			"bad argument #1 to ",
			"EVAL", `return cjson.decode()`, 0,
		),
		failWith(
			"bad argument #1 to ",
			"EVAL", `return cjson.decode(1, 2)`, 0,
		),
	)
}

func TestLuaCall(t *testing.T) {
	testCommands(t,
		succ("SET", "foo", 1),
		succ("EVAL", `local foo = redis.call("GET", "foo"); redis.call("SET", "foo", foo+1)`, 0),
		succ("GET", "foo"),
		succ("EVAL", `return redis.call("GET", "foo")`, 0),
		succ("EVAL", `return redis.call("SET", "foo", 42)`, 0),
	)

	// datatype errors
	testCommands(t,
		failWith(
			"Please specify at least one argument for redis.call()",
			"EVAL", `redis.call()`, 0,
		),
		failWith(
			"Lua redis() command arguments must be strings or integers",
			"EVAL", `redis.call({})`, 0,
		),
		failWith(
			"Unknown Redis command called from Lua script",
			"EVAL", `redis.call(1)`, 0,
		),
		failWith(
			"Lua redis() command arguments must be strings or integers",
			"EVAL", `redis.call("ECHO", true)`, 0,
		),
		failWith(
			"Lua redis() command arguments must be strings or integers",
			"EVAL", `redis.call("ECHO", false)`, 0,
		),
		failWith(
			"Lua redis() command arguments must be strings or integers",
			"EVAL", `redis.call("ECHO", nil)`, 0,
		),
		failWith(
			"Lua redis() command arguments must be strings or integers",
			"EVAL", `redis.call("HELLO", {})`, 0,
		),
		failLoosely("EVAL", `redis.call("HELLO", 1)`, 0),
		failLoosely("EVAL", `redis.call("HELLO", 3.14)`, 0),
		failWith(
			"Lua redis() command arguments must be strings or integers",
			"EVAL", `redis.call("GET", {})`, 0,
		),
	)

	// call() errors
	testCommands(t,
		succ("SET", "foo", 1),

		failLoosely("EVAL", `redis.call("HGET", "foo")`, 0),
		succ("GET", "foo"),
		failLoosely("EVAL", `local foo = redis.call("HGET", "foo"); redis.call("SET", "res", foo)`, 0),
		succ("GET", "foo"),
		succ("GET", "res"),
		failLoosely("EVAL", `local foo = redis.call("HGET", "foo", "bar"); redis.call("SET", "res", foo)`, 0),
		succ("GET", "foo"),
		succ("GET", "res"),
	)

	// pcall() errors
	testCommands(t,
		succ("SET", "foo", 1),
		failWith(
			"Lua redis() command arguments must be strings or integers",
			"EVAL", `local foo = redis.pcall("HGET", "foo"); redis.call("SET", "res", foo)`, 0,
		),
		succ("GET", "foo"),
		succ("GET", "res"),
		failWith(
			"Lua redis() command arguments must be strings or integers",
			"EVAL", `local foo = redis.pcall("HGET", "foo", "bar"); redis.call("SET", "res", foo)`, 0,
		),
		succ("GET", "foo"),
		succ("GET", "res"),
	)
}

func TestScriptNoAuth(t *testing.T) {
	testAuthCommands(t,
		"supersecret",
		failWith(
			"NOAUTH Authentication required.",
			"EVAL", `redis.call("ECHO", "foo")`, 0,
		),
		succ("AUTH", "supersecret"),
		succ(
			"EVAL", `redis.call("ECHO", "foo")`, 0,
		),
	)
}
