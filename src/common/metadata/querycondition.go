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

package metadata

import (
	"time"

	"github.com/coccyx/timeparser"

	"configcenter/src/common"
	"configcenter/src/common/util"
)

const (
	CC_time_type_parse_flag = "cc_time_type"
)

type ObjQueryInput struct {
	Condition interface{} `json:"condition"`
	Fields    string      `json:"fields"`
	Start     int         `json:"start"`
	Limit     int         `json:"limit"`
	Sort      string      `json:"sort"`
}

//ConvTime 将查询条件中字段包含cc_type key ，子节点变为time.Time
func (o *ObjQueryInput) ConvTime() error {
	conds, ok := o.Condition.(map[string]interface{})
	if true != ok && nil != conds {
		return nil
	}
	for key, item := range conds {
		convItem, err := o.convTimeItem(item)
		if nil != err {
			continue
		}
		conds[key] = convItem
	}

	return nil
}

//convTimeItem 转义具体的某一项,将查询条件中字段包含cc_time_type
func (o *ObjQueryInput) convTimeItem(item interface{}) (interface{}, error) {

	switch item.(type) {
	case map[string]interface{}:

		arrItem, ok := item.(map[string]interface{})
		if true == ok {
			_, timeTypeOk := arrItem[common.BKTimeTypeParseFlag]
			if timeTypeOk {
				delete(arrItem, common.BKTimeTypeParseFlag)
			}

			for key, value := range arrItem {
				switch value.(type) {

				case []interface{}:
					var err error
					arrItem[key], err = o.convTimeItem(value)
					if nil != err {
						return nil, err
					}
				case map[string]interface{}:
					arrItemVal, ok := value.(map[string]interface{})
					if ok {
						for key, value := range arrItemVal {
							var err error
							arrItemVal[key], err = o.convTimeItem(value)
							if nil != err {
								return nil, err
							}
						}
						arrItem[key] = value
					}

				default:
					if timeTypeOk {
						var err error
						arrItem[key], err = o.convInterfaceToTime(value)
						if nil != err {
							return nil, err
						}
					}

				}
			}
			item = arrItem
		}
	case []interface{}:
		//如果是数据，递归转换所有子项
		arrItem, ok := item.([]interface{})
		if true == ok {
			for index, value := range arrItem {
				newValue, err := o.convTimeItem(value)
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

func (O *ObjQueryInput) convInterfaceToTime(val interface{}) (interface{}, error) {
	switch val.(type) {
	case string:
		ts, err := timeparser.TimeParser(val.(string))
		if nil != err {
			return nil, err
		}
		return ts.UTC(), nil
	default:
		ts, err := util.GetInt64ByInterface(val)
		if nil != err {
			return 0, err
		}
		t := time.Unix(ts, 0).UTC()
		return t, nil
	}

}
