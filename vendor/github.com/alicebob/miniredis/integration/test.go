// +build int

package main

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/gomodule/redigo/redis"
)

type command struct {
	cmd      string // 'GET', 'SET', &c.
	args     []interface{}
	error    bool   // Whether the command should return an error or not.
	sort     bool   // Sort real redis's result. Used for 'keys'.
	loosely  bool   // Don't compare values, only structure. (for random things)
	errorSub string // Both errors need this substring
}

func succ(cmd string, args ...interface{}) command {
	return command{
		cmd:   cmd,
		args:  args,
		error: false,
	}
}

func succSorted(cmd string, args ...interface{}) command {
	return command{
		cmd:   cmd,
		args:  args,
		error: false,
		sort:  true,
	}
}

func succLoosely(cmd string, args ...interface{}) command {
	return command{
		cmd:     cmd,
		args:    args,
		error:   false,
		loosely: true,
	}
}

func fail(cmd string, args ...interface{}) command {
	return command{
		cmd:   cmd,
		args:  args,
		error: true,
	}
}

// expect an error, with `sub` in both errors
func failWith(sub string, cmd string, args ...interface{}) command {
	return command{
		cmd:      cmd,
		args:     args,
		error:    true,
		errorSub: sub,
	}
}

// only compare the error state, not the actual error message
func failLoosely(cmd string, args ...interface{}) command {
	return command{
		cmd:     cmd,
		args:    args,
		error:   true,
		loosely: true,
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("unexpected error: %s", err.Error())
	}
}

func testCommands(t *testing.T, commands ...command) {
	t.Helper()
	sMini, err := miniredis.Run()
	ok(t, err)
	defer sMini.Close()

	sReal, sRealAddr := Redis()
	defer sReal.Close()
	runCommands(t, sRealAddr, sMini.Addr(), commands)
}

// like testCommands, but multiple connections
func testMultiCommands(t *testing.T, cs ...func(chan<- command, *miniredis.Miniredis)) {
	t.Helper()
	sMini, err := miniredis.Run()
	ok(t, err)
	defer sMini.Close()

	sReal, realAddr := Redis()
	defer sReal.Close()

	var wg sync.WaitGroup
	for _, c := range cs {
		// one connections per cs
		cMini, err := redis.Dial("tcp", sMini.Addr())
		ok(t, err)

		cReal, err := redis.Dial("tcp", realAddr)
		ok(t, err)

		wg.Add(1)
		go func(c func(chan<- command, *miniredis.Miniredis)) {
			defer wg.Done()
			gen := make(chan command)
			wg.Add(1)
			go func() {
				defer wg.Done()
				c(gen, sMini)
				close(gen)
			}()
			for cm := range gen {
				runCommand(t, cMini, cReal, cm)
			}
		}(c)
	}
	wg.Wait()
}

func testAuthCommands(t *testing.T, passwd string, commands ...command) {
	sMini, err := miniredis.Run()
	ok(t, err)
	defer sMini.Close()
	sMini.RequireAuth(passwd)

	sReal, sRealAddr := RedisAuth(passwd)
	defer sReal.Close()
	runCommands(t, sRealAddr, sMini.Addr(), commands)
}

func runCommands(t *testing.T, realAddr, miniAddr string, commands []command) {
	t.Helper()
	cMini, err := redis.Dial("tcp", miniAddr)
	ok(t, err)

	cReal, err := redis.Dial("tcp", realAddr)
	ok(t, err)

	for _, c := range commands {
		runCommand(t, cMini, cReal, c)
	}
}

func runCommand(t *testing.T, cMini, cReal redis.Conn, p command) {
	t.Helper()
	vReal, errReal := cReal.Do(p.cmd, p.args...)
	vMini, errMini := cMini.Do(p.cmd, p.args...)
	if p.error {
		if errReal == nil {
			t.Errorf("got no error from realredis. case: %#v", p)
			return
		}
		if errMini == nil {
			t.Errorf("got no error from miniredis. case: %#v real error: %s", p, errReal)
			return
		}
		if p.loosely {
			return
		}
	} else {
		if errReal != nil {
			t.Errorf("got an error from realredis: %v. case: %#v", errReal, p)
			return
		}
		if errMini != nil {
			t.Errorf("got an error from miniredis: %v. case: %#v", errMini, p)
			return
		}
	}
	if p.errorSub != "" {
		if have, want := errReal.Error(), p.errorSub; !strings.Contains(have, want) {
			t.Errorf("realredis error error. expected: %q in %q case: %#v", want, have, p)
		}
		if have, want := errMini.Error(), p.errorSub; !strings.Contains(have, want) {
			t.Errorf("miniredis error error. expected: %q in %q case: %#v", want, have, p)
		}
		return
	}

	if !reflect.DeepEqual(errReal, errMini) {
		t.Errorf("error error. expected: %#v got: %#v case: %#v", vReal, vMini, p)
		return
	}
	// Sort the strings.
	if p.sort {
		sort.Sort(BytesList(vReal.([]interface{})))
		sort.Sort(BytesList(vMini.([]interface{})))
	}
	if p.loosely {
		if !looselyEqual(vReal, vMini) {
			t.Errorf("value error. expected: %#v got: %#v case: %#v", vReal, vMini, p)
			return
		}
	} else {
		if !reflect.DeepEqual(vReal, vMini) {
			t.Errorf("value error. expected: %#v got: %#v case: %#v", vReal, vMini, p)
			return
		}
	}
}

// BytesList implements the sort interface for things we know is a list of
// bytes.
type BytesList []interface{}

func (b BytesList) Len() int {
	return len(b)
}
func (b BytesList) Less(i, j int) bool {
	return bytes.Compare(b[i].([]byte), b[j].([]byte)) < 0
}
func (b BytesList) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func looselyEqual(a, b interface{}) bool {
	switch av := a.(type) {
	case string:
		_, ok := b.(string)
		return ok
	case []byte:
		_, ok := b.([]byte)
		return ok
	case []interface{}:
		bv, ok := b.([]interface{})
		if !ok {
			return false
		}
		if len(av) != len(bv) {
			return false
		}
		for i, v := range av {
			if !looselyEqual(v, bv[i]) {
				return false
			}
		}
		return true
	default:
		panic(fmt.Sprintf("unhandled case, got a %#v", a))
	}
}
