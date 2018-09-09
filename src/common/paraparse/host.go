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
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
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

func ParseHostParams(input []metadata.ConditionItem, output map[string]interface{}) error {
	for _, i := range input {
		switch i.Operator {
		case common.BKDBEQ:
			output[i.Field] = i.Value
		case common.BKDBIN:
			queryCondItem := make(map[string]interface{})
			queryCondItem[i.Operator] = i.Value
			output[i.Field] = queryCondItem
		case common.BKDBLIKE:
			//d := make(map[string]interface{})
			queryCondItem, ok := output[i.Field].(map[string]interface{})
			if !ok {
				queryCondItem = make(map[string]interface{})
			}
			valStr, ok := i.Value.(string)
			if ok {
				queryCondItem[i.Operator] = SpeceialCharChange(valStr)
			} else {
				queryCondItem[i.Operator] = i.Value
			}
			output[i.Field] = queryCondItem
		default:
			queryCondItem, ok := output[i.Field].(map[string]interface{})
			if !ok {
				queryCondItem = make(map[string]interface{})
			}
			switch rawVal := i.Value.(type) {
			case string:
				if util.IsTime(rawVal) {
					i.Value = util.Str2Time(rawVal)
				}
			}
			queryCondItem[i.Operator] = i.Value
			output[i.Field] = queryCondItem
		}
	}
	return nil
}

func ParseHostIPParams(ipCond metadata.IPInfo, output map[string]interface{}) error {
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
			c[common.BKDBLIKE] = SpeceialCharChange(ip)
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
