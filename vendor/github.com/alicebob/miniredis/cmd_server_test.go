package miniredis

import (
	"testing"

	"github.com/gomodule/redigo/redis"
)

// Test DBSIZE, FLUSHDB, and FLUSHALL.
func TestCmdServer(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c, err := redis.Dial("tcp", s.Addr())
	ok(t, err)

	// Set something
	{
		s.Set("aap", "niet")
		s.Set("roos", "vuur")
		s.DB(1).Set("noot", "mies")
	}

	{
		n, err := redis.Int(c.Do("DBSIZE"))
		ok(t, err)
		equals(t, 2, n)

		b, err := redis.String(c.Do("FLUSHDB"))
		ok(t, err)
		equals(t, "OK", b)

		n, err = redis.Int(c.Do("DBSIZE"))
		ok(t, err)
		equals(t, 0, n)

		_, err = c.Do("SELECT", 1)
		ok(t, err)

		n, err = redis.Int(c.Do("DBSIZE"))
		ok(t, err)
		equals(t, 1, n)

		b, err = redis.String(c.Do("FLUSHALL"))
		ok(t, err)
		equals(t, "OK", b)

		n, err = redis.Int(c.Do("DBSIZE"))
		ok(t, err)
		equals(t, 0, n)

		_, err = c.Do("SELECT", 4)
		ok(t, err)

		n, err = redis.Int(c.Do("DBSIZE"))
		ok(t, err)
		equals(t, 0, n)

	}

	{
		b, err := redis.String(c.Do("FLUSHDB", "ASYNC"))
		ok(t, err)
		equals(t, "OK", b)

		b, err = redis.String(c.Do("FLUSHALL", "ASYNC"))
		ok(t, err)
		equals(t, "OK", b)
	}

	{
		_, err := redis.Int(c.Do("DBSIZE", "FOO"))
		assert(t, err != nil, "no DBSIZE error")

		_, err = redis.Int(c.Do("FLUSHDB", "FOO"))
		assert(t, err != nil, "no FLUSHDB error")

		_, err = redis.Int(c.Do("FLUSHDB", "ASYNC", "FOO"))
		assert(t, err != nil, "no FLUSHDB error")

		_, err = redis.Int(c.Do("FLUSHALL", "FOO"))
		assert(t, err != nil, "no FLUSHALL error")

		_, err = redis.Int(c.Do("FLUSHALL", "ASYNC", "FOO"))
		assert(t, err != nil, "no FLUSHALL error")

		_, err = redis.Int(c.Do("FLUSHALL", "ASYNC", "ASYNC"))
		assert(t, err != nil, "no FLUSHALL error")
	}
}
