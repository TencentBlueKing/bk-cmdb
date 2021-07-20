package toml

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

type customString string

type stringer struct{}

func (s stringer) String() string {
	return "stringer"
}

func validate(t *testing.T, path string, object interface{}) {
	switch o := object.(type) {
	case *Tree:
		for key, tree := range o.values {
			validate(t, path+"."+key, tree)
		}
	case []*Tree:
		for index, tree := range o {
			validate(t, path+"."+strconv.Itoa(index), tree)
		}
	case *tomlValue:
		switch o.value.(type) {
		case int64, uint64, bool, string, float64, time.Time,
			[]int64, []uint64, []bool, []string, []float64, []time.Time:
		default:
			t.Fatalf("tomlValue at key %s containing incorrect type %T", path, o.value)
		}
	default:
		t.Fatalf("value at key %s is of incorrect type %T", path, object)
	}
	t.Logf("validation ok %s as %T", path, object)
}

func validateTree(t *testing.T, tree *Tree) {
	validate(t, "", tree)
}

func TestTreeCreateToTree(t *testing.T) {
	data := map[string]interface{}{
		"a_string": "bar",
		"an_int":   42,
		"time":     time.Now(),
		"int8":     int8(2),
		"int16":    int16(2),
		"int32":    int32(2),
		"uint8":    uint8(2),
		"uint16":   uint16(2),
		"uint32":   uint32(2),
		"float32":  float32(2),
		"a_bool":   false,
		"stringer": stringer{},
		"nested": map[string]interface{}{
			"foo": "bar",
		},
		"array":                 []string{"a", "b", "c"},
		"array_uint":            []uint{uint(1), uint(2)},
		"array_table":           []map[string]interface{}{{"sub_map": 52}},
		"array_times":           []time.Time{time.Now(), time.Now()},
		"map_times":             map[string]time.Time{"now": time.Now()},
		"custom_string_map_key": map[customString]interface{}{customString("custom"): "custom"},
	}
	tree, err := TreeFromMap(data)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	validateTree(t, tree)
}

func TestTreeCreateToTreeInvalidLeafType(t *testing.T) {
	_, err := TreeFromMap(map[string]interface{}{"foo": t})
	expected := "cannot convert type *testing.T to Tree"
	if err.Error() != expected {
		t.Fatalf("expected error %s, got %s", expected, err.Error())
	}
}

func TestTreeCreateToTreeInvalidMapKeyType(t *testing.T) {
	_, err := TreeFromMap(map[string]interface{}{"foo": map[int]interface{}{2: 1}})
	expected := "map key needs to be a string, not int (int)"
	if err.Error() != expected {
		t.Fatalf("expected error %s, got %s", expected, err.Error())
	}
}

func TestTreeCreateToTreeInvalidArrayMemberType(t *testing.T) {
	_, err := TreeFromMap(map[string]interface{}{"foo": []*testing.T{t}})
	expected := "cannot convert type *testing.T to Tree"
	if err.Error() != expected {
		t.Fatalf("expected error %s, got %s", expected, err.Error())
	}
}

func TestTreeCreateToTreeInvalidTableGroupType(t *testing.T) {
	_, err := TreeFromMap(map[string]interface{}{"foo": []map[string]interface{}{{"hello": t}}})
	expected := "cannot convert type *testing.T to Tree"
	if err.Error() != expected {
		t.Fatalf("expected error %s, got %s", expected, err.Error())
	}
}

func TestRoundTripArrayOfTables(t *testing.T) {
	orig := "\n[[stuff]]\n  name = \"foo\"\n  things = [\"a\", \"b\"]\n"
	tree, err := Load(orig)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	m := tree.ToMap()

	tree, err = TreeFromMap(m)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	want := orig
	got := tree.String()

	if got != want {
		t.Errorf("want:\n%s\ngot:\n%s", want, got)
	}
}

func TestTomlSliceOfSlice(t *testing.T) {
	tree, err := Load(` hosts=[["10.1.0.107:9092","10.1.0.107:9093", "192.168.0.40:9094"] ] `)
	m := tree.ToMap()
	tree, err = TreeFromMap(m)
	if err != nil {
		t.Error("should not error", err)
	}
	type Struct struct {
		Hosts [][]string
	}
	var actual Struct
	tree.Unmarshal(&actual)

	expected := Struct{Hosts: [][]string{[]string{"10.1.0.107:9092", "10.1.0.107:9093", "192.168.0.40:9094"}}}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Bad unmarshal: expected %+v, got %+v", expected, actual)
	}
}

func TestTomlSliceOfSliceOfSlice(t *testing.T) {
	tree, err := Load(` hosts=[[["10.1.0.107:9092","10.1.0.107:9093", "192.168.0.40:9094"] ]] `)
	m := tree.ToMap()
	tree, err = TreeFromMap(m)
	if err != nil {
		t.Error("should not error", err)
	}
	type Struct struct {
		Hosts [][][]string
	}
	var actual Struct
	tree.Unmarshal(&actual)

	expected := Struct{Hosts: [][][]string{[][]string{[]string{"10.1.0.107:9092", "10.1.0.107:9093", "192.168.0.40:9094"}}}}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Bad unmarshal: expected %+v, got %+v", expected, actual)
	}
}

func TestTomlSliceOfSliceInt(t *testing.T) {
	tree, err := Load(` hosts=[[1,2,3],[4,5,6] ] `)
	m := tree.ToMap()
	tree, err = TreeFromMap(m)
	if err != nil {
		t.Error("should not error", err)
	}
	type Struct struct {
		Hosts [][]int
	}
	var actual Struct
	err = tree.Unmarshal(&actual)
	if err != nil {
		t.Error("should not error", err)
	}

	expected := Struct{Hosts: [][]int{[]int{1, 2, 3}, []int{4, 5, 6}}}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Bad unmarshal: expected %+v, got %+v", expected, actual)
	}

}
func TestTomlSliceOfSliceInt64(t *testing.T) {
	tree, err := Load(` hosts=[[1,2,3],[4,5,6] ] `)
	m := tree.ToMap()
	tree, err = TreeFromMap(m)
	if err != nil {
		t.Error("should not error", err)
	}
	type Struct struct {
		Hosts [][]int64
	}
	var actual Struct
	err = tree.Unmarshal(&actual)
	if err != nil {
		t.Error("should not error", err)
	}

	expected := Struct{Hosts: [][]int64{[]int64{1, 2, 3}, []int64{4, 5, 6}}}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Bad unmarshal: expected %+v, got %+v", expected, actual)
	}

}

func TestTomlSliceOfSliceInt64FromMap(t *testing.T) {
	tree, err := TreeFromMap(map[string]interface{}{"hosts": [][]interface{}{[]interface{}{int32(1), int8(2), 3}}})
	if err != nil {
		t.Error("should not error", err)
	}
	type Struct struct {
		Hosts [][]int64
	}
	var actual Struct
	err = tree.Unmarshal(&actual)
	if err != nil {
		t.Error("should not error", err)
	}

	expected := Struct{Hosts: [][]int64{[]int64{1, 2, 3}}}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Bad unmarshal: expected %+v, got %+v", expected, actual)
	}

}
func TestTomlSliceOfSliceError(t *testing.T) { // make Codecov happy
	_, err := TreeFromMap(map[string]interface{}{"hosts": [][]interface{}{[]interface{}{1, 2, []struct{}{}}}})
	expected := "cannot convert type []struct {} to Tree"
	if err.Error() != expected {
		t.Fatalf("unexpected error: %s", err)
	}
}
