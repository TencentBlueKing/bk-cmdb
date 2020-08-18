// Copyright 2014 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package ini_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"gopkg.in/ini.v1"
)

type testNested struct {
	Cities      []string `delim:"|"`
	Visits      []time.Time
	Years       []int
	Numbers     []int64
	Ages        []uint
	Populations []uint64
	Coordinates []float64
	Flags       []bool
	Note        string
	Unused      int `ini:"-"`
}

type TestEmbeded struct {
	GPA float64
}

type testStruct struct {
	Name           string `ini:"NAME"`
	Age            int
	Male           bool
	Money          float64
	Born           time.Time
	Time           time.Duration `ini:"Duration"`
	OldVersionTime time.Duration
	Others         testNested
	OthersPtr      *testNested
	NilPtr         *testNested
	*TestEmbeded   `ini:"grade"`
	Unused         int `ini:"-"`
	Unsigned       uint
	Omitted        bool     `ini:"omitthis,omitempty"`
	Shadows        []string `ini:",,allowshadow"`
	ShadowInts     []int    `ini:"Shadows,,allowshadow"`
	BoolPtr        *bool
	BoolPtrNil     *bool
	FloatPtr       *float64
	FloatPtrNil    *float64
	IntPtr         *int
	IntPtrNil      *int
	UintPtr        *uint
	UintPtrNil     *uint
	StringPtr      *string
	StringPtrNil   *string
	TimePtr        *time.Time
	TimePtrNil     *time.Time
	DurationPtr    *time.Duration
	DurationPtrNil *time.Duration
}

type testInterface struct {
	Address    string
	ListenPort int
	PrivateKey string
}

type testPeer struct {
	PublicKey    string
	PresharedKey string
	AllowedIPs   []string `delim:","`
}

type testNonUniqueSectionsStruct struct {
	Interface testInterface
	Peer      []testPeer `ini:",,,nonunique"`
}

const confDataStruct = `
NAME = Unknwon
Age = 21
Male = true
Money = 1.25
Born = 1993-10-07T20:17:05Z
Duration = 2h45m
OldVersionTime = 30
Unsigned = 3
omitthis = true
Shadows = 1, 2
Shadows = 3, 4
BoolPtr = false
FloatPtr = 0
IntPtr = 0
UintPtr = 0
StringPtr = ""
TimePtr = 0001-01-01T00:00:00Z
DurationPtr = 0s

[Others]
Cities = HangZhou|Boston
Visits = 1993-10-07T20:17:05Z, 1993-10-07T20:17:05Z
Years = 1993,1994
Numbers = 10010,10086
Ages = 18,19
Populations = 12345678,98765432
Coordinates = 192.168,10.11
Flags       = true,false
Note = Hello world!

[OthersPtr]
Cities = HangZhou|Boston
Visits = 1993-10-07T20:17:05Z, 1993-10-07T20:17:05Z
Years = 1993,1994
Numbers = 10010,10086
Ages = 18,19
Populations = 12345678,98765432
Coordinates = 192.168,10.11
Flags       = true,false
Note = Hello world!

[grade]
GPA = 2.8

[foo.bar]
Here = there
When = then
`

const confNonUniqueSectionDataStruct = `[Interface]
Address    = 10.2.0.1/24
ListenPort = 34777
PrivateKey = privServerKey

[Peer]
PublicKey    = pubClientKey
PresharedKey = psKey
AllowedIPs   = 10.2.0.2/32,fd00:2::2/128

[Peer]
PublicKey    = pubClientKey2
PresharedKey = psKey2
AllowedIPs   = 10.2.0.3/32,fd00:2::3/128

`

type unsupport struct {
	Byte byte
}

type unsupport2 struct {
	Others struct {
		Cities byte
	}
}

type Unsupport3 struct {
	Cities byte
}

type unsupport4 struct {
	*Unsupport3 `ini:"Others"`
}

type defaultValue struct {
	Name     string
	Age      int
	Male     bool
	Optional *bool
	Money    float64
	Born     time.Time
	Cities   []string
}

type fooBar struct {
	Here, When string
}

const invalidDataConfStruct = `
Name = 
Age = age
Male = 123
Money = money
Born = nil
Cities = 
`

func Test_MapToStruct(t *testing.T) {
	Convey("Map to struct", t, func() {
		Convey("Map file to struct", func() {
			ts := new(testStruct)
			So(ini.MapTo(ts, []byte(confDataStruct)), ShouldBeNil)

			So(ts.Name, ShouldEqual, "Unknwon")
			So(ts.Age, ShouldEqual, 21)
			So(ts.Male, ShouldBeTrue)
			So(ts.Money, ShouldEqual, 1.25)
			So(ts.Unsigned, ShouldEqual, 3)

			t, err := time.Parse(time.RFC3339, "1993-10-07T20:17:05Z")
			So(err, ShouldBeNil)
			So(ts.Born.String(), ShouldEqual, t.String())

			dur, err := time.ParseDuration("2h45m")
			So(err, ShouldBeNil)
			So(ts.Time.Seconds(), ShouldEqual, dur.Seconds())

			So(ts.OldVersionTime*time.Second, ShouldEqual, 30*time.Second)

			So(strings.Join(ts.Others.Cities, ","), ShouldEqual, "HangZhou,Boston")
			So(ts.Others.Visits[0].String(), ShouldEqual, t.String())
			So(fmt.Sprint(ts.Others.Years), ShouldEqual, "[1993 1994]")
			So(fmt.Sprint(ts.Others.Numbers), ShouldEqual, "[10010 10086]")
			So(fmt.Sprint(ts.Others.Ages), ShouldEqual, "[18 19]")
			So(fmt.Sprint(ts.Others.Populations), ShouldEqual, "[12345678 98765432]")
			So(fmt.Sprint(ts.Others.Coordinates), ShouldEqual, "[192.168 10.11]")
			So(fmt.Sprint(ts.Others.Flags), ShouldEqual, "[true false]")
			So(ts.Others.Note, ShouldEqual, "Hello world!")
			So(ts.TestEmbeded.GPA, ShouldEqual, 2.8)

			So(strings.Join(ts.OthersPtr.Cities, ","), ShouldEqual, "HangZhou,Boston")
			So(ts.OthersPtr.Visits[0].String(), ShouldEqual, t.String())
			So(fmt.Sprint(ts.OthersPtr.Years), ShouldEqual, "[1993 1994]")
			So(fmt.Sprint(ts.OthersPtr.Numbers), ShouldEqual, "[10010 10086]")
			So(fmt.Sprint(ts.OthersPtr.Ages), ShouldEqual, "[18 19]")
			So(fmt.Sprint(ts.OthersPtr.Populations), ShouldEqual, "[12345678 98765432]")
			So(fmt.Sprint(ts.OthersPtr.Coordinates), ShouldEqual, "[192.168 10.11]")
			So(fmt.Sprint(ts.OthersPtr.Flags), ShouldEqual, "[true false]")
			So(ts.OthersPtr.Note, ShouldEqual, "Hello world!")

			So(ts.NilPtr, ShouldBeNil)

			So(*ts.BoolPtr, ShouldEqual, false)
			So(ts.BoolPtrNil, ShouldEqual, nil)
			So(*ts.FloatPtr, ShouldEqual, 0)
			So(ts.FloatPtrNil, ShouldEqual, nil)
			So(*ts.IntPtr, ShouldEqual, 0)
			So(ts.IntPtrNil, ShouldEqual, nil)
			So(*ts.UintPtr, ShouldEqual, 0)
			So(ts.UintPtrNil, ShouldEqual, nil)
			So(*ts.StringPtr, ShouldEqual, "")
			So(ts.StringPtrNil, ShouldEqual, nil)
			So(*ts.TimePtr, ShouldNotEqual, nil)
			So(ts.TimePtrNil, ShouldEqual, nil)
			So(*ts.DurationPtr, ShouldEqual, 0)
			So(ts.DurationPtrNil, ShouldEqual, nil)
		})

		Convey("Map section to struct", func() {
			foobar := new(fooBar)
			f, err := ini.Load([]byte(confDataStruct))
			So(err, ShouldBeNil)

			So(f.Section("foo.bar").MapTo(foobar), ShouldBeNil)
			So(foobar.Here, ShouldEqual, "there")
			So(foobar.When, ShouldEqual, "then")
		})

		Convey("Map to non-pointer struct", func() {
			f, err := ini.Load([]byte(confDataStruct))
			So(err, ShouldBeNil)
			So(f, ShouldNotBeNil)

			So(f.MapTo(testStruct{}), ShouldNotBeNil)
		})

		Convey("Map to unsupported type", func() {
			f, err := ini.Load([]byte(confDataStruct))
			So(err, ShouldBeNil)
			So(f, ShouldNotBeNil)

			f.NameMapper = func(raw string) string {
				if raw == "Byte" {
					return "NAME"
				}
				return raw
			}
			So(f.MapTo(&unsupport{}), ShouldNotBeNil)
			So(f.MapTo(&unsupport2{}), ShouldNotBeNil)
			So(f.MapTo(&unsupport4{}), ShouldNotBeNil)
		})

		Convey("Map to omitempty field", func() {
			ts := new(testStruct)
			So(ini.MapTo(ts, []byte(confDataStruct)), ShouldBeNil)

			So(ts.Omitted, ShouldEqual, true)
		})

		Convey("Map with shadows", func() {
			f, err := ini.LoadSources(ini.LoadOptions{AllowShadows: true}, []byte(confDataStruct))
			So(err, ShouldBeNil)
			ts := new(testStruct)
			So(f.MapTo(ts), ShouldBeNil)

			So(strings.Join(ts.Shadows, " "), ShouldEqual, "1 2 3 4")
			So(fmt.Sprintf("%v", ts.ShadowInts), ShouldEqual, "[1 2 3 4]")
		})

		Convey("Map from invalid data source", func() {
			So(ini.MapTo(&testStruct{}, "hi"), ShouldNotBeNil)
		})

		Convey("Map to wrong types and gain default values", func() {
			f, err := ini.Load([]byte(invalidDataConfStruct))
			So(err, ShouldBeNil)

			t, err := time.Parse(time.RFC3339, "1993-10-07T20:17:05Z")
			So(err, ShouldBeNil)
			dv := &defaultValue{"Joe", 10, true, nil, 1.25, t, []string{"HangZhou", "Boston"}}
			So(f.MapTo(dv), ShouldBeNil)
			So(dv.Name, ShouldEqual, "Joe")
			So(dv.Age, ShouldEqual, 10)
			So(dv.Male, ShouldBeTrue)
			So(dv.Money, ShouldEqual, 1.25)
			So(dv.Born.String(), ShouldEqual, t.String())
			So(strings.Join(dv.Cities, ","), ShouldEqual, "HangZhou,Boston")
		})
	})

	Convey("Map to struct in strict mode", t, func() {
		f, err := ini.Load([]byte(`
name=bruce
age=a30`))
		So(err, ShouldBeNil)

		type Strict struct {
			Name string `ini:"name"`
			Age  int    `ini:"age"`
		}
		s := new(Strict)

		So(f.Section("").StrictMapTo(s), ShouldNotBeNil)
	})

	Convey("Map slice in strict mode", t, func() {
		f, err := ini.Load([]byte(`
names=alice, bruce`))
		So(err, ShouldBeNil)

		type Strict struct {
			Names []string `ini:"names"`
		}
		s := new(Strict)

		So(f.Section("").StrictMapTo(s), ShouldBeNil)
		So(fmt.Sprint(s.Names), ShouldEqual, "[alice bruce]")
	})
}

func Test_MapToStructNonUniqueSections(t *testing.T) {
	Convey("Map to struct non unique", t, func() {
		Convey("Map file to struct non unique", func() {
			f, err := ini.LoadSources(ini.LoadOptions{AllowNonUniqueSections: true}, []byte(confNonUniqueSectionDataStruct))
			So(err, ShouldBeNil)
			ts := new(testNonUniqueSectionsStruct)

			So(f.MapTo(ts), ShouldBeNil)

			So(ts.Interface.Address, ShouldEqual, "10.2.0.1/24")
			So(ts.Interface.ListenPort, ShouldEqual, 34777)
			So(ts.Interface.PrivateKey, ShouldEqual, "privServerKey")

			So(ts.Peer[0].PublicKey, ShouldEqual, "pubClientKey")
			So(ts.Peer[0].PresharedKey, ShouldEqual, "psKey")
			So(ts.Peer[0].AllowedIPs[0], ShouldEqual, "10.2.0.2/32")
			So(ts.Peer[0].AllowedIPs[1], ShouldEqual, "fd00:2::2/128")

			So(ts.Peer[1].PublicKey, ShouldEqual, "pubClientKey2")
			So(ts.Peer[1].PresharedKey, ShouldEqual, "psKey2")
			So(ts.Peer[1].AllowedIPs[0], ShouldEqual, "10.2.0.3/32")
			So(ts.Peer[1].AllowedIPs[1], ShouldEqual, "fd00:2::3/128")
		})

		Convey("Map non unique section to struct", func() {
			newPeer := new(testPeer)
			newPeerSlice := make([]testPeer, 0)

			f, err := ini.LoadSources(ini.LoadOptions{AllowNonUniqueSections: true}, []byte(confNonUniqueSectionDataStruct))
			So(err, ShouldBeNil)

			// try only first one
			So(f.Section("Peer").MapTo(newPeer), ShouldBeNil)
			So(newPeer.PublicKey, ShouldEqual, "pubClientKey")
			So(newPeer.PresharedKey, ShouldEqual, "psKey")
			So(newPeer.AllowedIPs[0], ShouldEqual, "10.2.0.2/32")
			So(newPeer.AllowedIPs[1], ShouldEqual, "fd00:2::2/128")

			// try all
			So(f.Section("Peer").MapTo(&newPeerSlice), ShouldBeNil)
			So(newPeerSlice[0].PublicKey, ShouldEqual, "pubClientKey")
			So(newPeerSlice[0].PresharedKey, ShouldEqual, "psKey")
			So(newPeerSlice[0].AllowedIPs[0], ShouldEqual, "10.2.0.2/32")
			So(newPeerSlice[0].AllowedIPs[1], ShouldEqual, "fd00:2::2/128")

			So(newPeerSlice[1].PublicKey, ShouldEqual, "pubClientKey2")
			So(newPeerSlice[1].PresharedKey, ShouldEqual, "psKey2")
			So(newPeerSlice[1].AllowedIPs[0], ShouldEqual, "10.2.0.3/32")
			So(newPeerSlice[1].AllowedIPs[1], ShouldEqual, "fd00:2::3/128")
		})

		Convey("Map non unique sections with subsections to struct", func() {
			iniFile, err := ini.LoadSources(ini.LoadOptions{AllowNonUniqueSections: true}, strings.NewReader(`
[Section]
FieldInSubSection = 1
FieldInSubSection2 = 2
FieldInSection = 3

[Section]
FieldInSubSection = 4
FieldInSubSection2 = 5
FieldInSection = 6
`))
			So(err, ShouldBeNil)

			type SubSection struct {
				FieldInSubSection string `ini:"FieldInSubSection"`
			}
			type SubSection2 struct {
				FieldInSubSection2 string `ini:"FieldInSubSection2"`
			}

			type Section struct {
				SubSection     `ini:"Section"`
				SubSection2    `ini:"Section"`
				FieldInSection string `ini:"FieldInSection"`
			}

			type File struct {
				Sections []Section `ini:"Section,,,nonunique"`
			}

			f := new(File)
			err = iniFile.MapTo(f)
			So(err, ShouldBeNil)

			So(f.Sections[0].FieldInSubSection, ShouldEqual, "1")
			So(f.Sections[0].FieldInSubSection2, ShouldEqual, "2")
			So(f.Sections[0].FieldInSection, ShouldEqual, "3")

			So(f.Sections[1].FieldInSubSection, ShouldEqual, "4")
			So(f.Sections[1].FieldInSubSection2, ShouldEqual, "5")
			So(f.Sections[1].FieldInSection, ShouldEqual, "6")
		})
	})
}

func Test_ReflectFromStruct(t *testing.T) {
	Convey("Reflect from struct", t, func() {
		type Embeded struct {
			Dates       []time.Time `delim:"|" comment:"Time data"`
			Places      []string
			Years       []int
			Numbers     []int64
			Ages        []uint
			Populations []uint64
			Coordinates []float64
			Flags       []bool
			None        []int
		}
		type Author struct {
			Name      string `ini:"NAME"`
			Male      bool
			Optional  *bool
			Age       int `comment:"Author's age"`
			Height    uint
			GPA       float64
			Date      time.Time
			NeverMind string `ini:"-"`
			ignored   string
			*Embeded  `ini:"infos" comment:"Embeded section"`
		}

		t, err := time.Parse(time.RFC3339, "1993-10-07T20:17:05Z")
		So(err, ShouldBeNil)
		a := &Author{"Unknwon", true, nil, 21, 100, 2.8, t, "", "ignored",
			&Embeded{
				[]time.Time{t, t},
				[]string{"HangZhou", "Boston"},
				[]int{1993, 1994},
				[]int64{10010, 10086},
				[]uint{18, 19},
				[]uint64{12345678, 98765432},
				[]float64{192.168, 10.11},
				[]bool{true, false},
				[]int{},
			}}
		cfg := ini.Empty()
		So(ini.ReflectFrom(cfg, a), ShouldBeNil)

		var buf bytes.Buffer
		_, err = cfg.WriteTo(&buf)
		So(err, ShouldBeNil)
		So(buf.String(), ShouldEqual, `NAME     = Unknwon
Male     = true
Optional = 
; Author's age
Age      = 21
Height   = 100
GPA      = 2.8
Date     = 1993-10-07T20:17:05Z

; Embeded section
[infos]
; Time data
Dates       = 1993-10-07T20:17:05Z|1993-10-07T20:17:05Z
Places      = HangZhou,Boston
Years       = 1993,1994
Numbers     = 10010,10086
Ages        = 18,19
Populations = 12345678,98765432
Coordinates = 192.168,10.11
Flags       = true,false
None        = 

`)

		Convey("Reflect from non-point struct", func() {
			So(ini.ReflectFrom(cfg, Author{}), ShouldNotBeNil)
		})

		Convey("Reflect from struct with omitempty", func() {
			cfg := ini.Empty()
			type SpecialStruct struct {
				FirstName  string    `ini:"first_name"`
				LastName   string    `ini:"last_name"`
				JustOmitMe string    `ini:"omitempty"`
				LastLogin  time.Time `ini:"last_login,omitempty"`
				LastLogin2 time.Time `ini:",omitempty"`
				NotEmpty   int       `ini:"omitempty"`
			}

			So(ini.ReflectFrom(cfg, &SpecialStruct{FirstName: "John", LastName: "Doe", NotEmpty: 9}), ShouldBeNil)

			var buf bytes.Buffer
			_, err = cfg.WriteTo(&buf)
			So(buf.String(), ShouldEqual, `first_name = John
last_name  = Doe
omitempty  = 9

`)
		})
	})
}

func Test_ReflectFromStructNonUniqueSections(t *testing.T) {
	Convey("Reflect from struct with non unique sections", t, func() {
		nonUnique := &testNonUniqueSectionsStruct{
			Interface: testInterface{
				Address:    "10.2.0.1/24",
				ListenPort: 34777,
				PrivateKey: "privServerKey",
			},
			Peer: []testPeer{
				{
					PublicKey:    "pubClientKey",
					PresharedKey: "psKey",
					AllowedIPs:   []string{"10.2.0.2/32,fd00:2::2/128"},
				},
				{
					PublicKey:    "pubClientKey2",
					PresharedKey: "psKey2",
					AllowedIPs:   []string{"10.2.0.3/32,fd00:2::3/128"},
				},
			},
		}

		cfg := ini.Empty(ini.LoadOptions{
			AllowNonUniqueSections: true,
		})

		So(ini.ReflectFrom(cfg, nonUnique), ShouldBeNil)

		var buf bytes.Buffer
		_, err := cfg.WriteTo(&buf)
		So(err, ShouldBeNil)
		So(buf.String(), ShouldEqual, confNonUniqueSectionDataStruct)

		// note: using ReflectFrom from should overwrite the existing sections
		err = cfg.Section("Peer").ReflectFrom([]*testPeer{
			{
				PublicKey:    "pubClientKey3",
				PresharedKey: "psKey3",
				AllowedIPs:   []string{"10.2.0.4/32,fd00:2::4/128"},
			},
			{
				PublicKey:    "pubClientKey4",
				PresharedKey: "psKey4",
				AllowedIPs:   []string{"10.2.0.5/32,fd00:2::5/128"},
			},
		})

		So(err, ShouldBeNil)

		buf = bytes.Buffer{}
		_, err = cfg.WriteTo(&buf)
		So(err, ShouldBeNil)
		So(buf.String(), ShouldEqual, `[Interface]
Address    = 10.2.0.1/24
ListenPort = 34777
PrivateKey = privServerKey

[Peer]
PublicKey    = pubClientKey3
PresharedKey = psKey3
AllowedIPs   = 10.2.0.4/32,fd00:2::4/128

[Peer]
PublicKey    = pubClientKey4
PresharedKey = psKey4
AllowedIPs   = 10.2.0.5/32,fd00:2::5/128

`)

		// note: using ReflectFrom from should overwrite the existing sections
		err = cfg.Section("Peer").ReflectFrom(&testPeer{
			PublicKey:    "pubClientKey5",
			PresharedKey: "psKey5",
			AllowedIPs:   []string{"10.2.0.6/32,fd00:2::6/128"},
		})

		So(err, ShouldBeNil)

		buf = bytes.Buffer{}
		_, err = cfg.WriteTo(&buf)
		So(err, ShouldBeNil)
		So(buf.String(), ShouldEqual, `[Interface]
Address    = 10.2.0.1/24
ListenPort = 34777
PrivateKey = privServerKey

[Peer]
PublicKey    = pubClientKey5
PresharedKey = psKey5
AllowedIPs   = 10.2.0.6/32,fd00:2::6/128

`)
	})
}

// Inspired by https://github.com/go-ini/ini/issues/196
func TestMapToAndReflectFromStructWithShadows(t *testing.T) {
	Convey("Map to struct and then reflect with shadows should generate original config content", t, func() {
		type include struct {
			Paths []string `ini:"path,omitempty,allowshadow"`
		}

		cfg, err := ini.LoadSources(ini.LoadOptions{
			AllowShadows: true,
		}, []byte(`
[include]
path = /tmp/gpm-profiles/test5.profile
path = /tmp/gpm-profiles/test1.profile`))
		So(err, ShouldBeNil)

		sec := cfg.Section("include")
		inc := new(include)
		err = sec.MapTo(inc)
		So(err, ShouldBeNil)

		err = sec.ReflectFrom(inc)
		So(err, ShouldBeNil)

		var buf bytes.Buffer
		_, err = cfg.WriteTo(&buf)
		So(err, ShouldBeNil)
		So(buf.String(), ShouldEqual, `[include]
path = /tmp/gpm-profiles/test5.profile
path = /tmp/gpm-profiles/test1.profile

`)
	})
}

type testMapper struct {
	PackageName string
}

func Test_NameGetter(t *testing.T) {
	Convey("Test name mappers", t, func() {
		So(ini.MapToWithMapper(&testMapper{}, ini.TitleUnderscore, []byte("packag_name=ini")), ShouldBeNil)

		cfg, err := ini.Load([]byte("PACKAGE_NAME=ini"))
		So(err, ShouldBeNil)
		So(cfg, ShouldNotBeNil)

		cfg.NameMapper = ini.SnackCase
		tg := new(testMapper)
		So(cfg.MapTo(tg), ShouldBeNil)
		So(tg.PackageName, ShouldEqual, "ini")
	})
}

type testDurationStruct struct {
	Duration time.Duration `ini:"Duration"`
}

func Test_Duration(t *testing.T) {
	Convey("Duration less than 16m50s", t, func() {
		ds := new(testDurationStruct)
		So(ini.MapTo(ds, []byte("Duration=16m49s")), ShouldBeNil)

		dur, err := time.ParseDuration("16m49s")
		So(err, ShouldBeNil)
		So(ds.Duration.Seconds(), ShouldEqual, dur.Seconds())
	})
}

type Employer struct {
	Name  string
	Title string
}

type Employers []*Employer

func (es Employers) ReflectINIStruct(f *ini.File) error {
	for _, e := range es {
		f.Section(e.Name).Key("Title").SetValue(e.Title)
	}
	return nil
}

// Inspired by https://github.com/go-ini/ini/issues/199
func Test_StructReflector(t *testing.T) {
	Convey("Reflect with StructReflector interface", t, func() {
		p := &struct {
			FirstName string
			Employer  Employers
		}{
			FirstName: "Andrew",
			Employer: []*Employer{
				{
					Name:  `Employer "VMware"`,
					Title: "Staff II Engineer",
				},
				{
					Name:  `Employer "EMC"`,
					Title: "Consultant Engineer",
				},
			},
		}

		f := ini.Empty()
		So(f.ReflectFrom(p), ShouldBeNil)

		var buf bytes.Buffer
		_, err := f.WriteTo(&buf)
		So(err, ShouldBeNil)

		So(buf.String(), ShouldEqual, `FirstName = Andrew

[Employer "VMware"]
Title = Staff II Engineer

[Employer "EMC"]
Title = Consultant Engineer

`)
	})
}
