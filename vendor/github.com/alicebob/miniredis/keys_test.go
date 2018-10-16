package miniredis

import (
	"testing"
)

func TestKeysSel(t *testing.T) {
	// Helper to test the selection behind KEYS
	// pattern -> cases -> should match?
	test := func(pat string, chk map[string]bool) {
		t.Helper()
		patRe := patternRE(pat)
		if patRe == nil {
			t.Errorf("'%v' won't match anything. Didn't expect that.", pat)
			return
		}
		for key, expected := range chk {
			match := patRe.MatchString(key)
			if have, want := match, expected; have != want {
				t.Errorf("'%v' -> '%v'. have %v, want %v", pat, key, have, want)
			}
		}
	}
	test("aap", map[string]bool{
		"aap":         true,
		"aapnoot":     false,
		"nootaap":     false,
		"nootaapnoot": false,
		"AAP":         false,
	})
	test("aap*", map[string]bool{
		"aap":         true,
		"aapnoot":     true,
		"nootaap":     false,
		"nootaapnoot": false,
		"AAP":         false,
	})
	// No problem with regexp meta chars?
	test("(?:a)ap*", map[string]bool{
		"(?:a)ap!": true,
		"aap":      false,
	})
	test("*aap*", map[string]bool{
		"aap":         true,
		"aapnoot":     true,
		"nootaap":     true,
		"nootaapnoot": true,
		"AAP":         false,
		"a_a_p":       false,
	})
	test(`\*aap*`, map[string]bool{
		"*aap":     true,
		"aap":      false,
		"*aapnoot": true,
		"aapnoot":  false,
	})
	test(`aa?`, map[string]bool{
		"aap":  true,
		"aal":  true,
		"aaf":  true,
		"aa?":  true,
		"aap!": false,
	})
	test(`aa\?`, map[string]bool{
		"aap":  false,
		"aa?":  true,
		"aa?!": false,
	})
	test("aa[pl]", map[string]bool{
		"aap":  true,
		"aal":  true,
		"aaf":  false,
		"aa?":  false,
		"aap!": false,
	})
	test("[ab]a[pl]", map[string]bool{
		"aap":  true,
		"aal":  true,
		"bap":  true,
		"bal":  true,
		"aaf":  false,
		"cap":  false,
		"aa?":  false,
		"aap!": false,
	})
	test(`\[ab\]`, map[string]bool{
		"[ab]": true,
		"a":    false,
	})
	test(`[\[ab]`, map[string]bool{
		"[": true,
		"a": true,
		"b": true,
		"c": false,
		"]": false,
	})
	test(`[\[\]]`, map[string]bool{
		"[": true,
		"]": true,
		"c": false,
	})
	test(`\\ap`, map[string]bool{
		`\ap`:  true,
		`\\ap`: false,
	})
	// Escape a normal char
	test(`\foo`, map[string]bool{
		`foo`:  true,
		`\foo`: false,
	})

	// Patterns which won't match anything.
	test2 := func(pat string) {
		t.Helper()
		if patternRE(pat) != nil {
			t.Errorf("'%v' will match something. Didn't expect that.", pat)
		}
	}
	test2(`ap[\`) // trailing \ in char class
	test2(`ap[`)  // open char class
	test2(`[]ap`) // empty char class
	test2(`ap\`)  // trailing \
}
