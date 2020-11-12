/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"configcenter/src/common"

	"github.com/coccyx/timeparser"
)

var (
	// 需要转换的时间的标志
	convTimeFields = []string{common.CreateTimeField, common.LastTimeField, common.ConfirmTimeField}
)

func GetCurrentTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetCurrentTimePtr() *time.Time {
	now := time.Now()
	return &now
}

func ConvParamsTime(data interface{}) interface{} {
	conds, ok := data.(map[string]interface{})
	if true != ok && nil != conds {
		return data
	}
	convItem, _ := convTimeItem(data)

	return convItem
}

func convTimeItem(item interface{}) (interface{}, error) {

	switch item.(type) {
	case map[string]interface{}:
		arrItem, ok := item.(map[string]interface{})
		if true == ok {

			for key, value := range arrItem {
				var timeTypeOk bool = false
				for _, convTimeKey := range convTimeFields {
					if key == convTimeKey {
						timeTypeOk = true
						break
					}
				}
				// 如果当前不需要转换，递归转
				if !timeTypeOk {
					arrItem[key], _ = convTimeItem(value)
					continue
				}

				switch value.(type) {
				case []interface{}:
					arr := value.([]interface{})
					for index, tsValue := range arr {
						ts, err := convInterfaceToTime(tsValue)
						if nil != err {
							continue
						}
						arr[index] = ts
					}
					arrItem[key] = arr
				case map[string]interface{}:
					arr := value.(map[string]interface{})
					for mapKey, mapVal := range arr {
						ts, err := convInterfaceToTime(mapVal)
						if nil != err {
							continue
						}
						arr[mapKey] = ts
					}
					arrItem[key] = arr
				case string:
					ts, err := convInterfaceToTime(value)
					if nil == err {
						arrItem[key] = ts
					}

				}
			}
			item = arrItem
		}
	case []interface{}:
		// 如果是数据，递归转换所有子项
		arrItem, ok := item.([]interface{})
		if true == ok {
			for index, value := range arrItem {
				newValue, err := convTimeItem(value)
				if nil != err {
					return nil, err

				}
				arrItem[index] = newValue
			}
			item = arrItem

		}
	}

	return item, nil
}

func convInterfaceToTime(val interface{}) (interface{}, error) {
	switch val.(type) {
	case []interface{}:
		arrVal, _ := val.([]interface{})
		var ret []interface{}
		for _, itemVal := range arrVal {
			ts, err := convItemToTime(itemVal)
			if nil != err {
				ret = append(ret, itemVal)
			} else {
				ret = append(ret, ts)

			}
		}
		return ret, nil
	default:
		return convItemToTime(val)
	}

}

func convItemToTime(val interface{}) (interface{}, error) {
	switch val.(type) {
	case string:
		ts, err := timeparser.TimeParser(val.(string))
		if nil != err {
			return nil, err
		}
		return ts.UTC(), nil

	default:
		ts, err := GetInt64ByInterface(val)
		if nil != err {
			return 0, err
		}
		t := time.Unix(ts, 0).UTC()
		return t, nil
	}

}

var validPeriod = regexp.MustCompile("^\\d*[DHMS]$") // period regexp to check period

// 00002H --> 2H
// 0000D/0M ---> ∞
// empty string / ∞ ---> ∞
// regexp matched: positive integer (include positive integer begin with more the one '0') + [D/H/M/S]
// eg. 0H, 000H, 0002H, 32M，34S...
// examples of no matched:  1.4H, -2H, +2H ...
func FormatPeriod(period string) (string, error) {
	if common.Infinite == period || "" == period {
		return common.Infinite, nil
	}

	if !validPeriod.MatchString(period) {
		return "", fmt.Errorf("invalid period")
	}

	num, err := strconv.Atoi(period[:len(period)-1])
	if nil != err {
		return "", err
	}
	if 0 == num {
		return common.Infinite, nil
	}

	return strconv.Itoa(num) + period[len(period)-1:], nil
}
