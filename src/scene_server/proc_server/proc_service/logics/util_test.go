package logics

import (
	"testing"

	"configcenter/src/common"
)

type testData struct {
	input           string
	isInteger       bool
	output_condtion common.KvMap
	output_notParse bool
	isErr           bool
	isNil           bool
}

func TestCondtion(t *testing.T) {
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
			output_notParse: true,
		},
		testData{
			input:           "aaa[1,4-3,a]",
			output_condtion: nil, //"^app.$",
			output_notParse: true,
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
		cond, notParse, err := ParseProcInstMatchCondition(item.input, !item.isInteger)

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
		if notParse != item.output_notParse {
			t.Errorf("test not parse return error, need %v not %v", item.output_notParse, notParse)
			continue
		}
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
					t.Errorf("test item %v return cond %s key % val not string", item, cond, key)
					continue
				}
				if valRaw != condVal {
					t.Errorf("test item %v return cond %s key % val %s not equal %s ", item, cond, key, valRaw, condVal)
					continue
				}
			case []int64:
				condVal, ok := mapData[key].([]int64)
				if !ok {
					t.Errorf("test item %v return cond %s key % val not []int64", item, cond, key)
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
						t.Errorf("test item %v return cond %s key % val %d not equal %d ", item, cond, key, valRaw, condVal)
						continue
					}
				}
			case []string:
				condVal, ok := mapData[key].([]string)
				if !ok {
					t.Errorf("test item %v return cond %s key % val not []string", item, cond, key)
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
						t.Errorf("test item %v return cond %s key % val %v not equal %v ", item, cond, key, valRaw, condVal)
						continue
					}
				}
			case []interface{}:
				condVal, ok := mapData[key].([]interface{})
				if !ok {
					t.Errorf("test item %v return cond %s key % val not []interface", item, cond, key)
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
						t.Errorf("test item %v return cond %s key % val %v not equal %v ", item, cond, key, valRaw, condVal)
						continue
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
