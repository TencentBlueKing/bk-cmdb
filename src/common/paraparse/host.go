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

package params

import (
	"configcenter/src/common"
	"configcenter/src/common/util"
	"errors"
	"fmt"
)

//type Flag string

const (
	INNERONLY string = "bk_host_innerip"
	OUTERONLY string = "bk_host_outerip"
	IOBOTH    string = "bk_host_innerip|bk_host_outerip"
)

//common search struct
type HostCommonSearch struct {
	AppID     int               `json:"bk_biz_id,omitempty"`
	Ip        IPInfo            `json:"ip"`
	Condition []SearchCondition `json:"condition"`
	Page      PageInfo          `json:"page"`
	Pattern   string            `json:"pattern,omitempty"`
}

//ip search info
type IPInfo struct {
	Data  []string `json:"data"`
	Exact int      `json:"exact"`
	Flag  string   `json:"flag"`
}

//common page info
type PageInfo struct {
	Start int    `json:"start"`
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
}

//search condition
type SearchCondition struct {
	Fields    []string      `json:"fields"`
	Condition []interface{} `json:"condition"`
	ObjectID  string        `json:"bk_obj_id"`
}

//NewHostCommonSearch new host common search
type NewHostCommonSearch struct {
	IP        IPInfo                            `json:"ip"`
	Condition map[string]map[string]interface{} `json:"condition"`
	Page      PageInfo                          `json:"page"`
	Pattern   string                            `json:"pattern,omitempty"`
	Field     map[string][]string               `json:"field"`
}

func ParseCondParams(input map[string]interface{}) ([]interface{}, error) {
	output := make([]interface{}, 0)
	for k, v := range input {
		condition := make(map[string]interface{})
		condition["field"] = k
		condition["operator"] = common.BKDBEQ
		condition["value"] = v
		output = append(output, condition)
	}
	return output, nil
}

func ParseHostParams(input []interface{}, output map[string]interface{}) error {
	fmt.Println(input)
	for _, i := range input {
		j, ok := i.(map[string]interface{})
		if false == ok {
			return errors.New("condition error")
		}
		field, ok := j["field"].(string)
		if false == ok {
			return errors.New("condition error")
		}
		operator, ok := j["operator"].(string)
		if false == ok {
			return errors.New("condition error")
		}
		value := j["value"]

		switch operator {
		case common.BKDBEQ:
			output[field] = value
		case common.BKDBIN:
			d := make(map[string]interface{})
			d[operator] = value
			output[field] = d
		default:
			d := make(map[string]interface{})
			switch value.(type) {
			case string:

				valStr := value.(string)
				if util.IsTime(valStr) {

					value = util.Str2Time(valStr)
				}

			}
			d[operator] = value
			output[field] = d
		}

	}
	fmt.Println(output)
	return nil
}

func ParseHostIPParams(ipCond IPInfo, output map[string]interface{}) error {
	ipArr := ipCond.Data
	exact := ipCond.Exact
	flag := ipCond.Flag
	if 0 == len(ipArr) {
		return nil
	}
	if 1 == exact {
		//exact search
		c := make(map[string]interface{})
		c[common.BKDBIN] = ipArr
		if INNERONLY == flag {
			output[common.BKHostInnerIPField] = c

		} else if OUTERONLY == flag {
			output[common.BKHostOuterIPField] = c
		} else if IOBOTH == flag {
			io := make([]map[string]interface{}, 2)
			i := make(map[string]interface{})
			o := make(map[string]interface{})
			ic := make(map[string]interface{})
			oc := make(map[string]interface{})
			i[common.BKDBIN] = ipArr
			o[common.BKDBIN] = ipArr
			ic[common.BKHostInnerIPField] = i
			oc[common.BKHostOuterIPField] = o
			io[0] = ic
			io[1] = oc
			output[common.BKDBOR] = io
		}
	} else {
		//not exact search
		orCond := make([]map[string]map[string]interface{}, 0)
		for _, ip := range ipArr {
			c := make(map[string]interface{})
			c[common.BKDBLIKE] = ip
			if INNERONLY == flag {
				ipCon := make(map[string]map[string]interface{})
				ipCon[common.BKHostInnerIPField] = c
				orCond = append(orCond, ipCon)
			} else if OUTERONLY == flag {
				ipCon := make(map[string]map[string]interface{})
				ipCon[common.BKHostOuterIPField] = c
				orCond = append(orCond, ipCon)
			} else if IOBOTH == flag {
				ipiCon := make(map[string]map[string]interface{})
				ipoCon := make(map[string]map[string]interface{})
				ipoCon[common.BKHostOuterIPField] = c
				ipiCon[common.BKHostInnerIPField] = c
				orCond = append(orCond, ipoCon)
				orCond = append(orCond, ipiCon)
			}
			output[common.BKDBOR] = orCond
		}

	}
	return nil
}
