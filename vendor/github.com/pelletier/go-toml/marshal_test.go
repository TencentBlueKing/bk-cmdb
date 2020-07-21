package toml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

type basicMarshalTestStruct struct {
	String     string   `toml:"Zstring"`
	StringList []string `toml:"Ystrlist"`
	BasicMarshalTestSubAnonymousStruct
	Sub     basicMarshalTestSubStruct   `toml:"Xsubdoc"`
	SubList []basicMarshalTestSubStruct `toml:"Wsublist"`
}

type basicMarshalTestSubStruct struct {
	String2 string
}

type BasicMarshalTestSubAnonymousStruct struct {
	String3 string
}

var basicTestData = basicMarshalTestStruct{
	String:                             "Hello",
	StringList:                         []string{"Howdy", "Hey There"},
	BasicMarshalTestSubAnonymousStruct: BasicMarshalTestSubAnonymousStruct{"One"},
	Sub:                                basicMarshalTestSubStruct{"Two"},
	SubList:                            []basicMarshalTestSubStruct{{"Three"}, {"Four"}},
}

var basicTestToml = []byte(`String3 = "One"
Ystrlist = ["Howdy", "Hey There"]
Zstring = "Hello"

[[Wsublist]]
  String2 = "Three"

[[Wsublist]]
  String2 = "Four"

[Xsubdoc]
  String2 = "Two"
`)

var basicTestTomlCustomIndentation = []byte(`String3 = "One"
Ystrlist = ["Howdy", "Hey There"]
Zstring = "Hello"

[[Wsublist]]
	String2 = "Three"

[[Wsublist]]
	String2 = "Four"

[Xsubdoc]
	String2 = "Two"
`)

var basicTestTomlOrdered = []byte(`Zstring = "Hello"
Ystrlist = ["Howdy", "Hey There"]
String3 = "One"

[Xsubdoc]
  String2 = "Two"

[[Wsublist]]
  String2 = "Three"

[[Wsublist]]
  String2 = "Four"
`)

var marshalTestToml = []byte(`title = "TOML Marshal Testing"

[basic]
  bool = true
  date = 1979-05-27T07:32:00Z
  float = 123.4
  float64 = 123.456782132399
  int = 5000
  string = "Bite me"
  uint = 5001

[basic_lists]
  bools = [true, false, true]
  dates = [1979-05-27T07:32:00Z, 1980-05-27T07:32:00Z]
  floats = [12.3, 45.6, 78.9]
  ints = [8001, 8001, 8002]
  strings = ["One", "Two", "Three"]
  uints = [5002, 5003]

[basic_map]
  one = "one"
  two = "two"

[subdoc]

  [subdoc.first]
    name = "First"

  [subdoc.second]
    name = "Second"

[[subdoclist]]
  name = "List.First"

[[subdoclist]]
  name = "List.Second"

[[subdocptrs]]
  name = "Second"
`)

var marshalOrderPreserveToml = []byte(`title = "TOML Marshal Testing"

[basic_lists]
  floats = [12.3, 45.6, 78.9]
  bools = [true, false, true]
  dates = [1979-05-27T07:32:00Z, 1980-05-27T07:32:00Z]
  ints = [8001, 8001, 8002]
  uints = [5002, 5003]
  strings = ["One", "Two", "Three"]

[[subdocptrs]]
  name = "Second"

[basic_map]
  one = "one"
  two = "two"

[subdoc]

  [subdoc.second]
    name = "Second"

  [subdoc.first]
    name = "First"

[basic]
  uint = 5001
  bool = true
  float = 123.4
  float64 = 123.456782132399
  int = 5000
  string = "Bite me"
  date = 1979-05-27T07:32:00Z

[[subdoclist]]
  name = "List.First"

[[subdoclist]]
  name = "List.Second"
`)

var mashalOrderPreserveMapToml = []byte(`title = "TOML Marshal Testing"

[basic_map]
  one = "one"
  two = "two"

[long_map]
  a7 = "1"
  b3 = "2"
  c8 = "3"
  d4 = "4"
  e6 = "5"
  f5 = "6"
  g10 = "7"
  h1 = "8"
  i2 = "9"
  j9 = "10"
`)

type Conf struct {
	Name  string
	Age   int
	Inter interface{}
}

type NestedStruct struct {
	FirstName string
	LastName  string
	Age       int
}

var doc = []byte(`Name = "rui"
Age = 18

[Inter]
  FirstName = "wang"
  LastName = "jl"
  Age = 100`)

func TestInterface(t *testing.T) {
	var config Conf
	config.Inter = &NestedStruct{}
	err := Unmarshal(doc, &config)
	expected := Conf{
		Name: "rui",
		Age:  18,
		Inter: &NestedStruct{
			FirstName: "wang",
			LastName:  "jl",
			Age:       100,
		},
	}
	if err != nil || !reflect.DeepEqual(config, expected) {
		t.Errorf("Bad unmarshal: expected %v, got %v", expected, config)
	}
}

func TestBasicMarshal(t *testing.T) {
	result, err := Marshal(basicTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := basicTestToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestBasicMarshalCustomIndentation(t *testing.T) {
	var result bytes.Buffer
	err := NewEncoder(&result).Indentation("\t").Encode(basicTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := basicTestTomlCustomIndentation
	if !bytes.Equal(result.Bytes(), expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result.Bytes())
	}
}

func TestBasicMarshalWrongIndentation(t *testing.T) {
	var result bytes.Buffer
	err := NewEncoder(&result).Indentation("  \n").Encode(basicTestData)
	if err.Error() != "invalid indentation: must only contains space or tab characters" {
		t.Error("expect err:invalid indentation: must only contains space or tab characters but got:", err)
	}
}

func TestBasicMarshalOrdered(t *testing.T) {
	var result bytes.Buffer
	err := NewEncoder(&result).Order(OrderPreserve).Encode(basicTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := basicTestTomlOrdered
	if !bytes.Equal(result.Bytes(), expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result.Bytes())
	}
}

func TestBasicMarshalWithPointer(t *testing.T) {
	result, err := Marshal(&basicTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := basicTestToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestBasicMarshalOrderedWithPointer(t *testing.T) {
	var result bytes.Buffer
	err := NewEncoder(&result).Order(OrderPreserve).Encode(&basicTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := basicTestTomlOrdered
	if !bytes.Equal(result.Bytes(), expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result.Bytes())
	}
}

func TestBasicUnmarshal(t *testing.T) {
	result := basicMarshalTestStruct{}
	err := Unmarshal(basicTestToml, &result)
	expected := basicTestData
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad unmarshal: expected %v, got %v", expected, result)
	}
}

type quotedKeyMarshalTestStruct struct {
	String  string                      `toml:"Z.string-àéù"`
	Float   float64                     `toml:"Yfloat-𝟘"`
	Sub     basicMarshalTestSubStruct   `toml:"Xsubdoc-àéù"`
	SubList []basicMarshalTestSubStruct `toml:"W.sublist-𝟘"`
}

var quotedKeyMarshalTestData = quotedKeyMarshalTestStruct{
	String:  "Hello",
	Float:   3.5,
	Sub:     basicMarshalTestSubStruct{"One"},
	SubList: []basicMarshalTestSubStruct{{"Two"}, {"Three"}},
}

var quotedKeyMarshalTestToml = []byte(`"Yfloat-𝟘" = 3.5
"Z.string-àéù" = "Hello"

[["W.sublist-𝟘"]]
  String2 = "Two"

[["W.sublist-𝟘"]]
  String2 = "Three"

["Xsubdoc-àéù"]
  String2 = "One"
`)

func TestBasicMarshalQuotedKey(t *testing.T) {
	result, err := Marshal(quotedKeyMarshalTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := quotedKeyMarshalTestToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestBasicUnmarshalQuotedKey(t *testing.T) {
	tree, err := LoadBytes(quotedKeyMarshalTestToml)
	if err != nil {
		t.Fatal(err)
	}

	var q quotedKeyMarshalTestStruct
	tree.Unmarshal(&q)
	fmt.Println(q)

	if !reflect.DeepEqual(quotedKeyMarshalTestData, q) {
		t.Errorf("Bad unmarshal: expected\n-----\n%v\n-----\ngot\n-----\n%v\n-----\n", quotedKeyMarshalTestData, q)
	}
}

type testDoc struct {
	Title       string            `toml:"title"`
	BasicLists  testDocBasicLists `toml:"basic_lists"`
	SubDocPtrs  []*testSubDoc     `toml:"subdocptrs"`
	BasicMap    map[string]string `toml:"basic_map"`
	Subdocs     testDocSubs       `toml:"subdoc"`
	Basics      testDocBasics     `toml:"basic"`
	SubDocList  []testSubDoc      `toml:"subdoclist"`
	err         int               `toml:"shouldntBeHere"`
	unexported  int               `toml:"shouldntBeHere"`
	Unexported2 int               `toml:"-"`
}

type testMapDoc struct {
	Title    string            `toml:"title"`
	BasicMap map[string]string `toml:"basic_map"`
	LongMap  map[string]string `toml:"long_map"`
}

type testDocBasics struct {
	Uint       uint      `toml:"uint"`
	Bool       bool      `toml:"bool"`
	Float32    float32   `toml:"float"`
	Float64    float64   `toml:"float64"`
	Int        int       `toml:"int"`
	String     *string   `toml:"string"`
	Date       time.Time `toml:"date"`
	unexported int       `toml:"shouldntBeHere"`
}

type testDocBasicLists struct {
	Floats  []*float32  `toml:"floats"`
	Bools   []bool      `toml:"bools"`
	Dates   []time.Time `toml:"dates"`
	Ints    []int       `toml:"ints"`
	UInts   []uint      `toml:"uints"`
	Strings []string    `toml:"strings"`
}

type testDocSubs struct {
	Second *testSubDoc `toml:"second"`
	First  testSubDoc  `toml:"first"`
}

type testSubDoc struct {
	Name       string `toml:"name"`
	unexported int    `toml:"shouldntBeHere"`
}

var biteMe = "Bite me"
var float1 float32 = 12.3
var float2 float32 = 45.6
var float3 float32 = 78.9
var subdoc = testSubDoc{"Second", 0}

var docData = testDoc{
	Title:       "TOML Marshal Testing",
	unexported:  0,
	Unexported2: 0,
	Basics: testDocBasics{
		Bool:       true,
		Date:       time.Date(1979, 5, 27, 7, 32, 0, 0, time.UTC),
		Float32:    123.4,
		Float64:    123.456782132399,
		Int:        5000,
		Uint:       5001,
		String:     &biteMe,
		unexported: 0,
	},
	BasicLists: testDocBasicLists{
		Bools: []bool{true, false, true},
		Dates: []time.Time{
			time.Date(1979, 5, 27, 7, 32, 0, 0, time.UTC),
			time.Date(1980, 5, 27, 7, 32, 0, 0, time.UTC),
		},
		Floats:  []*float32{&float1, &float2, &float3},
		Ints:    []int{8001, 8001, 8002},
		Strings: []string{"One", "Two", "Three"},
		UInts:   []uint{5002, 5003},
	},
	BasicMap: map[string]string{
		"one": "one",
		"two": "two",
	},
	Subdocs: testDocSubs{
		First:  testSubDoc{"First", 0},
		Second: &subdoc,
	},
	SubDocList: []testSubDoc{
		{"List.First", 0},
		{"List.Second", 0},
	},
	SubDocPtrs: []*testSubDoc{&subdoc},
}

var mapTestDoc = testMapDoc{
	Title: "TOML Marshal Testing",
	BasicMap: map[string]string{
		"one": "one",
		"two": "two",
	},
	LongMap: map[string]string{
		"h1":  "8",
		"i2":  "9",
		"b3":  "2",
		"d4":  "4",
		"f5":  "6",
		"e6":  "5",
		"a7":  "1",
		"c8":  "3",
		"j9":  "10",
		"g10": "7",
	},
}

func TestDocMarshal(t *testing.T) {
	result, err := Marshal(docData)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(result, marshalTestToml) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", marshalTestToml, result)
	}
}

func TestDocMarshalOrdered(t *testing.T) {
	var result bytes.Buffer
	err := NewEncoder(&result).Order(OrderPreserve).Encode(docData)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(result.Bytes(), marshalOrderPreserveToml) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", marshalOrderPreserveToml, result.Bytes())
	}
}

func TestDocMarshalMaps(t *testing.T) {
	result, err := Marshal(mapTestDoc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(result, mashalOrderPreserveMapToml) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", mashalOrderPreserveMapToml, result)
	}
}

func TestDocMarshalOrderedMaps(t *testing.T) {
	var result bytes.Buffer
	err := NewEncoder(&result).Order(OrderPreserve).Encode(mapTestDoc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(result.Bytes(), mashalOrderPreserveMapToml) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", mashalOrderPreserveMapToml, result.Bytes())
	}
}

func TestDocMarshalPointer(t *testing.T) {
	result, err := Marshal(&docData)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, marshalTestToml) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", marshalTestToml, result)
	}
}

func TestDocUnmarshal(t *testing.T) {
	result := testDoc{}
	err := Unmarshal(marshalTestToml, &result)
	expected := docData
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		resStr, _ := json.MarshalIndent(result, "", "  ")
		expStr, _ := json.MarshalIndent(expected, "", "  ")
		t.Errorf("Bad unmarshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expStr, resStr)
	}
}

func TestDocPartialUnmarshal(t *testing.T) {
	file, err := ioutil.TempFile("", "test-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	err = ioutil.WriteFile(file.Name(), marshalTestToml, 0)
	if err != nil {
		t.Fatal(err)
	}

	tree, _ := LoadFile(file.Name())
	subTree := tree.Get("subdoc").(*Tree)

	result := testDocSubs{}
	err = subTree.Unmarshal(&result)
	expected := docData.Subdocs
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		resStr, _ := json.MarshalIndent(result, "", "  ")
		expStr, _ := json.MarshalIndent(expected, "", "  ")
		t.Errorf("Bad partial unmartial: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expStr, resStr)
	}
}

type tomlTypeCheckTest struct {
	name string
	item interface{}
	typ  int //0=primitive, 1=otherslice, 2=treeslice, 3=tree
}

func TestTypeChecks(t *testing.T) {
	tests := []tomlTypeCheckTest{
		{"bool", true, 0},
		{"bool", false, 0},
		{"int", int(2), 0},
		{"int8", int8(2), 0},
		{"int16", int16(2), 0},
		{"int32", int32(2), 0},
		{"int64", int64(2), 0},
		{"uint", uint(2), 0},
		{"uint8", uint8(2), 0},
		{"uint16", uint16(2), 0},
		{"uint32", uint32(2), 0},
		{"uint64", uint64(2), 0},
		{"float32", float32(3.14), 0},
		{"float64", float64(3.14), 0},
		{"string", "lorem ipsum", 0},
		{"time", time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), 0},
		{"stringlist", []string{"hello", "hi"}, 1},
		{"stringlistptr", &[]string{"hello", "hi"}, 1},
		{"stringarray", [2]string{"hello", "hi"}, 1},
		{"stringarrayptr", &[2]string{"hello", "hi"}, 1},
		{"timelist", []time.Time{time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)}, 1},
		{"timelistptr", &[]time.Time{time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)}, 1},
		{"timearray", [1]time.Time{time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)}, 1},
		{"timearrayptr", &[1]time.Time{time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)}, 1},
		{"objectlist", []tomlTypeCheckTest{}, 2},
		{"objectlistptr", &[]tomlTypeCheckTest{}, 2},
		{"objectarray", [2]tomlTypeCheckTest{{}, {}}, 2},
		{"objectlistptr", &[2]tomlTypeCheckTest{{}, {}}, 2},
		{"object", tomlTypeCheckTest{}, 3},
		{"objectptr", &tomlTypeCheckTest{}, 3},
	}

	for _, test := range tests {
		expected := []bool{false, false, false, false}
		expected[test.typ] = true
		result := []bool{
			isPrimitive(reflect.TypeOf(test.item)),
			isOtherSequence(reflect.TypeOf(test.item)),
			isTreeSequence(reflect.TypeOf(test.item)),
			isTree(reflect.TypeOf(test.item)),
		}
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("Bad type check on %q: expected %v, got %v", test.name, expected, result)
		}
	}
}

type unexportedMarshalTestStruct struct {
	String      string                      `toml:"string"`
	StringList  []string                    `toml:"strlist"`
	Sub         basicMarshalTestSubStruct   `toml:"subdoc"`
	SubList     []basicMarshalTestSubStruct `toml:"sublist"`
	unexported  int                         `toml:"shouldntBeHere"`
	Unexported2 int                         `toml:"-"`
}

var unexportedTestData = unexportedMarshalTestStruct{
	String:      "Hello",
	StringList:  []string{"Howdy", "Hey There"},
	Sub:         basicMarshalTestSubStruct{"One"},
	SubList:     []basicMarshalTestSubStruct{{"Two"}, {"Three"}},
	unexported:  0,
	Unexported2: 0,
}

var unexportedTestToml = []byte(`string = "Hello"
strlist = ["Howdy","Hey There"]
unexported = 1
shouldntBeHere = 2

[subdoc]
  String2 = "One"

[[sublist]]
  String2 = "Two"

[[sublist]]
  String2 = "Three"
`)

func TestUnexportedUnmarshal(t *testing.T) {
	result := unexportedMarshalTestStruct{}
	err := Unmarshal(unexportedTestToml, &result)
	expected := unexportedTestData
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad unexported unmarshal: expected %v, got %v", expected, result)
	}
}

type errStruct struct {
	Bool   bool      `toml:"bool"`
	Date   time.Time `toml:"date"`
	Float  float64   `toml:"float"`
	Int    int16     `toml:"int"`
	String *string   `toml:"string"`
}

var errTomls = []string{
	"bool = truly\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:3200Z\nfloat = 123.4\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123a4\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = j000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = 5000\nstring = Bite me",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = 5000\nstring = Bite me",
	"bool = 1\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1\nfloat = 123.4\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\n\"sorry\"\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = \"sorry\"\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = 5000\nstring = 1",
}

type mapErr struct {
	Vals map[string]float64
}

type intErr struct {
	Int1  int
	Int2  int8
	Int3  int16
	Int4  int32
	Int5  int64
	UInt1 uint
	UInt2 uint8
	UInt3 uint16
	UInt4 uint32
	UInt5 uint64
	Flt1  float32
	Flt2  float64
}

var intErrTomls = []string{
	"Int1 = []\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = []\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = []\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = []\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = []\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = []\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = []\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = []\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = []\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = []\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = []\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = []",
}

func TestErrUnmarshal(t *testing.T) {
	for ind, toml := range errTomls {
		result := errStruct{}
		err := Unmarshal([]byte(toml), &result)
		if err == nil {
			t.Errorf("Expected err from case %d\n", ind)
		}
	}
	result2 := mapErr{}
	err := Unmarshal([]byte("[Vals]\nfred=\"1.2\""), &result2)
	if err == nil {
		t.Errorf("Expected err from map")
	}
	for ind, toml := range intErrTomls {
		result3 := intErr{}
		err := Unmarshal([]byte(toml), &result3)
		if err == nil {
			t.Errorf("Expected int err from case %d\n", ind)
		}
	}
}

type emptyMarshalTestStruct struct {
	Title      string                  `toml:"title"`
	Bool       bool                    `toml:"bool"`
	Int        int                     `toml:"int"`
	String     string                  `toml:"string"`
	StringList []string                `toml:"stringlist"`
	Ptr        *basicMarshalTestStruct `toml:"ptr"`
	Map        map[string]string       `toml:"map"`
}

var emptyTestData = emptyMarshalTestStruct{
	Title:      "Placeholder",
	Bool:       false,
	Int:        0,
	String:     "",
	StringList: []string{},
	Ptr:        nil,
	Map:        map[string]string{},
}

var emptyTestToml = []byte(`bool = false
int = 0
string = ""
stringlist = []
title = "Placeholder"

[map]
`)

type emptyMarshalTestStruct2 struct {
	Title      string                  `toml:"title"`
	Bool       bool                    `toml:"bool,omitempty"`
	Int        int                     `toml:"int, omitempty"`
	String     string                  `toml:"string,omitempty "`
	StringList []string                `toml:"stringlist,omitempty"`
	Ptr        *basicMarshalTestStruct `toml:"ptr,omitempty"`
	Map        map[string]string       `toml:"map,omitempty"`
}

var emptyTestData2 = emptyMarshalTestStruct2{
	Title:      "Placeholder",
	Bool:       false,
	Int:        0,
	String:     "",
	StringList: []string{},
	Ptr:        nil,
	Map:        map[string]string{},
}

var emptyTestToml2 = []byte(`title = "Placeholder"
`)

func TestEmptyMarshal(t *testing.T) {
	result, err := Marshal(emptyTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := emptyTestToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad empty marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestEmptyMarshalOmit(t *testing.T) {
	result, err := Marshal(emptyTestData2)
	if err != nil {
		t.Fatal(err)
	}
	expected := emptyTestToml2
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad empty omit marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestEmptyUnmarshal(t *testing.T) {
	result := emptyMarshalTestStruct{}
	err := Unmarshal(emptyTestToml, &result)
	expected := emptyTestData
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad empty unmarshal: expected %v, got %v", expected, result)
	}
}

func TestEmptyUnmarshalOmit(t *testing.T) {
	result := emptyMarshalTestStruct2{}
	err := Unmarshal(emptyTestToml, &result)
	expected := emptyTestData2
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad empty omit unmarshal: expected %v, got %v", expected, result)
	}
}

type pointerMarshalTestStruct struct {
	Str       *string
	List      *[]string
	ListPtr   *[]*string
	Map       *map[string]string
	MapPtr    *map[string]*string
	EmptyStr  *string
	EmptyList *[]string
	EmptyMap  *map[string]string
	DblPtr    *[]*[]*string
}

var pointerStr = "Hello"
var pointerList = []string{"Hello back"}
var pointerListPtr = []*string{&pointerStr}
var pointerMap = map[string]string{"response": "Goodbye"}
var pointerMapPtr = map[string]*string{"alternate": &pointerStr}
var pointerTestData = pointerMarshalTestStruct{
	Str:       &pointerStr,
	List:      &pointerList,
	ListPtr:   &pointerListPtr,
	Map:       &pointerMap,
	MapPtr:    &pointerMapPtr,
	EmptyStr:  nil,
	EmptyList: nil,
	EmptyMap:  nil,
}

var pointerTestToml = []byte(`List = ["Hello back"]
ListPtr = ["Hello"]
Str = "Hello"

[Map]
  response = "Goodbye"

[MapPtr]
  alternate = "Hello"
`)

func TestPointerMarshal(t *testing.T) {
	result, err := Marshal(pointerTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := pointerTestToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad pointer marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestPointerUnmarshal(t *testing.T) {
	result := pointerMarshalTestStruct{}
	err := Unmarshal(pointerTestToml, &result)
	expected := pointerTestData
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad pointer unmarshal: expected %v, got %v", expected, result)
	}
}

func TestUnmarshalTypeMismatch(t *testing.T) {
	result := pointerMarshalTestStruct{}
	err := Unmarshal([]byte("List = 123"), &result)
	if !strings.HasPrefix(err.Error(), "(1, 1): Can't convert 123(int64) to []string(slice)") {
		t.Errorf("Type mismatch must be reported: got %v", err.Error())
	}
}

type nestedMarshalTestStruct struct {
	String [][]string
	//Struct [][]basicMarshalTestSubStruct
	StringPtr *[]*[]*string
	// StructPtr *[]*[]*basicMarshalTestSubStruct
}

var str1 = "Three"
var str2 = "Four"
var strPtr = []*string{&str1, &str2}
var strPtr2 = []*[]*string{&strPtr}

var nestedTestData = nestedMarshalTestStruct{
	String:    [][]string{{"Five", "Six"}, {"One", "Two"}},
	StringPtr: &strPtr2,
}

var nestedTestToml = []byte(`String = [["Five", "Six"], ["One", "Two"]]
StringPtr = [["Three", "Four"]]
`)

func TestNestedMarshal(t *testing.T) {
	result, err := Marshal(nestedTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := nestedTestToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad nested marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestNestedUnmarshal(t *testing.T) {
	result := nestedMarshalTestStruct{}
	err := Unmarshal(nestedTestToml, &result)
	expected := nestedTestData
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad nested unmarshal: expected %v, got %v", expected, result)
	}
}

type customMarshalerParent struct {
	Self    customMarshaler   `toml:"me"`
	Friends []customMarshaler `toml:"friends"`
}

type customMarshaler struct {
	FirstName string
	LastName  string
}

func (c customMarshaler) MarshalTOML() ([]byte, error) {
	fullName := fmt.Sprintf("%s %s", c.FirstName, c.LastName)
	return []byte(fullName), nil
}

var customMarshalerData = customMarshaler{FirstName: "Sally", LastName: "Fields"}
var customMarshalerToml = []byte(`Sally Fields`)
var nestedCustomMarshalerData = customMarshalerParent{
	Self:    customMarshaler{FirstName: "Maiku", LastName: "Suteda"},
	Friends: []customMarshaler{customMarshalerData},
}
var nestedCustomMarshalerToml = []byte(`friends = ["Sally Fields"]
me = "Maiku Suteda"
`)
var nestedCustomMarshalerTomlForUnmarshal = []byte(`[friends]
FirstName = "Sally"
LastName = "Fields"`)

func TestCustomMarshaler(t *testing.T) {
	result, err := Marshal(customMarshalerData)
	if err != nil {
		t.Fatal(err)
	}
	expected := customMarshalerToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad custom marshaler: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

type IntOrString string

func (x *IntOrString) MarshalTOML() ([]byte, error) {
	s := *(*string)(x)
	_, err := strconv.Atoi(s)
	if err != nil {
		return []byte(fmt.Sprintf(`"%s"`, s)), nil
	}
	return []byte(s), nil
}

func TestNestedCustomMarshaler(t *testing.T) {
	num := IntOrString("100")
	str := IntOrString("hello")
	var parent = struct {
		IntField    *IntOrString `toml:"int"`
		StringField *IntOrString `toml:"string"`
	}{
		&num,
		&str,
	}

	result, err := Marshal(parent)
	if err != nil {
		t.Fatal(err)
	}
	expected := `int = 100
string = "hello"
`
	if !bytes.Equal(result, []byte(expected)) {
		t.Errorf("Bad nested text marshaler: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

type textMarshaler struct {
	FirstName string
	LastName  string
}

func (m textMarshaler) MarshalText() ([]byte, error) {
	fullName := fmt.Sprintf("%s %s", m.FirstName, m.LastName)
	return []byte(fullName), nil
}

func TestTextMarshaler(t *testing.T) {
	m := textMarshaler{FirstName: "Sally", LastName: "Fields"}

	result, err := Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	expected := `Sally Fields`
	if !bytes.Equal(result, []byte(expected)) {
		t.Errorf("Bad text marshaler: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestUnmarshalTextMarshaler(t *testing.T) {
	var nested = struct {
		Friends textMarshaler `toml:"friends"`
	}{}

	var expected = struct {
		Friends textMarshaler `toml:"friends"`
	}{
		Friends: textMarshaler{FirstName: "Sally", LastName: "Fields"},
	}

	err := Unmarshal(nestedCustomMarshalerTomlForUnmarshal, &nested)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(nested, expected) {
		t.Errorf("Bad unmarshal: expected %v, got %v", expected, nested)
	}
}

func TestNestedTextMarshaler(t *testing.T) {
	var parent = struct {
		Self     textMarshaler   `toml:"me"`
		Friends  []textMarshaler `toml:"friends"`
		Stranger *textMarshaler  `toml:"stranger"`
	}{
		Self:     textMarshaler{FirstName: "Maiku", LastName: "Suteda"},
		Friends:  []textMarshaler{textMarshaler{FirstName: "Sally", LastName: "Fields"}},
		Stranger: &textMarshaler{FirstName: "Earl", LastName: "Henson"},
	}

	result, err := Marshal(parent)
	if err != nil {
		t.Fatal(err)
	}
	expected := `friends = ["Sally Fields"]
me = "Maiku Suteda"
stranger = "Earl Henson"
`
	if !bytes.Equal(result, []byte(expected)) {
		t.Errorf("Bad nested text marshaler: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

type precedentMarshaler struct {
	FirstName string
	LastName  string
}

func (m precedentMarshaler) MarshalText() ([]byte, error) {
	return []byte("shadowed"), nil
}

func (m precedentMarshaler) MarshalTOML() ([]byte, error) {
	fullName := fmt.Sprintf("%s %s", m.FirstName, m.LastName)
	return []byte(fullName), nil
}

func TestPrecedentMarshaler(t *testing.T) {
	m := textMarshaler{FirstName: "Sally", LastName: "Fields"}

	result, err := Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	expected := `Sally Fields`
	if !bytes.Equal(result, []byte(expected)) {
		t.Errorf("Bad text marshaler: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

type customPointerMarshaler struct {
	FirstName string
	LastName  string
}

func (m *customPointerMarshaler) MarshalTOML() ([]byte, error) {
	return []byte(`"hidden"`), nil
}

type textPointerMarshaler struct {
	FirstName string
	LastName  string
}

func (m *textPointerMarshaler) MarshalText() ([]byte, error) {
	return []byte("hidden"), nil
}

func TestPointerMarshaler(t *testing.T) {
	var parent = struct {
		Self     customPointerMarshaler  `toml:"me"`
		Stranger *customPointerMarshaler `toml:"stranger"`
		Friend   textPointerMarshaler    `toml:"friend"`
		Fiend    *textPointerMarshaler   `toml:"fiend"`
	}{
		Self:     customPointerMarshaler{FirstName: "Maiku", LastName: "Suteda"},
		Stranger: &customPointerMarshaler{FirstName: "Earl", LastName: "Henson"},
		Friend:   textPointerMarshaler{FirstName: "Sally", LastName: "Fields"},
		Fiend:    &textPointerMarshaler{FirstName: "Casper", LastName: "Snider"},
	}

	result, err := Marshal(parent)
	if err != nil {
		t.Fatal(err)
	}
	expected := `fiend = "hidden"
stranger = "hidden"

[friend]
  FirstName = "Sally"
  LastName = "Fields"

[me]
  FirstName = "Maiku"
  LastName = "Suteda"
`
	if !bytes.Equal(result, []byte(expected)) {
		t.Errorf("Bad nested text marshaler: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestPointerCustomMarshalerSequence(t *testing.T) {
	var customPointerMarshalerSlice *[]*customPointerMarshaler
	var customPointerMarshalerArray *[2]*customPointerMarshaler

	if !isCustomMarshalerSequence(reflect.TypeOf(customPointerMarshalerSlice)) {
		t.Errorf("error: should be a sequence of custom marshaler interfaces")
	}
	if !isCustomMarshalerSequence(reflect.TypeOf(customPointerMarshalerArray)) {
		t.Errorf("error: should be a sequence of custom marshaler interfaces")
	}
}

func TestPointerTextMarshalerSequence(t *testing.T) {
	var textPointerMarshalerSlice *[]*textPointerMarshaler
	var textPointerMarshalerArray *[2]*textPointerMarshaler

	if !isTextMarshalerSequence(reflect.TypeOf(textPointerMarshalerSlice)) {
		t.Errorf("error: should be a sequence of text marshaler interfaces")
	}
	if !isTextMarshalerSequence(reflect.TypeOf(textPointerMarshalerArray)) {
		t.Errorf("error: should be a sequence of text marshaler interfaces")
	}
}

var commentTestToml = []byte(`
# it's a comment on type
[postgres]
  # isCommented = "dvalue"
  noComment = "cvalue"

  # A comment on AttrB with a
  # break line
  password = "bvalue"

  # A comment on AttrA
  user = "avalue"

  [[postgres.My]]

    # a comment on my on typeC
    My = "Foo"

  [[postgres.My]]

    # a comment on my on typeC
    My = "Baar"
`)

func TestMarshalComment(t *testing.T) {
	type TypeC struct {
		My string `comment:"a comment on my on typeC"`
	}
	type TypeB struct {
		AttrA string `toml:"user" comment:"A comment on AttrA"`
		AttrB string `toml:"password" comment:"A comment on AttrB with a\n break line"`
		AttrC string `toml:"noComment"`
		AttrD string `toml:"isCommented" commented:"true"`
		My    []TypeC
	}
	type TypeA struct {
		TypeB TypeB `toml:"postgres" comment:"it's a comment on type"`
	}

	ta := []TypeC{{My: "Foo"}, {My: "Baar"}}
	config := TypeA{TypeB{AttrA: "avalue", AttrB: "bvalue", AttrC: "cvalue", AttrD: "dvalue", My: ta}}
	result, err := Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	expected := commentTestToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestMarshalMultilineCommented(t *testing.T) {
	expectedToml := []byte(`# MultilineArray = [
  # 100,
  # 200,
  # 300,
# ]
# MultilineNestedArray = [
  # [
  # "a",
  # "b",
  # "c",
# ],
  # [
  # "d",
  # "e",
  # "f",
# ],
# ]
# MultilineString = """
# I
# am
# Allen"""
NonCommented = "Not commented line"
`)
	type StructWithMultiline struct {
		NonCommented         string
		MultilineString      string     `commented:"true" multiline:"true"`
		MultilineArray       []int      `commented:"true"`
		MultilineNestedArray [][]string `commented:"true"`
	}

	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	if err := enc.ArraysWithOneElementPerLine(true).Encode(StructWithMultiline{
		NonCommented:    "Not commented line",
		MultilineString: "I\nam\nAllen",
		MultilineArray:  []int{100, 200, 300},
		MultilineNestedArray: [][]string{
			{"a", "b", "c"},
			{"d", "e", "f"},
		},
	}); err == nil {
		result := buf.Bytes()
		if !bytes.Equal(result, expectedToml) {
			t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expectedToml, result)
		}
	} else {
		t.Fatal(err)
	}
}

func TestMarshalNonPrimitiveTypeCommented(t *testing.T) {
	expectedToml := []byte(`
# [CommentedMapField]

  # [CommentedMapField.CommentedMapField1]
    # SingleLineString = "This line should be commented out"

  # [CommentedMapField.CommentedMapField2]
    # SingleLineString = "This line should be commented out"

# [CommentedStructField]

  # [CommentedStructField.CommentedStructField]
    # MultilineArray = [
      # 1,
      # 2,
    # ]
    # MultilineNestedArray = [
      # [
      # 10,
      # 20,
    # ],
      # [
      # 100,
      # 200,
    # ],
    # ]
    # MultilineString = """
# This line
# should be
# commented out"""

  # [CommentedStructField.NotCommentedStructField]
    # MultilineArray = [
      # 1,
      # 2,
    # ]
    # MultilineNestedArray = [
      # [
      # 10,
      # 20,
    # ],
      # [
      # 100,
      # 200,
    # ],
    # ]
    # MultilineString = """
# This line
# should be
# commented out"""

[NotCommentedStructField]

  # [NotCommentedStructField.CommentedStructField]
    # MultilineArray = [
      # 1,
      # 2,
    # ]
    # MultilineNestedArray = [
      # [
      # 10,
      # 20,
    # ],
      # [
      # 100,
      # 200,
    # ],
    # ]
    # MultilineString = """
# This line
# should be
# commented out"""

  [NotCommentedStructField.NotCommentedStructField]
    MultilineArray = [
      3,
      4,
    ]
    MultilineNestedArray = [
      [
      30,
      40,
    ],
      [
      300,
      400,
    ],
    ]
    MultilineString = """
This line
should NOT be
commented out"""
`)
	type InnerStruct struct {
		MultilineString      string `multiline:"true"`
		MultilineArray       []int
		MultilineNestedArray [][]int
	}
	type MiddleStruct struct {
		NotCommentedStructField InnerStruct
		CommentedStructField    InnerStruct `commented:"true"`
	}
	type OuterStruct struct {
		CommentedStructField    MiddleStruct `commented:"true"`
		NotCommentedStructField MiddleStruct
		CommentedMapField       map[string]struct{ SingleLineString string } `commented:"true"`
	}

	commentedTestStruct := OuterStruct{
		CommentedStructField: MiddleStruct{
			NotCommentedStructField: InnerStruct{
				MultilineString:      "This line\nshould be\ncommented out",
				MultilineArray:       []int{1, 2},
				MultilineNestedArray: [][]int{{10, 20}, {100, 200}},
			},
			CommentedStructField: InnerStruct{
				MultilineString:      "This line\nshould be\ncommented out",
				MultilineArray:       []int{1, 2},
				MultilineNestedArray: [][]int{{10, 20}, {100, 200}},
			},
		},
		NotCommentedStructField: MiddleStruct{
			NotCommentedStructField: InnerStruct{
				MultilineString:      "This line\nshould NOT be\ncommented out",
				MultilineArray:       []int{3, 4},
				MultilineNestedArray: [][]int{{30, 40}, {300, 400}},
			},
			CommentedStructField: InnerStruct{
				MultilineString:      "This line\nshould be\ncommented out",
				MultilineArray:       []int{1, 2},
				MultilineNestedArray: [][]int{{10, 20}, {100, 200}},
			},
		},
		CommentedMapField: map[string]struct{ SingleLineString string }{
			"CommentedMapField1": {
				SingleLineString: "This line should be commented out",
			},
			"CommentedMapField2": {
				SingleLineString: "This line should be commented out",
			},
		},
	}

	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	if err := enc.ArraysWithOneElementPerLine(true).Encode(commentedTestStruct); err == nil {
		result := buf.Bytes()
		if !bytes.Equal(result, expectedToml) {
			t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expectedToml, result)
		}
	} else {
		t.Fatal(err)
	}
}

type mapsTestStruct struct {
	Simple map[string]string
	Paths  map[string]string
	Other  map[string]float64
	X      struct {
		Y struct {
			Z map[string]bool
		}
	}
}

var mapsTestData = mapsTestStruct{
	Simple: map[string]string{
		"one plus one": "two",
		"next":         "three",
	},
	Paths: map[string]string{
		"/this/is/a/path": "/this/is/also/a/path",
		"/heloo.txt":      "/tmp/lololo.txt",
	},
	Other: map[string]float64{
		"testing": 3.9999,
	},
	X: struct{ Y struct{ Z map[string]bool } }{
		Y: struct{ Z map[string]bool }{
			Z: map[string]bool{
				"is.Nested": true,
			},
		},
	},
}
var mapsTestToml = []byte(`
[Other]
  "testing" = 3.9999

[Paths]
  "/heloo.txt" = "/tmp/lololo.txt"
  "/this/is/a/path" = "/this/is/also/a/path"

[Simple]
  "next" = "three"
  "one plus one" = "two"

[X]

  [X.Y]

    [X.Y.Z]
      "is.Nested" = true
`)

func TestEncodeQuotedMapKeys(t *testing.T) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).QuoteMapKeys(true).Encode(mapsTestData); err != nil {
		t.Fatal(err)
	}
	result := buf.Bytes()
	expected := mapsTestToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad maps marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestDecodeQuotedMapKeys(t *testing.T) {
	result := mapsTestStruct{}
	err := NewDecoder(bytes.NewBuffer(mapsTestToml)).Decode(&result)
	expected := mapsTestData
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad maps unmarshal: expected %v, got %v", expected, result)
	}
}

type structArrayNoTag struct {
	A struct {
		B []int64
		C []int64
	}
}

func TestMarshalArray(t *testing.T) {
	expected := []byte(`
[A]
  B = [1, 2, 3]
  C = [1]
`)

	m := structArrayNoTag{
		A: struct {
			B []int64
			C []int64
		}{
			B: []int64{1, 2, 3},
			C: []int64{1},
		},
	}

	b, err := Marshal(m)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(b, expected) {
		t.Errorf("Bad arrays marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, b)
	}
}

func TestMarshalArrayOnePerLine(t *testing.T) {
	expected := []byte(`
[A]
  B = [
    1,
    2,
    3,
  ]
  C = [1]
`)

	m := structArrayNoTag{
		A: struct {
			B []int64
			C []int64
		}{
			B: []int64{1, 2, 3},
			C: []int64{1},
		},
	}

	var buf bytes.Buffer
	encoder := NewEncoder(&buf).ArraysWithOneElementPerLine(true)
	err := encoder.Encode(m)

	if err != nil {
		t.Fatal(err)
	}

	b := buf.Bytes()

	if !bytes.Equal(b, expected) {
		t.Errorf("Bad arrays marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, b)
	}
}

var customTagTestToml = []byte(`
[postgres]
  password = "bvalue"
  user = "avalue"

  [[postgres.My]]
    My = "Foo"

  [[postgres.My]]
    My = "Baar"
`)

func TestMarshalCustomTag(t *testing.T) {
	type TypeC struct {
		My string
	}
	type TypeB struct {
		AttrA string `file:"user"`
		AttrB string `file:"password"`
		My    []TypeC
	}
	type TypeA struct {
		TypeB TypeB `file:"postgres"`
	}

	ta := []TypeC{{My: "Foo"}, {My: "Baar"}}
	config := TypeA{TypeB{AttrA: "avalue", AttrB: "bvalue", My: ta}}
	var buf bytes.Buffer
	err := NewEncoder(&buf).SetTagName("file").Encode(config)
	if err != nil {
		t.Fatal(err)
	}
	expected := customTagTestToml
	result := buf.Bytes()
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

var customCommentTagTestToml = []byte(`
# db connection
[postgres]

  # db pass
  password = "bvalue"

  # db user
  user = "avalue"
`)

func TestMarshalCustomComment(t *testing.T) {
	type TypeB struct {
		AttrA string `toml:"user" descr:"db user"`
		AttrB string `toml:"password" descr:"db pass"`
	}
	type TypeA struct {
		TypeB TypeB `toml:"postgres" descr:"db connection"`
	}

	config := TypeA{TypeB{AttrA: "avalue", AttrB: "bvalue"}}
	var buf bytes.Buffer
	err := NewEncoder(&buf).SetTagComment("descr").Encode(config)
	if err != nil {
		t.Fatal(err)
	}
	expected := customCommentTagTestToml
	result := buf.Bytes()
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

var customCommentedTagTestToml = []byte(`
[postgres]
  # password = "bvalue"
  # user = "avalue"
`)

func TestMarshalCustomCommented(t *testing.T) {
	type TypeB struct {
		AttrA string `toml:"user" disable:"true"`
		AttrB string `toml:"password" disable:"true"`
	}
	type TypeA struct {
		TypeB TypeB `toml:"postgres"`
	}

	config := TypeA{TypeB{AttrA: "avalue", AttrB: "bvalue"}}
	var buf bytes.Buffer
	err := NewEncoder(&buf).SetTagCommented("disable").Encode(config)
	if err != nil {
		t.Fatal(err)
	}
	expected := customCommentedTagTestToml
	result := buf.Bytes()
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestMarshalDirectMultilineString(t *testing.T) {
	tree := newTree()
	tree.SetWithOptions("mykey", SetOptions{
		Multiline: true,
	}, "my\x11multiline\nstring\ba\tb\fc\rd\"e\\!")
	result, err := tree.Marshal()
	if err != nil {
		t.Fatal("marshal should not error:", err)
	}
	expected := []byte("mykey = \"\"\"\nmy\\u0011multiline\nstring\\ba\tb\\fc\rd\"e\\!\"\"\"\n")
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestUnmarshalTabInStringAndQuotedKey(t *testing.T) {
	type Test struct {
		Field1 string `toml:"Fie	ld1"`
		Field2 string
	}

	type TestCase struct {
		desc     string
		input    []byte
		expected Test
	}

	testCases := []TestCase{
		{
			desc:  "multiline string with tab",
			input: []byte("Field2 = \"\"\"\nhello\tworld\"\"\""),
			expected: Test{
				Field2: "hello\tworld",
			},
		},
		{
			desc:  "quoted key with tab",
			input: []byte("\"Fie\tld1\" = \"key with tab\""),
			expected: Test{
				Field1: "key with tab",
			},
		},
		{
			desc:  "basic string tab",
			input: []byte("Field2 = \"hello\tworld\""),
			expected: Test{
				Field2: "hello\tworld",
			},
		},
	}

	for i := range testCases {
		result := Test{}
		err := Unmarshal(testCases[i].input, &result)
		if err != nil {
			t.Errorf("%s test error:%v", testCases[i].desc, err)
			continue
		}

		if !reflect.DeepEqual(result, testCases[i].expected) {
			t.Errorf("%s test error: expected\n-----\n%+v\n-----\ngot\n-----\n%+v\n-----\n",
				testCases[i].desc, testCases[i].expected, result)
		}
	}
}

var customMultilineTagTestToml = []byte(`int_slice = [
  1,
  2,
  3,
]
`)

func TestMarshalCustomMultiline(t *testing.T) {
	type TypeA struct {
		AttrA []int `toml:"int_slice" mltln:"true"`
	}

	config := TypeA{AttrA: []int{1, 2, 3}}
	var buf bytes.Buffer
	err := NewEncoder(&buf).ArraysWithOneElementPerLine(true).SetTagMultiline("mltln").Encode(config)
	if err != nil {
		t.Fatal(err)
	}
	expected := customMultilineTagTestToml
	result := buf.Bytes()
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestMultilineWithAdjacentQuotationMarks(t *testing.T) {
	type testStruct struct {
		Str string `multiline:"true"`
	}
	type testCase struct {
		expected []byte
		data     testStruct
	}

	testCases := []testCase{
		{
			expected: []byte(`Str = """
hello\""""
`),
			data: testStruct{
				Str: "hello\"",
			},
		},
		{
			expected: []byte(`Str = """
""\"""\"""\""""
`),
			data: testStruct{
				Str: "\"\"\"\"\"\"\"\"\"",
			},
		},
	}
	for i := range testCases {
		result, err := Marshal(testCases[i].data)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(result, testCases[i].expected) {
			t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n",
				testCases[i].expected, result)
		} else {
			var data testStruct
			if err = Unmarshal(result, &data); err != nil {
				t.Fatal(err)
			}
			if data.Str != testCases[i].data.Str {
				t.Errorf("Round trip test fail: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n",
					testCases[i].data.Str, data.Str)
			}
		}
	}
}

func TestMarshalEmbedTree(t *testing.T) {
	expected := []byte(`OuterField1 = "Out"
OuterField2 = 1024

[TreeField]
  InnerField1 = "In"
  InnerField2 = 2048

  [TreeField.EmbedStruct]
    EmbedField = "Embed"
`)
	type InnerStruct struct {
		InnerField1 string
		InnerField2 int
		EmbedStruct struct {
			EmbedField string
		}
	}

	type OuterStruct struct {
		OuterField1 string
		OuterField2 int
		TreeField   *Tree
	}

	tree, err := Load(`
InnerField1 = "In"
InnerField2 = 2048

[EmbedStruct]
	EmbedField = "Embed"
`)
	if err != nil {
		t.Fatal(err)
	}

	out := OuterStruct{
		"Out",
		1024,
		tree,
	}
	actual, _ := Marshal(out)

	if !bytes.Equal(actual, expected) {
		t.Errorf("Bad marshal: expected %s, got %s", expected, actual)
	}
}

var testDocBasicToml = []byte(`
[document]
  bool_val = true
  date_val = 1979-05-27T07:32:00Z
  float_val = 123.4
  int_val = 5000
  string_val = "Bite me"
  uint_val = 5001
`)

type testDocCustomTag struct {
	Doc testDocBasicsCustomTag `file:"document"`
}
type testDocBasicsCustomTag struct {
	Bool       bool      `file:"bool_val"`
	Date       time.Time `file:"date_val"`
	Float      float32   `file:"float_val"`
	Int        int       `file:"int_val"`
	Uint       uint      `file:"uint_val"`
	String     *string   `file:"string_val"`
	unexported int       `file:"shouldntBeHere"`
}

var testDocCustomTagData = testDocCustomTag{
	Doc: testDocBasicsCustomTag{
		Bool:       true,
		Date:       time.Date(1979, 5, 27, 7, 32, 0, 0, time.UTC),
		Float:      123.4,
		Int:        5000,
		Uint:       5001,
		String:     &biteMe,
		unexported: 0,
	},
}

func TestUnmarshalCustomTag(t *testing.T) {
	buf := bytes.NewBuffer(testDocBasicToml)

	result := testDocCustomTag{}
	err := NewDecoder(buf).SetTagName("file").Decode(&result)
	if err != nil {
		t.Fatal(err)
	}
	expected := testDocCustomTagData
	if !reflect.DeepEqual(result, expected) {
		resStr, _ := json.MarshalIndent(result, "", "  ")
		expStr, _ := json.MarshalIndent(expected, "", "  ")
		t.Errorf("Bad unmarshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expStr, resStr)

	}
}

func TestUnmarshalMap(t *testing.T) {
	testToml := []byte(`
		a = 1
		b = 2
		c = 3
		`)
	var result map[string]int
	err := Unmarshal(testToml, &result)
	if err != nil {
		t.Errorf("Received unexpected error: %s", err)
		return
	}

	expected := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad unmarshal: expected %v, got %v", expected, result)
	}
}

func TestUnmarshalMapWithTypedKey(t *testing.T) {
	testToml := []byte(`
		a = 1
		b = 2
		c = 3
		`)

	type letter string
	var result map[letter]int
	err := Unmarshal(testToml, &result)
	if err != nil {
		t.Errorf("Received unexpected error: %s", err)
		return
	}

	expected := map[letter]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad unmarshal: expected %v, got %v", expected, result)
	}
}

func TestUnmarshalNonPointer(t *testing.T) {
	a := 1
	err := Unmarshal([]byte{}, a)
	if err == nil {
		t.Fatal("unmarshal should err when given a non pointer")
	}
}

func TestUnmarshalInvalidPointerKind(t *testing.T) {
	a := 1
	err := Unmarshal([]byte{}, &a)
	if err == nil {
		t.Fatal("unmarshal should err when given an invalid pointer type")
	}
}

func TestMarshalSlice(t *testing.T) {
	m := make([]int, 1)
	m[0] = 1

	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(&m)
	if err == nil {
		t.Error("expected error, got nil")
		return
	}
	if err.Error() != "Only pointer to struct can be marshaled to TOML" {
		t.Fail()
	}
}

func TestMarshalSlicePointer(t *testing.T) {
	m := make([]int, 1)
	m[0] = 1

	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(m)
	if err == nil {
		t.Error("expected error, got nil")
		return
	}
	if err.Error() != "Only a struct or map can be marshaled to TOML" {
		t.Fail()
	}
}

func TestMarshalNestedArrayInlineTables(t *testing.T) {
	type table struct {
		Value1 int `toml:"ZValue1"`
		Value2 int `toml:"YValue2"`
		Value3 int `toml:"XValue3"`
	}

	type nestedTable struct {
		Table table
	}

	nestedArray := struct {
		Simple        [][]table
		SimplePointer *[]*[]table
		Nested        [][]nestedTable
		NestedPointer *[]*[]nestedTable
	}{
		Simple:        [][]table{{{Value1: 1}, {Value1: 10}}},
		SimplePointer: &[]*[]table{{{Value2: 2}}},
		Nested:        [][]nestedTable{{{Table: table{Value3: 3}}}},
		NestedPointer: &[]*[]nestedTable{{{Table: table{Value3: -3}}}},
	}

	expectedPreserve := `Simple = [[{ ZValue1 = 1, YValue2 = 0, XValue3 = 0 }, { ZValue1 = 10, YValue2 = 0, XValue3 = 0 }]]
SimplePointer = [[{ ZValue1 = 0, YValue2 = 2, XValue3 = 0 }]]
Nested = [[{ Table = { ZValue1 = 0, YValue2 = 0, XValue3 = 3 } }]]
NestedPointer = [[{ Table = { ZValue1 = 0, YValue2 = 0, XValue3 = -3 } }]]
`

	expectedAlphabetical := `Nested = [[{ Table = { XValue3 = 3, YValue2 = 0, ZValue1 = 0 } }]]
NestedPointer = [[{ Table = { XValue3 = -3, YValue2 = 0, ZValue1 = 0 } }]]
Simple = [[{ XValue3 = 0, YValue2 = 0, ZValue1 = 1 }, { XValue3 = 0, YValue2 = 0, ZValue1 = 10 }]]
SimplePointer = [[{ XValue3 = 0, YValue2 = 2, ZValue1 = 0 }]]
`

	var bufPreserve bytes.Buffer
	if err := NewEncoder(&bufPreserve).Order(OrderPreserve).Encode(nestedArray); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if !bytes.Equal(bufPreserve.Bytes(), []byte(expectedPreserve)) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expectedPreserve, bufPreserve.String())
	}

	var bufAlphabetical bytes.Buffer
	if err := NewEncoder(&bufAlphabetical).Order(OrderAlphabetical).Encode(nestedArray); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if !bytes.Equal(bufAlphabetical.Bytes(), []byte(expectedAlphabetical)) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expectedAlphabetical, bufAlphabetical.String())
	}
}

type testDuration struct {
	Nanosec   time.Duration  `toml:"nanosec"`
	Microsec1 time.Duration  `toml:"microsec1"`
	Microsec2 *time.Duration `toml:"microsec2"`
	Millisec  time.Duration  `toml:"millisec"`
	Sec       time.Duration  `toml:"sec"`
	Min       time.Duration  `toml:"min"`
	Hour      time.Duration  `toml:"hour"`
	Mixed     time.Duration  `toml:"mixed"`
	AString   string         `toml:"a_string"`
}

var testDurationToml = []byte(`
nanosec = "1ns"
microsec1 = "1us"
microsec2 = "1µs"
millisec = "1ms"
sec = "1s"
min = "1m"
hour = "1h"
mixed = "1h1m1s1ms1µs1ns"
a_string = "15s"
`)

func TestUnmarshalDuration(t *testing.T) {
	buf := bytes.NewBuffer(testDurationToml)

	result := testDuration{}
	err := NewDecoder(buf).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}
	ms := time.Duration(1) * time.Microsecond
	expected := testDuration{
		Nanosec:   1,
		Microsec1: time.Microsecond,
		Microsec2: &ms,
		Millisec:  time.Millisecond,
		Sec:       time.Second,
		Min:       time.Minute,
		Hour:      time.Hour,
		Mixed: time.Hour +
			time.Minute +
			time.Second +
			time.Millisecond +
			time.Microsecond +
			time.Nanosecond,
		AString: "15s",
	}
	if !reflect.DeepEqual(result, expected) {
		resStr, _ := json.MarshalIndent(result, "", "  ")
		expStr, _ := json.MarshalIndent(expected, "", "  ")
		t.Errorf("Bad unmarshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expStr, resStr)

	}
}

var testDurationToml2 = []byte(`a_string = "15s"
hour = "1h0m0s"
microsec1 = "1µs"
microsec2 = "1µs"
millisec = "1ms"
min = "1m0s"
mixed = "1h1m1.001001001s"
nanosec = "1ns"
sec = "1s"
`)

func TestMarshalDuration(t *testing.T) {
	ms := time.Duration(1) * time.Microsecond
	data := testDuration{
		Nanosec:   1,
		Microsec1: time.Microsecond,
		Microsec2: &ms,
		Millisec:  time.Millisecond,
		Sec:       time.Second,
		Min:       time.Minute,
		Hour:      time.Hour,
		Mixed: time.Hour +
			time.Minute +
			time.Second +
			time.Millisecond +
			time.Microsecond +
			time.Nanosecond,
		AString: "15s",
	}

	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(data)
	if err != nil {
		t.Fatal(err)
	}
	expected := testDurationToml2
	result := buf.Bytes()
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

type testBadDuration struct {
	Val time.Duration `toml:"val"`
}

var testBadDurationToml = []byte(`val = "1z"`)

func TestUnmarshalBadDuration(t *testing.T) {
	buf := bytes.NewBuffer(testBadDurationToml)

	result := testBadDuration{}
	err := NewDecoder(buf).Decode(&result)
	if err == nil {
		t.Fatal()
	}
	if err.Error() != "(1, 1): Can't convert 1z(string) to time.Duration. time: unknown unit z in duration 1z" {
		t.Fatalf("unexpected error: %s", err)
	}
}

var testCamelCaseKeyToml = []byte(`fooBar = 10`)

func TestUnmarshalCamelCaseKey(t *testing.T) {
	var x struct {
		FooBar int
		B      int
	}

	if err := Unmarshal(testCamelCaseKeyToml, &x); err != nil {
		t.Fatal(err)
	}

	if x.FooBar != 10 {
		t.Fatal("Did not set camelCase'd key")
	}
}

func TestUnmarshalNegativeUint(t *testing.T) {
	type check struct{ U uint }

	tree, _ := Load("u = -1")
	err := tree.Unmarshal(&check{})
	if err.Error() != "(1, 1): -1(int64) is negative so does not fit in uint" {
		t.Error("expect err:(1, 1): -1(int64) is negative so does not fit in uint but got:", err)
	}
}

func TestUnmarshalCheckConversionFloatInt(t *testing.T) {
	type conversionCheck struct {
		U uint
		I int
		F float64
	}

	treeU, _ := Load("u = 1e300")
	treeI, _ := Load("i = 1e300")
	treeF, _ := Load("f = 9223372036854775806")

	errU := treeU.Unmarshal(&conversionCheck{})
	errI := treeI.Unmarshal(&conversionCheck{})
	errF := treeF.Unmarshal(&conversionCheck{})

	if errU.Error() != "(1, 1): Can't convert 1e+300(float64) to uint" {
		t.Error("expect err:(1, 1): Can't convert 1e+300(float64) to uint but got:", errU)
	}
	if errI.Error() != "(1, 1): Can't convert 1e+300(float64) to int" {
		t.Error("expect err:(1, 1): Can't convert 1e+300(float64) to int but got:", errI)
	}
	if errF.Error() != "(1, 1): Can't convert 9223372036854775806(int64) to float64" {
		t.Error("expect err:(1, 1): Can't convert 9223372036854775806(int64) to float64 but got:", errF)
	}
}

func TestUnmarshalOverflow(t *testing.T) {
	type overflow struct {
		U8  uint8
		I8  int8
		F32 float32
	}

	treeU8, _ := Load("u8 = 300")
	treeI8, _ := Load("i8 = 300")
	treeF32, _ := Load("f32 = 1e300")

	errU8 := treeU8.Unmarshal(&overflow{})
	errI8 := treeI8.Unmarshal(&overflow{})
	errF32 := treeF32.Unmarshal(&overflow{})

	if errU8.Error() != "(1, 1): 300(int64) would overflow uint8" {
		t.Error("expect err:(1, 1): 300(int64) would overflow uint8 but got:", errU8)
	}
	if errI8.Error() != "(1, 1): 300(int64) would overflow int8" {
		t.Error("expect err:(1, 1): 300(int64) would overflow int8 but got:", errI8)
	}
	if errF32.Error() != "(1, 1): 1e+300(float64) would overflow float32" {
		t.Error("expect err:(1, 1): 1e+300(float64) would overflow float32 but got:", errF32)
	}
}

func TestUnmarshalDefault(t *testing.T) {
	type EmbeddedStruct struct {
		StringField string `default:"c"`
	}

	type aliasUint uint

	var doc struct {
		StringField       string  `default:"a"`
		BoolField         bool    `default:"true"`
		UintField         uint    `default:"1"`
		Uint8Field        uint8   `default:"8"`
		Uint16Field       uint16  `default:"16"`
		Uint32Field       uint32  `default:"32"`
		Uint64Field       uint64  `default:"64"`
		IntField          int     `default:"-1"`
		Int8Field         int8    `default:"-8"`
		Int16Field        int16   `default:"-16"`
		Int32Field        int32   `default:"-32"`
		Int64Field        int64   `default:"-64"`
		Float32Field      float32 `default:"32.1"`
		Float64Field      float64 `default:"64.1"`
		NonEmbeddedStruct struct {
			StringField string `default:"b"`
		}
		EmbeddedStruct
		AliasUintField aliasUint `default:"1000"`
	}

	err := Unmarshal([]byte(``), &doc)
	if err != nil {
		t.Fatal(err)
	}
	if doc.BoolField != true {
		t.Errorf("BoolField should be true, not %t", doc.BoolField)
	}
	if doc.StringField != "a" {
		t.Errorf("StringField should be \"a\", not %s", doc.StringField)
	}
	if doc.UintField != 1 {
		t.Errorf("UintField should be 1, not %d", doc.UintField)
	}
	if doc.Uint8Field != 8 {
		t.Errorf("Uint8Field should be 8, not %d", doc.Uint8Field)
	}
	if doc.Uint16Field != 16 {
		t.Errorf("Uint16Field should be 16, not %d", doc.Uint16Field)
	}
	if doc.Uint32Field != 32 {
		t.Errorf("Uint32Field should be 32, not %d", doc.Uint32Field)
	}
	if doc.Uint64Field != 64 {
		t.Errorf("Uint64Field should be 64, not %d", doc.Uint64Field)
	}
	if doc.IntField != -1 {
		t.Errorf("IntField should be -1, not %d", doc.IntField)
	}
	if doc.Int8Field != -8 {
		t.Errorf("Int8Field should be -8, not %d", doc.Int8Field)
	}
	if doc.Int16Field != -16 {
		t.Errorf("Int16Field should be -16, not %d", doc.Int16Field)
	}
	if doc.Int32Field != -32 {
		t.Errorf("Int32Field should be -32, not %d", doc.Int32Field)
	}
	if doc.Int64Field != -64 {
		t.Errorf("Int64Field should be -64, not %d", doc.Int64Field)
	}
	if doc.Float32Field != 32.1 {
		t.Errorf("Float32Field should be 32.1, not %f", doc.Float32Field)
	}
	if doc.Float64Field != 64.1 {
		t.Errorf("Float64Field should be 64.1, not %f", doc.Float64Field)
	}
	if doc.NonEmbeddedStruct.StringField != "b" {
		t.Errorf("StringField should be \"b\", not %s", doc.NonEmbeddedStruct.StringField)
	}
	if doc.EmbeddedStruct.StringField != "c" {
		t.Errorf("StringField should be \"c\", not %s", doc.EmbeddedStruct.StringField)
	}
	if doc.AliasUintField != 1000 {
		t.Errorf("AliasUintField should be 1000, not %d", doc.AliasUintField)
	}
}

func TestUnmarshalDefaultFailureBool(t *testing.T) {
	var doc struct {
		Field bool `default:"blah"`
	}

	err := Unmarshal([]byte(``), &doc)
	if err == nil {
		t.Fatal("should error")
	}
}

func TestUnmarshalDefaultFailureInt(t *testing.T) {
	var doc struct {
		Field int `default:"blah"`
	}

	err := Unmarshal([]byte(``), &doc)
	if err == nil {
		t.Fatal("should error")
	}
}

func TestUnmarshalDefaultFailureInt64(t *testing.T) {
	var doc struct {
		Field int64 `default:"blah"`
	}

	err := Unmarshal([]byte(``), &doc)
	if err == nil {
		t.Fatal("should error")
	}
}

func TestUnmarshalDefaultFailureFloat64(t *testing.T) {
	var doc struct {
		Field float64 `default:"blah"`
	}

	err := Unmarshal([]byte(``), &doc)
	if err == nil {
		t.Fatal("should error")
	}
}

func TestUnmarshalDefaultFailureUnsupported(t *testing.T) {
	var doc struct {
		Field struct{} `default:"blah"`
	}

	err := Unmarshal([]byte(``), &doc)
	if err == nil {
		t.Fatal("should error")
	}
}

func TestMarshalNestedAnonymousStructs(t *testing.T) {
	type Embedded struct {
		Value string `toml:"value"`
		Top   struct {
			Value string `toml:"value"`
		} `toml:"top"`
	}

	type Named struct {
		Value string `toml:"value"`
	}

	var doc struct {
		Embedded
		Named     `toml:"named"`
		Anonymous struct {
			Value string `toml:"value"`
		} `toml:"anonymous"`
	}

	expected := `value = ""

[anonymous]
  value = ""

[named]
  value = ""

[top]
  value = ""
`

	result, err := Marshal(doc)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if !bytes.Equal(result, []byte(expected)) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, string(result))
	}
}

func TestEncoderPromoteNestedAnonymousStructs(t *testing.T) {
	type Embedded struct {
		Value string `toml:"value"`
	}

	var doc struct {
		Embedded
	}

	expected := `
[Embedded]
  value = ""
`
	var buf bytes.Buffer
	if err := NewEncoder(&buf).PromoteAnonymous(true).Encode(doc); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if !bytes.Equal(buf.Bytes(), []byte(expected)) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, buf.String())
	}
}

func TestMarshalNestedAnonymousStructs_DuplicateField(t *testing.T) {
	type Embedded struct {
		Value string `toml:"value"`
		Top   struct {
			Value string `toml:"value"`
		} `toml:"top"`
	}

	var doc struct {
		Value string `toml:"value"`
		Embedded
	}
	doc.Embedded.Value = "shadowed"
	doc.Value = "shadows"

	expected := `value = "shadows"

[top]
  value = ""
`

	result, err := Marshal(doc)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if !bytes.Equal(result, []byte(expected)) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, string(result))
	}
}

func TestUnmarshalNestedAnonymousStructs(t *testing.T) {
	type Nested struct {
		Value string `toml:"nested_field"`
	}
	type Deep struct {
		Nested
	}
	type Document struct {
		Deep
		Value string `toml:"own_field"`
	}

	var doc Document

	err := Unmarshal([]byte(`nested_field = "nested value"`+"\n"+`own_field = "own value"`), &doc)
	if err != nil {
		t.Fatal("should not error")
	}
	if doc.Value != "own value" || doc.Nested.Value != "nested value" {
		t.Fatal("unexpected values")
	}
}

func TestUnmarshalNestedAnonymousStructs_Controversial(t *testing.T) {
	type Nested struct {
		Value string `toml:"nested"`
	}
	type Deep struct {
		Nested
	}
	type Document struct {
		Deep
		Value string `toml:"own"`
	}

	var doc Document

	err := Unmarshal([]byte(`nested = "nested value"`+"\n"+`own = "own value"`), &doc)
	if err == nil {
		t.Fatal("should error")
	}
}

type unexportedFieldPreservationTest struct {
	Exported   string `toml:"exported"`
	unexported string
	Nested1    unexportedFieldPreservationTestNested    `toml:"nested1"`
	Nested2    *unexportedFieldPreservationTestNested   `toml:"nested2"`
	Nested3    *unexportedFieldPreservationTestNested   `toml:"nested3"`
	Slice1     []unexportedFieldPreservationTestNested  `toml:"slice1"`
	Slice2     []*unexportedFieldPreservationTestNested `toml:"slice2"`
}

type unexportedFieldPreservationTestNested struct {
	Exported1   string `toml:"exported1"`
	unexported1 string
}

func TestUnmarshalPreservesUnexportedFields(t *testing.T) {
	toml := `
	exported = "visible"
	unexported = "ignored"

	[nested1]
	exported1 = "visible1"
	unexported1 = "ignored1"

	[nested2]
	exported1 = "visible2"
	unexported1 = "ignored2"

	[nested3]
	exported1 = "visible3"
	unexported1 = "ignored3"

	[[slice1]]
	exported1 = "visible3"

	[[slice1]]
	exported1 = "visible4"

	[[slice2]]
	exported1 = "visible5"
	`

	t.Run("unexported field should not be set from toml", func(t *testing.T) {
		var actual unexportedFieldPreservationTest
		err := Unmarshal([]byte(toml), &actual)

		if err != nil {
			t.Fatal("did not expect an error")
		}

		expect := unexportedFieldPreservationTest{
			Exported:   "visible",
			unexported: "",
			Nested1:    unexportedFieldPreservationTestNested{"visible1", ""},
			Nested2:    &unexportedFieldPreservationTestNested{"visible2", ""},
			Nested3:    &unexportedFieldPreservationTestNested{"visible3", ""},
			Slice1: []unexportedFieldPreservationTestNested{
				{Exported1: "visible3"},
				{Exported1: "visible4"},
			},
			Slice2: []*unexportedFieldPreservationTestNested{
				{Exported1: "visible5"},
			},
		}

		if !reflect.DeepEqual(actual, expect) {
			t.Fatalf("%+v did not equal %+v", actual, expect)
		}
	})

	t.Run("unexported field should be preserved", func(t *testing.T) {
		actual := unexportedFieldPreservationTest{
			Exported:   "foo",
			unexported: "bar",
			Nested1:    unexportedFieldPreservationTestNested{"baz", "bax"},
			Nested2:    nil,
			Nested3:    &unexportedFieldPreservationTestNested{"baz", "bax"},
		}
		err := Unmarshal([]byte(toml), &actual)

		if err != nil {
			t.Fatal("did not expect an error")
		}

		expect := unexportedFieldPreservationTest{
			Exported:   "visible",
			unexported: "bar",
			Nested1:    unexportedFieldPreservationTestNested{"visible1", "bax"},
			Nested2:    &unexportedFieldPreservationTestNested{"visible2", ""},
			Nested3:    &unexportedFieldPreservationTestNested{"visible3", "bax"},
			Slice1: []unexportedFieldPreservationTestNested{
				{Exported1: "visible3"},
				{Exported1: "visible4"},
			},
			Slice2: []*unexportedFieldPreservationTestNested{
				{Exported1: "visible5"},
			},
		}

		if !reflect.DeepEqual(actual, expect) {
			t.Fatalf("%+v did not equal %+v", actual, expect)
		}
	})
}

func TestTreeMarshal(t *testing.T) {
	cases := [][]byte{
		basicTestToml,
		marshalTestToml,
		emptyTestToml,
		pointerTestToml,
	}
	for _, expected := range cases {
		t.Run("", func(t *testing.T) {
			tree, err := LoadBytes(expected)
			if err != nil {
				t.Fatal(err)
			}
			result, err := tree.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(result, expected) {
				t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
			}
		})
	}
}

func TestMarshalArrays(t *testing.T) {
	cases := []struct {
		Data     interface{}
		Expected string
	}{
		{
			Data: struct {
				XY [2]int
			}{
				XY: [2]int{1, 2},
			},
			Expected: `XY = [1, 2]
`,
		},
		{
			Data: struct {
				XY [1][2]int
			}{
				XY: [1][2]int{{1, 2}},
			},
			Expected: `XY = [[1, 2]]
`,
		},
		{
			Data: struct {
				XY [1][]int
			}{
				XY: [1][]int{{1, 2}},
			},
			Expected: `XY = [[1, 2]]
`,
		},
		{
			Data: struct {
				XY [][2]int
			}{
				XY: [][2]int{{1, 2}},
			},
			Expected: `XY = [[1, 2]]
`,
		},
	}
	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			result, err := Marshal(tc.Data)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(result, []byte(tc.Expected)) {
				t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", []byte(tc.Expected), result)
			}
		})
	}
}

func TestUnmarshalLocalDate(t *testing.T) {
	t.Run("ToLocalDate", func(t *testing.T) {
		type dateStruct struct {
			Date LocalDate
		}

		toml := `date = 1979-05-27`

		var obj dateStruct

		err := Unmarshal([]byte(toml), &obj)

		if err != nil {
			t.Fatal(err)
		}

		if obj.Date.Year != 1979 {
			t.Errorf("expected year 1979, got %d", obj.Date.Year)
		}
		if obj.Date.Month != 5 {
			t.Errorf("expected month 5, got %d", obj.Date.Month)
		}
		if obj.Date.Day != 27 {
			t.Errorf("expected day 27, got %d", obj.Date.Day)
		}
	})

	t.Run("ToLocalDate", func(t *testing.T) {
		type dateStruct struct {
			Date time.Time
		}

		toml := `date = 1979-05-27`

		var obj dateStruct

		err := Unmarshal([]byte(toml), &obj)

		if err != nil {
			t.Fatal(err)
		}

		if obj.Date.Year() != 1979 {
			t.Errorf("expected year 1979, got %d", obj.Date.Year())
		}
		if obj.Date.Month() != 5 {
			t.Errorf("expected month 5, got %d", obj.Date.Month())
		}
		if obj.Date.Day() != 27 {
			t.Errorf("expected day 27, got %d", obj.Date.Day())
		}
	})
}

func TestMarshalLocalDate(t *testing.T) {
	type dateStruct struct {
		Date LocalDate
	}

	obj := dateStruct{Date: LocalDate{
		Year:  1979,
		Month: 5,
		Day:   27,
	}}

	b, err := Marshal(obj)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := string(b)
	expected := `Date = 1979-05-27
`

	if got != expected {
		t.Errorf("expected '%s', got '%s'", expected, got)
	}
}

func TestUnmarshalLocalDateTime(t *testing.T) {
	examples := []struct {
		name string
		in   string
		out  LocalDateTime
	}{
		{
			name: "normal",
			in:   "1979-05-27T07:32:00",
			out: LocalDateTime{
				Date: LocalDate{
					Year:  1979,
					Month: 5,
					Day:   27,
				},
				Time: LocalTime{
					Hour:       7,
					Minute:     32,
					Second:     0,
					Nanosecond: 0,
				},
			}},
		{
			name: "with nanoseconds",
			in:   "1979-05-27T00:32:00.999999",
			out: LocalDateTime{
				Date: LocalDate{
					Year:  1979,
					Month: 5,
					Day:   27,
				},
				Time: LocalTime{
					Hour:       0,
					Minute:     32,
					Second:     0,
					Nanosecond: 999999000,
				},
			},
		},
	}

	for i, example := range examples {
		toml := fmt.Sprintf(`date = %s`, example.in)

		t.Run(fmt.Sprintf("ToLocalDateTime_%d_%s", i, example.name), func(t *testing.T) {
			type dateStruct struct {
				Date LocalDateTime
			}

			var obj dateStruct

			err := Unmarshal([]byte(toml), &obj)

			if err != nil {
				t.Fatal(err)
			}

			if obj.Date != example.out {
				t.Errorf("expected '%s', got '%s'", example.out, obj.Date)
			}
		})

		t.Run(fmt.Sprintf("ToTime_%d_%s", i, example.name), func(t *testing.T) {
			type dateStruct struct {
				Date time.Time
			}

			var obj dateStruct

			err := Unmarshal([]byte(toml), &obj)

			if err != nil {
				t.Fatal(err)
			}

			if obj.Date.Year() != example.out.Date.Year {
				t.Errorf("expected year %d, got %d", example.out.Date.Year, obj.Date.Year())
			}
			if obj.Date.Month() != example.out.Date.Month {
				t.Errorf("expected month %d, got %d", example.out.Date.Month, obj.Date.Month())
			}
			if obj.Date.Day() != example.out.Date.Day {
				t.Errorf("expected day %d, got %d", example.out.Date.Day, obj.Date.Day())
			}
			if obj.Date.Hour() != example.out.Time.Hour {
				t.Errorf("expected hour %d, got %d", example.out.Time.Hour, obj.Date.Hour())
			}
			if obj.Date.Minute() != example.out.Time.Minute {
				t.Errorf("expected minute %d, got %d", example.out.Time.Minute, obj.Date.Minute())
			}
			if obj.Date.Second() != example.out.Time.Second {
				t.Errorf("expected second %d, got %d", example.out.Time.Second, obj.Date.Second())
			}
			if obj.Date.Nanosecond() != example.out.Time.Nanosecond {
				t.Errorf("expected nanoseconds %d, got %d", example.out.Time.Nanosecond, obj.Date.Nanosecond())
			}
		})
	}
}

func TestMarshalLocalDateTime(t *testing.T) {
	type dateStruct struct {
		DateTime LocalDateTime
	}

	examples := []struct {
		name string
		in   LocalDateTime
		out  string
	}{
		{
			name: "normal",
			out:  "DateTime = 1979-05-27T07:32:00\n",
			in: LocalDateTime{
				Date: LocalDate{
					Year:  1979,
					Month: 5,
					Day:   27,
				},
				Time: LocalTime{
					Hour:       7,
					Minute:     32,
					Second:     0,
					Nanosecond: 0,
				},
			}},
		{
			name: "with nanoseconds",
			out:  "DateTime = 1979-05-27T00:32:00.999999000\n",
			in: LocalDateTime{
				Date: LocalDate{
					Year:  1979,
					Month: 5,
					Day:   27,
				},
				Time: LocalTime{
					Hour:       0,
					Minute:     32,
					Second:     0,
					Nanosecond: 999999000,
				},
			},
		},
	}

	for i, example := range examples {
		t.Run(fmt.Sprintf("%d_%s", i, example.name), func(t *testing.T) {
			obj := dateStruct{
				DateTime: example.in,
			}
			b, err := Marshal(obj)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got := string(b)

			if got != example.out {
				t.Errorf("expected '%s', got '%s'", example.out, got)
			}
		})
	}
}

func TestUnmarshalLocalTime(t *testing.T) {
	examples := []struct {
		name string
		in   string
		out  LocalTime
	}{
		{
			name: "normal",
			in:   "07:32:00",
			out: LocalTime{
				Hour:       7,
				Minute:     32,
				Second:     0,
				Nanosecond: 0,
			},
		},
		{
			name: "with nanoseconds",
			in:   "00:32:00.999999",
			out: LocalTime{
				Hour:       0,
				Minute:     32,
				Second:     0,
				Nanosecond: 999999000,
			},
		},
	}

	for i, example := range examples {
		toml := fmt.Sprintf(`Time = %s`, example.in)

		t.Run(fmt.Sprintf("ToLocalTime_%d_%s", i, example.name), func(t *testing.T) {
			type dateStruct struct {
				Time LocalTime
			}

			var obj dateStruct

			err := Unmarshal([]byte(toml), &obj)

			if err != nil {
				t.Fatal(err)
			}

			if obj.Time != example.out {
				t.Errorf("expected '%s', got '%s'", example.out, obj.Time)
			}
		})
	}
}

func TestMarshalLocalTime(t *testing.T) {
	type timeStruct struct {
		Time LocalTime
	}

	examples := []struct {
		name string
		in   LocalTime
		out  string
	}{
		{
			name: "normal",
			out:  "Time = 07:32:00\n",
			in: LocalTime{
				Hour:       7,
				Minute:     32,
				Second:     0,
				Nanosecond: 0,
			}},
		{
			name: "with nanoseconds",
			out:  "Time = 00:32:00.999999000\n",
			in: LocalTime{
				Hour:       0,
				Minute:     32,
				Second:     0,
				Nanosecond: 999999000,
			},
		},
	}

	for i, example := range examples {
		t.Run(fmt.Sprintf("%d_%s", i, example.name), func(t *testing.T) {
			obj := timeStruct{
				Time: example.in,
			}
			b, err := Marshal(obj)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got := string(b)

			if got != example.out {
				t.Errorf("expected '%s', got '%s'", example.out, got)
			}
		})
	}
}

// test case for issue #339
func TestUnmarshalSameInnerField(t *testing.T) {
	type InterStruct2 struct {
		Test string
		Name string
		Age  int
	}
	type Inter2 struct {
		Name         string
		Age          int
		InterStruct2 InterStruct2
	}
	type Server struct {
		Name   string `toml:"name"`
		Inter2 Inter2 `toml:"inter2"`
	}

	var server Server

	if err := Unmarshal([]byte(`name = "123"
[inter2]
name = "inter2"
age = 222`), &server); err == nil {
		expected := Server{
			Name: "123",
			Inter2: Inter2{
				Name: "inter2",
				Age:  222,
			},
		}
		if !reflect.DeepEqual(server, expected) {
			t.Errorf("Bad unmarshal: expected %v, got %v", expected, server)
		}
	} else {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMarshalInterface(t *testing.T) {
	type InnerStruct struct {
		InnerField string
	}

	type OuterStruct struct {
		PrimitiveField        interface{}
		ArrayField            interface{}
		StructArrayField      interface{}
		MapField              map[string]interface{}
		StructField           interface{}
		PointerField          interface{}
		NilField              interface{}
		InterfacePointerField *interface{}
	}

	expected := []byte(`ArrayField = [1, 2, 3]
InterfacePointerField = "hello world"
PrimitiveField = "string"

[MapField]
  key1 = "value1"
  key2 = false

  [MapField.key3]
    InnerField = "value3"

[PointerField]
  InnerField = "yyy"

[[StructArrayField]]
  InnerField = "s1"

[[StructArrayField]]
  InnerField = "s2"

[StructField]
  InnerField = "xxx"
`)

	var h interface{} = "hello world"
	if result, err := Marshal(OuterStruct{
		"string",
		[]int{1, 2, 3},
		[]InnerStruct{{"s1"}, {"s2"}},
		map[string]interface{}{
			"key1":      "value1",
			"key2":      false,
			"key3":      InnerStruct{"value3"},
			"nil value": nil,
		},
		InnerStruct{
			"xxx",
		},
		&InnerStruct{
			"yyy",
		},
		nil,
		&h,
	}); err == nil {
		if !bytes.Equal(result, expected) {
			t.Errorf("Bad marshal: expected\n----\n%s\n----\ngot\n----\n%s\n----\n", expected, result)
		}
	} else {
		t.Fatal(err)
	}
}

func TestUnmarshalToNilInterface(t *testing.T) {
	toml := []byte(`
PrimitiveField = "Hello"
ArrayField = [1,2,3]
InterfacePointerField = "World"

[StructField]
Field1 = 123
Field2 = "Field2"

[MapField]
MapField1 = [4,5,6]
MapField2 = {A = "A"}
MapField3 = false

[[StructArrayField]]
Name = "Allen"
Age = 20

[[StructArrayField]]
Name = "Jack"
Age = 23
`)

	type OuterStruct struct {
		PrimitiveField        interface{}
		ArrayField            interface{}
		StructArrayField      interface{}
		MapField              map[string]interface{}
		StructField           interface{}
		NilField              interface{}
		InterfacePointerField *interface{}
	}

	var s interface{} = "World"
	expected := OuterStruct{
		PrimitiveField: "Hello",
		ArrayField:     []interface{}{int64(1), int64(2), int64(3)},
		StructField: map[string]interface{}{
			"Field1": int64(123),
			"Field2": "Field2",
		},
		MapField: map[string]interface{}{
			"MapField1": []interface{}{int64(4), int64(5), int64(6)},
			"MapField2": map[string]interface{}{
				"A": "A",
			},
			"MapField3": false,
		},
		NilField:              nil,
		InterfacePointerField: &s,
		StructArrayField: []map[string]interface{}{
			{
				"Name": "Allen",
				"Age":  int64(20),
			},
			{
				"Name": "Jack",
				"Age":  int64(23),
			},
		},
	}
	actual := OuterStruct{}
	if err := Unmarshal(toml, &actual); err == nil {
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Bad unmarshal: expected %v, got %v", expected, actual)
		}
	} else {
		t.Fatal(err)
	}
}

func TestUnmarshalToNonNilInterface(t *testing.T) {
	toml := []byte(`
PrimitiveField = "Allen"
ArrayField = [1,2,3]

[StructField]
InnerField = "After1"

[PointerField]
InnerField = "After2"

[InterfacePointerField]
InnerField = "After"

[MapField]
MapField1 = [4,5,6]
MapField2 = {A = "A"}
MapField3 = false

[[StructArrayField]]
InnerField = "After3"

[[StructArrayField]]
InnerField = "After4"
`)
	type InnerStruct struct {
		InnerField interface{}
	}

	type OuterStruct struct {
		PrimitiveField        interface{}
		ArrayField            interface{}
		StructArrayField      interface{}
		MapField              map[string]interface{}
		StructField           interface{}
		PointerField          interface{}
		NilField              interface{}
		InterfacePointerField *interface{}
	}

	var s interface{} = InnerStruct{"After"}
	expected := OuterStruct{
		PrimitiveField: "Allen",
		ArrayField:     []int{1, 2, 3},
		StructField:    InnerStruct{InnerField: "After1"},
		MapField: map[string]interface{}{
			"MapField1": []interface{}{int64(4), int64(5), int64(6)},
			"MapField2": map[string]interface{}{
				"A": "A",
			},
			"MapField3": false,
		},
		PointerField:          &InnerStruct{InnerField: "After2"},
		NilField:              nil,
		InterfacePointerField: &s,
		StructArrayField: []InnerStruct{
			{InnerField: "After3"},
			{InnerField: "After4"},
		},
	}
	actual := OuterStruct{
		PrimitiveField: "aaa",
		ArrayField:     []int{100, 200, 300, 400},
		StructField:    InnerStruct{InnerField: "Before1"},
		MapField: map[string]interface{}{
			"MapField1": []int{4, 5, 6},
			"MapField2": map[string]string{
				"B": "BBB",
			},
			"MapField3": true,
		},
		PointerField:          &InnerStruct{InnerField: "Before2"},
		NilField:              nil,
		InterfacePointerField: &s,
		StructArrayField: []InnerStruct{
			{InnerField: "Before3"},
			{InnerField: "Before4"},
		},
	}
	if err := Unmarshal(toml, &actual); err == nil {
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Bad unmarshal: expected %v, got %v", expected, actual)
		}
	} else {
		t.Fatal(err)
	}
}

func TestUnmarshalEmbedTree(t *testing.T) {
	toml := []byte(`
OuterField1 = "Out"
OuterField2 = 1024

[TreeField]
InnerField1 = "In"
InnerField2 = 2048

	[TreeField.EmbedStruct]
		EmbedField = "Embed"

`)
	type InnerStruct struct {
		InnerField1 string
		InnerField2 int
		EmbedStruct struct {
			EmbedField string
		}
	}

	type OuterStruct struct {
		OuterField1 string
		OuterField2 int
		TreeField   *Tree
	}

	out := OuterStruct{}
	actual := InnerStruct{}
	expected := InnerStruct{
		"In",
		2048,
		struct {
			EmbedField string
		}{
			EmbedField: "Embed",
		},
	}
	if err := Unmarshal(toml, &out); err != nil {
		t.Fatal(err)
	}
	if err := out.TreeField.Unmarshal(&actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Bad unmarshal: expected %v, got %v", expected, actual)
	}
}

func TestMarshalNil(t *testing.T) {
	if _, err := Marshal(nil); err == nil {
		t.Errorf("Expected err from nil marshal")
	}
	if _, err := Marshal((*struct{})(nil)); err == nil {
		t.Errorf("Expected err from nil marshal")
	}
}

func TestUnmarshalNil(t *testing.T) {
	if err := Unmarshal([]byte(`whatever = "whatever"`), nil); err == nil {
		t.Errorf("Expected err from nil marshal")
	}
	if err := Unmarshal([]byte(`whatever = "whatever"`), (*struct{})(nil)); err == nil {
		t.Errorf("Expected err from nil marshal")
	}
}

var sliceTomlDemo = []byte(`str_slice = ["Howdy","Hey There"]
str_slice_ptr= ["Howdy","Hey There"]
int_slice=[1,2]
int_slice_ptr=[1,2]
[[struct_slice]]
String2="1"
[[struct_slice]]
String2="2"
[[struct_slice_ptr]]
String2="1"
[[struct_slice_ptr]]
String2="2"
`)

type sliceStruct struct {
	Slice          []string                     `  toml:"str_slice"  `
	SlicePtr       *[]string                    `  toml:"str_slice_ptr"  `
	IntSlice       []int                        `  toml:"int_slice"  `
	IntSlicePtr    *[]int                       `  toml:"int_slice_ptr"  `
	StructSlice    []basicMarshalTestSubStruct  `  toml:"struct_slice"  `
	StructSlicePtr *[]basicMarshalTestSubStruct `  toml:"struct_slice_ptr"  `
}

type arrayStruct struct {
	Slice          [4]string                     `  toml:"str_slice"  `
	SlicePtr       *[4]string                    `  toml:"str_slice_ptr"  `
	IntSlice       [4]int                        `  toml:"int_slice"  `
	IntSlicePtr    *[4]int                       `  toml:"int_slice_ptr"  `
	StructSlice    [4]basicMarshalTestSubStruct  `  toml:"struct_slice"  `
	StructSlicePtr *[4]basicMarshalTestSubStruct `  toml:"struct_slice_ptr"  `
}

type arrayTooSmallStruct struct {
	Slice       [1]string                    `  toml:"str_slice"  `
	StructSlice [1]basicMarshalTestSubStruct `  toml:"struct_slice"  `
}

func TestUnmarshalSlice(t *testing.T) {
	tree, _ := LoadBytes(sliceTomlDemo)
	tree, _ = TreeFromMap(tree.ToMap())

	var actual sliceStruct
	err := tree.Unmarshal(&actual)
	if err != nil {
		t.Error("shound not err", err)
	}
	expected := sliceStruct{
		Slice:          []string{"Howdy", "Hey There"},
		SlicePtr:       &[]string{"Howdy", "Hey There"},
		IntSlice:       []int{1, 2},
		IntSlicePtr:    &[]int{1, 2},
		StructSlice:    []basicMarshalTestSubStruct{{"1"}, {"2"}},
		StructSlicePtr: &[]basicMarshalTestSubStruct{{"1"}, {"2"}},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Bad unmarshal: expected %v, got %v", expected, actual)
	}

}

func TestUnmarshalSliceFail(t *testing.T) {
	tree, _ := TreeFromMap(map[string]interface{}{
		"str_slice": []int{1, 2},
	})

	var actual sliceStruct
	err := tree.Unmarshal(&actual)
	if err.Error() != "(0, 0): Can't convert 1(int64) to string" {
		t.Error("expect err:(0, 0): Can't convert 1(int64) to string but got ", err)
	}
}

func TestUnmarshalSliceFail2(t *testing.T) {
	tree, _ := Load(`str_slice=[1,2]`)

	var actual sliceStruct
	err := tree.Unmarshal(&actual)
	if err.Error() != "(1, 1): Can't convert 1(int64) to string" {
		t.Error("expect err:(1, 1): Can't convert 1(int64) to string but got ", err)
	}

}

func TestMarshalMixedTypeArray(t *testing.T) {
	type InnerStruct struct {
		IntField int
		StrField string
	}

	type TestStruct struct {
		ArrayField []interface{}
	}

	expected := []byte(`ArrayField = [3.14, 100, true, "hello world", { IntField = 100, StrField = "inner1" }, [{ IntField = 200, StrField = "inner2" }, { IntField = 300, StrField = "inner3" }]]
`)

	if result, err := Marshal(TestStruct{
		ArrayField: []interface{}{
			3.14,
			100,
			true,
			"hello world",
			InnerStruct{
				IntField: 100,
				StrField: "inner1",
			},
			[]InnerStruct{
				{IntField: 200, StrField: "inner2"},
				{IntField: 300, StrField: "inner3"},
			},
		},
	}); err == nil {
		if !bytes.Equal(result, expected) {
			t.Errorf("Bad marshal: expected\n----\n%s\n----\ngot\n----\n%s\n----\n", expected, result)
		}
	} else {
		t.Fatal(err)
	}
}

func TestUnmarshalMixedTypeArray(t *testing.T) {
	type TestStruct struct {
		ArrayField []interface{}
	}

	toml := []byte(`ArrayField = [3.14,100,true,"hello world",{Field = "inner1"},[{Field = "inner2"},{Field = "inner3"}]]
`)

	actual := TestStruct{}
	expected := TestStruct{
		ArrayField: []interface{}{
			3.14,
			int64(100),
			true,
			"hello world",
			map[string]interface{}{
				"Field": "inner1",
			},
			[]map[string]interface{}{
				{"Field": "inner2"},
				{"Field": "inner3"},
			},
		},
	}

	if err := Unmarshal(toml, &actual); err == nil {
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Bad unmarshal: expected %#v, got %#v", expected, actual)
		}
	} else {
		t.Fatal(err)
	}
}

func TestUnmarshalArray(t *testing.T) {
	var tree *Tree
	var err error

	tree, _ = LoadBytes(sliceTomlDemo)
	var actual1 arrayStruct
	err = tree.Unmarshal(&actual1)
	if err != nil {
		t.Error("shound not err", err)
	}

	tree, _ = TreeFromMap(tree.ToMap())
	var actual2 arrayStruct
	err = tree.Unmarshal(&actual2)
	if err != nil {
		t.Error("shound not err", err)
	}

	expected := arrayStruct{
		Slice:          [4]string{"Howdy", "Hey There"},
		SlicePtr:       &[4]string{"Howdy", "Hey There"},
		IntSlice:       [4]int{1, 2},
		IntSlicePtr:    &[4]int{1, 2},
		StructSlice:    [4]basicMarshalTestSubStruct{{"1"}, {"2"}},
		StructSlicePtr: &[4]basicMarshalTestSubStruct{{"1"}, {"2"}},
	}
	if !reflect.DeepEqual(actual1, expected) {
		t.Errorf("Bad unmarshal: expected %v, got %v", expected, actual1)
	}
	if !reflect.DeepEqual(actual2, expected) {
		t.Errorf("Bad unmarshal: expected %v, got %v", expected, actual2)
	}
}

func TestUnmarshalArrayFail(t *testing.T) {
	tree, _ := TreeFromMap(map[string]interface{}{
		"str_slice": []string{"Howdy", "Hey There"},
	})

	var actual arrayTooSmallStruct
	err := tree.Unmarshal(&actual)
	if err.Error() != "(0, 0): unmarshal: TOML array length (2) exceeds destination array length (1)" {
		t.Error("expect err:(0, 0): unmarshal: TOML array length (2) exceeds destination array length (1) but got ", err)
	}
}

func TestUnmarshalArrayFail2(t *testing.T) {
	tree, _ := Load(`str_slice=["Howdy","Hey There"]`)

	var actual arrayTooSmallStruct
	err := tree.Unmarshal(&actual)
	if err.Error() != "(1, 1): unmarshal: TOML array length (2) exceeds destination array length (1)" {
		t.Error("expect err:(1, 1): unmarshal: TOML array length (2) exceeds destination array length (1) but got ", err)
	}
}

func TestUnmarshalArrayFail3(t *testing.T) {
	tree, _ := Load(`[[struct_slice]]
String2="1"
[[struct_slice]]
String2="2"`)

	var actual arrayTooSmallStruct
	err := tree.Unmarshal(&actual)
	if err.Error() != "(3, 1): unmarshal: TOML array length (2) exceeds destination array length (1)" {
		t.Error("expect err:(3, 1): unmarshal: TOML array length (2) exceeds destination array length (1) but got ", err)
	}
}

func TestDecoderStrict(t *testing.T) {
	input := `
[decoded]
  key = ""

[undecoded]
  key = ""

  [undecoded.inner]
	key = ""

  [[undecoded.array]]
	key = ""

  [[undecoded.array]]
	key = ""

`
	var doc struct {
		Decoded struct {
			Key string
		}
	}

	expected := `undecoded keys: ["undecoded.array.0.key" "undecoded.array.1.key" "undecoded.inner.key" "undecoded.key"]`

	err := NewDecoder(bytes.NewReader([]byte(input))).Strict(true).Decode(&doc)
	if err == nil {
		t.Error("expected error, got none")
	} else if err.Error() != expected {
		t.Errorf("expect err: %s, got: %s", expected, err.Error())
	}

	if err := NewDecoder(bytes.NewReader([]byte(input))).Decode(&doc); err != nil {
		t.Errorf("unexpected err: %s", err)
	}

	var m map[string]interface{}
	if err := NewDecoder(bytes.NewReader([]byte(input))).Decode(&m); err != nil {
		t.Errorf("unexpected err: %s", err)
	}
}

func TestDecoderStrictValid(t *testing.T) {
	input := `
[decoded]
  key = ""
`
	var doc struct {
		Decoded struct {
			Key string
		}
	}

	err := NewDecoder(bytes.NewReader([]byte(input))).Strict(true).Decode(&doc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

type docUnmarshalTOML struct {
	Decoded struct {
		Key string
	}
}

func (d *docUnmarshalTOML) UnmarshalTOML(i interface{}) error {
	if iMap, ok := i.(map[string]interface{}); !ok {
		return fmt.Errorf("type assertion error: wants %T, have %T", map[string]interface{}{}, i)
	} else if key, ok := iMap["key"]; !ok {
		return fmt.Errorf("key '%s' not in map", "key")
	} else if keyString, ok := key.(string); !ok {
		return fmt.Errorf("type assertion error: wants %T, have %T", "", key)
	} else {
		d.Decoded.Key = keyString
	}
	return nil
}

func TestDecoderStrictCustomUnmarshal(t *testing.T) {
	input := `key = "ok"`
	var doc docUnmarshalTOML
	err := NewDecoder(bytes.NewReader([]byte(input))).Strict(true).Decode(&doc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if doc.Decoded.Key != "ok" {
		t.Errorf("Bad unmarshal: expected ok, got %v", doc.Decoded.Key)
	}
}

type parent struct {
	Doc        docUnmarshalTOML
	DocPointer *docUnmarshalTOML
}

func TestCustomUnmarshal(t *testing.T) {
	input := `
[Doc]
    key = "ok1"
[DocPointer]
    key = "ok2"
`

	var d parent
	if err := Unmarshal([]byte(input), &d); err != nil {
		t.Fatalf("unexpected err: %s", err.Error())
	}
	if d.Doc.Decoded.Key != "ok1" {
		t.Errorf("Bad unmarshal: expected ok, got %v", d.Doc.Decoded.Key)
	}
	if d.DocPointer.Decoded.Key != "ok2" {
		t.Errorf("Bad unmarshal: expected ok, got %v", d.DocPointer.Decoded.Key)
	}
}

func TestCustomUnmarshalError(t *testing.T) {
	input := `
[Doc]
    key = 1
[DocPointer]
    key = "ok2"
`

	expected := "(2, 1): unmarshal toml: type assertion error: wants string, have int64"

	var d parent
	err := Unmarshal([]byte(input), &d)
	if err == nil {
		t.Error("expected error, got none")
	} else if err.Error() != expected {
		t.Errorf("expect err: %s, got: %s", expected, err.Error())
	}
}

type intWrapper struct {
	Value int
}

func (w *intWrapper) UnmarshalText(text []byte) error {
	var err error
	if w.Value, err = strconv.Atoi(string(text)); err == nil {
		return nil
	}
	if b, err := strconv.ParseBool(string(text)); err == nil {
		if b {
			w.Value = 1
		}
		return nil
	}
	if f, err := strconv.ParseFloat(string(text), 32); err == nil {
		w.Value = int(f)
		return nil
	}
	return fmt.Errorf("unsupported: %s", text)
}

func TestTextUnmarshal(t *testing.T) {
	var doc struct {
		UnixTime intWrapper
		Version  *intWrapper

		Bool  intWrapper
		Int   intWrapper
		Float intWrapper
	}

	input := `
UnixTime = "12"
Version = "42"
Bool = true
Int = 21
Float = 2.0
`

	if err := Unmarshal([]byte(input), &doc); err != nil {
		t.Fatalf("unexpected err: %s", err.Error())
	}
	if doc.UnixTime.Value != 12 {
		t.Fatalf("expected UnixTime: 12 got: %d", doc.UnixTime.Value)
	}
	if doc.Version.Value != 42 {
		t.Fatalf("expected Version: 42 got: %d", doc.Version.Value)
	}
	if doc.Bool.Value != 1 {
		t.Fatalf("expected Bool: 1 got: %d", doc.Bool.Value)
	}
	if doc.Int.Value != 21 {
		t.Fatalf("expected Int: 21 got: %d", doc.Int.Value)
	}
	if doc.Float.Value != 2 {
		t.Fatalf("expected Float: 2 got: %d", doc.Float.Value)
	}
}

func TestTextUnmarshalError(t *testing.T) {
	var doc struct {
		Failer intWrapper
	}

	input := `Failer = "hello"`
	if err := Unmarshal([]byte(input), &doc); err == nil {
		t.Fatalf("expected err, got none")
	}
}

// issue406
func TestPreserveNotEmptyField(t *testing.T) {
	toml := []byte(`Field1 = "ccc"`)
	type Inner struct {
		InnerField1 string
		InnerField2 int
	}
	type TestStruct struct {
		Field1 string
		Field2 int
		Field3 Inner
	}

	actual := TestStruct{
		"aaa",
		100,
		Inner{
			"bbb",
			200,
		},
	}

	expected := TestStruct{
		"ccc",
		100,
		Inner{
			"bbb",
			200,
		},
	}

	err := Unmarshal(toml, &actual)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Bad unmarshal: expected %+v, got %+v", expected, actual)
	}
}
