package metadata

import (
	"testing"
)

func TestSearchSort(t *testing.T) {
	type testData struct {
		input string
		ss    []SearchSort
	}
	testDataArr := []testData{
		testData{
			input: "aa,bb",
			ss: []SearchSort{
				SearchSort{
					Field: "aa",
				},
				SearchSort{
					Field: "bb",
				},
			},
		},
		testData{
			input: "aa:-1,bb:1,cc,dd:abc",
			ss: []SearchSort{
				SearchSort{
					Field: "aa",
					IsDsc: true,
				},
				SearchSort{
					Field: "bb",
				},
				SearchSort{
					Field: "cc",
				},
				SearchSort{
					Field: "dd",
				},
			},
		},
	}
	for _, testDataItem := range testDataArr {
		testSSArr := NewSearchSortParse().Str(testDataItem.input)
		if len(testSSArr) != len(testDataItem.ss) {
			t.Errorf("str parse to Search Sort error!")
			return
		}
		for idx, ssItem := range testSSArr {
			if ssItem.Field != testDataItem.ss[idx].Field ||
				ssItem.IsDsc != testDataItem.ss[idx].IsDsc {
				t.Errorf("%+v, %+v not equal", ssItem, testDataItem.ss[idx])
				return
			}
		}
	}

}
