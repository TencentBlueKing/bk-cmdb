package logics

import (
	//"math"
	//"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"configcenter/src/common"
)

func TestCondtion(t *testing.T) {
	type testData struct {
		input           string
		isInteger       bool
		output_condtion common.KvMap
		needExtCompare  bool
		isErr           bool
		isNil           bool
	}

	td := []testData{
		testData{
			input:           "jo?",
			output_condtion: common.KvMap{common.BKDBLIKE: "^jo.$"},
		},
		testData{
			input:           "?ppt",
			output_condtion: common.KvMap{common.BKDBLIKE: "^.ppt$"},
		},
		testData{
			input:           "ap?t",
			output_condtion: common.KvMap{common.BKDBLIKE: "^ap.t$"},
		},
		testData{
			input:           "app?",
			output_condtion: common.KvMap{common.BKDBLIKE: "^app.$"},
		},
		testData{
			input:           "11?1",
			isInteger:       true,
			output_condtion: common.KvMap{common.BKDBIN: []int64{1101, 1111, 1121, 1131, 1141, 1151, 1161, 1171, 1181, 1191}}, //"^app.$",
		},
		testData{
			input:           "?11",
			isInteger:       true,
			output_condtion: common.KvMap{common.BKDBIN: []int64{111, 211, 311, 411, 511, 611, 711, 811, 911}}, //"^app.$",
		},
		testData{
			input:           "11?",
			isInteger:       true,
			output_condtion: common.KvMap{common.BKDBIN: []int64{110, 111, 112, 113, 114, 115, 116, 117, 118, 119}}, //"^app.$",
		},
		testData{
			input:           "11[2,3,10]",
			isInteger:       true,
			output_condtion: common.KvMap{common.BKDBIN: []int64{112, 113, 1110}}, //"^app.$",
		},
		testData{
			input:           "app[o,t]",
			output_condtion: common.KvMap{common.BKDBIN: []string{"appo", "appt"}}, //"^app.$",
		},
		testData{
			input:           "aaa[1-10]",
			output_condtion: nil, //"^app.$",
			needExtCompare:  true,
		},
		testData{
			input:           "aaa[1,4-3,a]",
			output_condtion: nil, //"^app.$",
			needExtCompare:  true,
		},
		testData{
			input:           "aaa?1111",
			output_condtion: nil, //"^app.$",
			isErr:           true,
			isInteger:       true,
		},
		testData{
			input:           "1?1a11",
			output_condtion: nil, //"^app.$",
			isErr:           true,
			isInteger:       true,
		},
		testData{
			input:           "11[1,2,a]",
			output_condtion: nil, //"^app.$",
			isErr:           true,
			isInteger:       true,
		},
		testData{
			input:           "11[a,2,3]",
			output_condtion: nil, //"^app.$",
			isErr:           true,
			isInteger:       true,
		},
		testData{
			input:           "11[a,1,2]",
			output_condtion: nil, //"^app.$",
			isErr:           true,
			isInteger:       true,
		},
	}

	for _, item := range td {
		sm := NewScopeMatch(item.input, !item.isInteger)

		cond, err := sm.ParseConds()
		if item.isErr {
			if nil == err {
				t.Errorf("test item %v need error", item)
				continue
			} else {
				continue
			}
		}
		if nil != err {
			t.Errorf("test item %v error:%s", item, err.Error())
			continue
		}
		if sm.needExtCompare != item.needExtCompare {
			t.Errorf("test not parse return error, need %v not %v, role:%s", item.needExtCompare, sm.needExtCompare, item.input)
			continue
		}
		//a, _ := json.Marshal(cond)
		//fmt.Println(string(a))
		for key, val := range item.output_condtion {
			mapData, ok := cond.(common.KvMap)
			if !ok {
				t.Errorf("test item %v return cond %s not kvmap", item, cond)
				continue
			}
			switch valRaw := val.(type) {
			case string:
				condVal, ok := mapData[key].(string)
				if !ok {
					t.Errorf("test item %v return cond %s key %s val not string", item, cond, key)
					continue
				}
				if valRaw != condVal {
					t.Errorf("test item %v return cond %s key %s val %s not equal %s ", item, cond, key, valRaw, condVal)
					continue
				}
			case []int64:
				condVal, ok := mapData[key].([]int64)
				if !ok {
					t.Errorf("test item %v return cond %s key %s val not []int64", item, cond, key)
					continue
				}
				if len(valRaw) != len(condVal) {
					t.Errorf("test item %v return cond %s key %s val %v not equal %v ", item, cond, key, valRaw, condVal)
					continue
				}
				for _, valRawItem := range valRaw {
					isExist := false
					for _, condValItem := range condVal {
						if valRawItem == condValItem {
							isExist = true
							break
						}
					}
					if !isExist {
						t.Errorf("test item %v return cond %s key %s val %v not equal %v ", item, cond, key, valRaw, condVal)
						break
					}
				}
			case []string:
				condVal, ok := mapData[key].([]string)
				if !ok {
					t.Errorf("test item %v return cond %s key %s val not []string", item, cond, key)
					continue
				}
				if len(valRaw) != len(condVal) {
					t.Errorf("test item %v return cond %s key %s val %v not equal %v ", item, cond, key, valRaw, condVal)
					continue
				}
				for _, valRawItem := range valRaw {
					isExist := false
					for _, condValItem := range condVal {
						if valRawItem == condValItem {
							isExist = true
							break
						}
					}
					if !isExist {
						t.Errorf("test item %v return cond %s key %s val %v not equal %v ", item, cond, key, valRaw, condVal)
						break
					}
				}
			case []interface{}:
				condVal, ok := mapData[key].([]interface{})
				if !ok {
					t.Errorf("test item %v return cond %s key %s val not []interface", item, cond, key)
					continue
				}
				if len(valRaw) != len(condVal) {
					t.Errorf("test item %v return cond %s key %s val %v not equal %v ", item, cond, key, valRaw, condVal)
					continue
				}
				for _, valRawItem := range valRaw {
					isExist := false
					for _, condValItem := range condVal {
						if valRawItem == condValItem {
							isExist = true
							break
						}
					}
					if !isExist {
						t.Errorf("test item %v return cond %s key %s val %v not equal %v ", item, cond, key, valRaw, condVal)
						break
					}
				}
			default:
				t.Errorf("not support data type, return cond %s ", valRaw)
				continue

			}

		}

	}

	//jo?
	//?ppt
	//ap?t
	//app?
	//11?1
	//?11
	//11?
	//11[2,3,10]
	//app[o,t]
	//aaa[1-10]
	//aaa[1,4-3,a]
	//aaa[1,4-3,a]

}

func TestPrivateSplitRegexRange(t *testing.T) {
	type checkData struct {
		reg      string
		isString bool
		result   ScopeMatch
		isErr    bool
	}
	tc := []checkData{
		checkData{
			reg:      "xxx[1-2]",
			isString: true,
			result: ScopeMatch{
				prefix: "xxx",
				ranges: []scopeItem{scopeItem{min: 1, max: 2}},
			},
		},
		checkData{
			reg:      "xxx[x1,1-2,x100,3-100]",
			isString: true,
			result: ScopeMatch{
				prefix: "xxx",
				mixed:  []string{"xxxx1", "xxxx100"},
				ranges: []scopeItem{scopeItem{min: 1, max: 2}},
			},
		},
		checkData{
			reg: "1[1-2]",
			result: ScopeMatch{
				prefix: "1",
				ranges: []scopeItem{scopeItem{min: 1, max: 2}},
			},
		},
		checkData{
			reg: "1[1,2-3,400,5-100]",
			result: ScopeMatch{
				prefix: "1",
				mixed:  []string{"11", "1400"},
				ranges: []scopeItem{scopeItem{min: 2, max: 3}, scopeItem{min: 5, max: 100}},
			},
		},
		checkData{
			reg: "1[1,6-3,400,5-100]",
			result: ScopeMatch{
				prefix: "1",
				mixed:  []string{"11", "1400"},
				ranges: []scopeItem{scopeItem{min: 5, max: 100}},
			},
		},
		checkData{
			reg:    "1[1,a-3,400,5-100]",
			isErr:  true,
			result: ScopeMatch{},
		},
		checkData{
			reg:    "1[1,6-d,400,5-100]",
			isErr:  true,
			result: ScopeMatch{},
		},
		checkData{
			reg:    "1[1,6-d,400,5-100-]",
			isErr:  true,
			result: ScopeMatch{},
		},
		checkData{
			reg:    "1[1,6-d,400,5-]",
			isErr:  true,
			result: ScopeMatch{},
		},
	}

	for _, tc := range tc {
		reg := NewScopeMatch(tc.reg, tc.isString)
		_, err := reg.ParseConds()
		if tc.isErr {
			if nil == err {
				t.Errorf(" item %v need return error", reg)
			}
			continue
		}

		if len(reg.mixed) != len(tc.result.mixed) {
			t.Errorf("test item %v return regex %s len error val %v not equal %v ", tc, tc.reg, reg.mixed, tc.result.mixed)
			continue
		}
		if reg.prefix != tc.result.prefix {
			t.Errorf("test item %v prefix %s not equal %s ", tc, reg.prefix, tc.result.prefix)
			continue
		}
		for _, mixed := range tc.result.mixed {
			isExist := false
			for _, retMixed := range reg.mixed {
				if mixed == retMixed {
					isExist = true
					break
				}
			}
			if !isExist {
				t.Errorf("test item %v return regex %s val %v not equal %v ", tc, tc.reg, reg.mixed, tc.result.mixed)
				continue
			}
		}
		for _, rangeItem := range tc.result.ranges {
			isExist := false
			for _, retRangeItem := range reg.ranges {
				if rangeItem.min == retRangeItem.min && rangeItem.max == retRangeItem.max {
					isExist = true
					break
				}
			}
			if !isExist {
				t.Errorf("test item %v return regex %s val %v not equal %v ", tc, tc.reg, reg.ranges, tc.result.ranges)
				continue
			}

		}
	}
}

func TestMatchStr(t *testing.T) {
	max := 100
	min := 1
	//regStr := fmt.Sprintf("aa[x1,x2, %d-%d]", min, max)
	type testMatchItem struct {
		result bool
		val    string
	}
	type testMatch struct {
		regStr         string
		testMatchItems []testMatchItem
		isRand         bool
		prefix         string
	}

	tm := []testMatch{
		testMatch{
			regStr: "aa[x1,x2, 50-60]",
			testMatchItems: []testMatchItem{
				testMatchItem{result: true, val: "aax1"},
				testMatchItem{result: true, val: "aax2"},
				testMatchItem{result: true, val: "aa50"},
				testMatchItem{result: true, val: "aa60"},
				testMatchItem{result: true, val: "aa55"},
				testMatchItem{result: true, val: "aa59"},
				testMatchItem{result: false, val: "aa"},
				testMatchItem{result: false, val: "aax11"},
				testMatchItem{result: false, val: "aax12"},
				testMatchItem{result: false, val: "aa0"},
				testMatchItem{result: false, val: "aa61"},
			},
		},
		testMatch{
			regStr: fmt.Sprintf("aa[%d-%d]", min, max),
			prefix: "aa",
			isRand: true,
			testMatchItems: []testMatchItem{
				testMatchItem{result: true, val: fmt.Sprintf("aa%d", min)},
				testMatchItem{result: true, val: fmt.Sprintf("aa%d", max)},
			},
		},
		testMatch{
			regStr: "[x1,x2, 50-60]",
			testMatchItems: []testMatchItem{
				testMatchItem{result: true, val: "x1"},
				testMatchItem{result: true, val: "x2"},
				testMatchItem{result: true, val: "50"},
				testMatchItem{result: true, val: "60"},
				testMatchItem{result: true, val: "55"},
				testMatchItem{result: true, val: "59"},
				testMatchItem{result: false, val: ""},
				testMatchItem{result: false, val: "x11"},
				testMatchItem{result: false, val: "x12"},
				testMatchItem{result: false, val: "0"},
				testMatchItem{result: false, val: "61"},
			},
		},
		testMatch{
			regStr: fmt.Sprintf("[%d-%d]", min, max),
			isRand: true,
			testMatchItems: []testMatchItem{
				testMatchItem{result: true, val: fmt.Sprintf("%d", min)},
				testMatchItem{result: true, val: fmt.Sprintf("%d", max)},
			},
		},
	}

	for _, tmItem := range tm {
		regStr := tmItem.regStr
		reg := NewScopeMatch(regStr, true)

		_, err := reg.ParseConds()
		if nil != err {
			t.Errorf("test regex %s new regex role errror %s", regStr, err.Error())
			return
		}
		if tmItem.isRand {
			for idx := 0; idx < int((max-min)/2); idx++ {
				num := rand.Intn(max + max)
				if num > max || num < min {
					tmItem.testMatchItems = append(tmItem.testMatchItems, testMatchItem{result: false, val: fmt.Sprintf("%s%d", tmItem.prefix, num)})
				} else {
					tmItem.testMatchItems = append(tmItem.testMatchItems, testMatchItem{result: true, val: fmt.Sprintf("%s%d", tmItem.prefix, num)})
				}
			}
		}
		for _, item := range tmItem.testMatchItems {
			result := reg.MatchStr(item.val)
			if result != item.result {
				t.Errorf("test regex %s match %s must be %v not %v", regStr, item.val, item.result, result)
				continue
			}
		}
	}

}

func TestMatchInt(t *testing.T) {
	max := int64(100)
	min := int64(20)
	//baseMin := int64(math.Pow10(len(string(min)) - 1))
	//baseMax := int64(math.Pow10(len(string(max)) - 1))
	//regStr := fmt.Sprintf("aa[x1,x2, %d-%d]", min, max)
	type testMatchItem struct {
		result bool
		val    int64
	}
	type testMatch struct {
		regStr         string
		isRand         bool
		testMatchItems []testMatchItem
	}

	tm := []testMatch{
		testMatch{
			regStr: "1[1,2, 50-60]",
			testMatchItems: []testMatchItem{
				testMatchItem{result: true, val: 11},
				testMatchItem{result: true, val: 12},
				testMatchItem{result: true, val: 150},
				testMatchItem{result: true, val: 160},
				testMatchItem{result: true, val: 157},
				testMatchItem{result: false, val: 1},
				testMatchItem{result: false, val: 10},
				testMatchItem{result: false, val: 13},
				testMatchItem{result: false, val: 149},
				testMatchItem{result: false, val: 161},
				testMatchItem{result: false, val: 1151},
			},
		},
		testMatch{
			regStr: "11[1-100]",
			testMatchItems: []testMatchItem{
				testMatchItem{result: true, val: 111},
				testMatchItem{result: true, val: 11100},
				testMatchItem{result: true, val: 1102},
				testMatchItem{result: true, val: 1101},
				testMatchItem{result: true, val: 11001},
				testMatchItem{result: true, val: 1199},
				testMatchItem{result: false, val: 11},
				testMatchItem{result: false, val: 110},
			},
		},
		testMatch{
			regStr: fmt.Sprintf("[%d-%d]", min, max),
			isRand: true,
			testMatchItems: []testMatchItem{
				testMatchItem{result: true, val: min},
				testMatchItem{result: true, val: max},
			},
		},
	}

	for _, tmItem := range tm {
		regStr := tmItem.regStr
		reg := NewScopeMatch(regStr, true)

		_, err := reg.ParseConds()
		if nil != err {
			t.Errorf("test regex %s new regex role errror %s", regStr, err.Error())
			return
		}
		if tmItem.isRand {
			for idx := 0; idx < int((max-min)/2); idx++ {
				num := rand.Int63n(max+(max+1)/2) + (min+1)/2
				if num > max || num < min {
					tmItem.testMatchItems = append(tmItem.testMatchItems, testMatchItem{result: false, val: num})
				} else {
					tmItem.testMatchItems = append(tmItem.testMatchItems, testMatchItem{result: true, val: num})
				}
			}
		}
		for _, item := range tmItem.testMatchItems {
			result := reg.MatchInt64(item.val)
			if result != item.result {
				t.Errorf("test regex %s match %d must be %v not %v", regStr, item.val, item.result, result)
				continue
			}
		}
	}

}
