package common

import (
	"testing"
)

func TestStringPtrsValues(t *testing.T) {
	vals := []string{"a", "b", "c", "d"}
	ptrs := StringPtrs(vals)
	for i := 0; i < len(vals); i++ {
		if *ptrs[i] != vals[i] {
			t.Errorf("[ERROR] value %s != ptr value %s", vals[i], *ptrs[i])
		}
	}
	newVals := StringValues(ptrs)
	for i := 0; i < len(vals); i++ {
		if newVals[i] != vals[i] {
			t.Errorf("[ERROR] new val %s != val %s", newVals[i], vals[i])
		}
	}
}
