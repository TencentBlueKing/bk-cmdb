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
	"fmt"
	"net"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

const (
	// INNERONLY TODO
	INNERONLY string = "bk_host_innerip"
	// OUTERONLY TODO
	OUTERONLY string = "bk_host_outerip"
	// IOBOTH TODO
	IOBOTH string = "bk_host_innerip|bk_host_outerip"
)

// ParseHostParams parses the host condition into query statement
func ParseHostParams(input []metadata.ConditionItem) (map[string]interface{}, error) {
	var err error
	output := make(map[string]interface{})
	for _, i := range input {
		switch i.Operator {
		case common.BKDBEQ:
			output[i.Field], err = common.ConvertIpv6ToFullWord(i.Field, i.Value)
			if err != nil {
				return nil, err
			}
		case common.BKDBNE:
			value, err := common.ConvertIpv6ToFullWord(i.Field, i.Value)
			if err != nil {
				return nil, err
			}
			output[i.Field] = mapstr.MapStr{i.Operator: value}
		case common.BKDBIN, common.BKDBNIN:
			queryCondItem := make(map[string]interface{})
			if i.Value == nil {
				queryCondItem[i.Operator] = make([]interface{}, 0)
			} else {
				queryCondItem[i.Operator], err = common.ConvertIpv6ToFullWord(i.Field, i.Value)
				if err != nil {
					return nil, err
				}
			}
			output[i.Field] = queryCondItem
		case common.BKDBLIKE:
			regex := make(map[string]interface{})
			regex[common.BKDBLIKE] = i.Value
			output[i.Field] = regex
		case common.BKDBMULTIPLELike:
			multi, ok := i.Value.([]interface{})
			if !ok {
				return output, fmt.Errorf("operator %s only support for string array", common.BKDBMULTIPLELike)
			}
			fields := make([]interface{}, 0)
			for _, m := range multi {
				mstr, ok := m.(string)
				if !ok {
					return output, fmt.Errorf("operator %s only support for string array", common.BKDBMULTIPLELike)
				}
				fields = append(fields, mapstr.MapStr{i.Field: mapstr.MapStr{common.BKDBLIKE: mstr}})
			}
			if len(fields) != 0 {
				// only when the fields is none empty, then the fields is valid.
				// a or operator can not have a empty value in mongodb.
				output[common.BKDBOR] = fields
			}
		default:
			queryCondItem, ok := output[i.Field].(map[string]interface{})
			if !ok {
				queryCondItem = make(map[string]interface{})
			}
			switch rawVal := i.Value.(type) {
			case string:
				if timeType, isTime := util.IsTime(rawVal); isTime {
					i.Value = util.Str2Time(rawVal, timeType)
				}
			}
			queryCondItem[i.Operator] = i.Value
			output[i.Field] = queryCondItem
		}
	}
	return output, nil
}

// ParseHostIPParams parses the IP address information into query statement
func ParseHostIPParams(ipv4Cond metadata.IPInfo, ipv6Cond metadata.IPInfo,
	output map[string]interface{}) (map[string]interface{}, error) {
	var err error
	exactOr := make([]map[string]interface{}, 0)
	embeddedIPv4Addrs := make([]string, 0)
	if len(ipv6Cond.Data) != 0 {
		exactOr, embeddedIPv4Addrs, err = parseIPv6Condition(ipv6Cond, exactOr)
		if err != nil {
			return nil, fmt.Errorf("failed to add ipv6 addresses to condition, err: %v", err)
		}
	}

	ipArr := append(ipv4Cond.Data, embeddedIPv4Addrs...)
	exact := ipv4Cond.Exact
	flag := ipv4Cond.Flag
	if len(ipArr) == 0 && len(exactOr) == 0 {
		return output, nil
	}

	if exact == 1 {
		// exact search
		// filter out illegal IPv4 addresses
		ipArr = filterHostIP(ipArr)
		exactIP := map[string]interface{}{common.BKDBIN: deduplication(ipArr)}
		exactOr, err = addExactSearchCondition(exactOr, exactIP, flag, "ipv4")
		if err != nil {
			return nil, err
		}
		output[common.BKDBOR] = exactOr
	} else {
		// not exact search
		orCond := make([]map[string]map[string]interface{}, 0)
		orCond, err = addInexactSearchCondition(orCond, ipArr, flag)
		if err != nil {
			return nil, err
		}
		output[common.BKDBOR] = orCond
	}
	return output, nil
}

// parseIPv6Condition parse IPv6 conditions to full Ipv6 addresses and embedded IPv4 addresses
// only full or abbreviated IPv6 addresses can be used for exact queries, not exact search is not supported
func parseIPv6Condition(ipCond metadata.IPInfo, exactOr []map[string]interface{}) ([]map[string]interface{},
	[]string, error) {
	ipArr := ipCond.Data
	legalIpArr := filterHostIP(ipArr)
	flag := ipCond.Flag
	if len(legalIpArr) == 0 || ipCond.Exact != 1 {
		return exactOr, nil, nil
	}

	fullIpv6Addrs := make([]string, 0)
	embeddedIPv4Addrs := make([]string, 0)
	for _, address := range legalIpArr {
		ip, err := common.GetIPv4IfEmbeddedInIPv6(address)
		if err != nil {
			continue
		}
		// 对于兼容IPv4的嵌入式IPv6地址，::127.0.0.1和::ffff:127.0.0.1这两种格式的地址，存放于ipv4字段中，所以使用ipv4的字段查询
		if !strings.Contains(ip, ":") {
			embeddedIPv4Addrs = append(embeddedIPv4Addrs, ip)
			continue
		}
		fullIpv6Addr, err := common.ConvertIPv6ToStandardFormat(ip)
		if err != nil {
			continue
		}
		fullIpv6Addrs = append(fullIpv6Addrs, fullIpv6Addr)
	}

	if len(fullIpv6Addrs) == 0 {
		return exactOr, embeddedIPv4Addrs, nil
	}

	// exact search ipv6 addr
	var err error
	exactIP := map[string]interface{}{common.BKDBIN: deduplication(fullIpv6Addrs)}
	exactOr, err = addExactSearchCondition(exactOr, exactIP, flag, "ipv6")
	if err != nil {
		return nil, nil, err
	}
	return exactOr, embeddedIPv4Addrs, nil
}

// filterHostIP filter out illegal IP addresses
func filterHostIP(ipArr []string) []string {
	legalAddress := make([]string, 0)
	for _, address := range ipArr {
		if ip := net.ParseIP(address); ip == nil {
			continue
		}
		legalAddress = append(legalAddress, address)
	}
	return legalAddress
}

// addExactSearchCondition combine query statements based on exact ip conditions
func addExactSearchCondition(exactOr []map[string]interface{}, exactIP map[string]interface{}, flag string,
	ipType string) ([]map[string]interface{}, error) {
	switch ipType {
	case "ipv4":
		switch flag {
		case INNERONLY:
			exactOr = append(exactOr, mapstr.MapStr{common.BKHostInnerIPField: exactIP})
		case OUTERONLY:
			exactOr = append(exactOr, mapstr.MapStr{common.BKHostOuterIPField: exactIP})
		case IOBOTH:
			exactOr = append(exactOr, mapstr.MapStr{common.BKHostInnerIPField: exactIP},
				mapstr.MapStr{common.BKHostOuterIPField: exactIP})
		default:
			return exactOr, fmt.Errorf("unsupported ip.flag %s", flag)
		}
	case "ipv6":
		switch flag {
		case INNERONLY:
			exactOr = append(exactOr, mapstr.MapStr{common.BKHostInnerIPv6Field: exactIP})
		case OUTERONLY:
			exactOr = append(exactOr, mapstr.MapStr{common.BKHostOuterIPv6Field: exactIP})
		case IOBOTH:
			exactOr = append(exactOr, mapstr.MapStr{common.BKHostInnerIPv6Field: exactIP},
				mapstr.MapStr{common.BKHostOuterIPv6Field: exactIP})
		default:
			return exactOr, fmt.Errorf("unsupported ip.flag %s", flag)
		}
	default:
		return exactOr, fmt.Errorf("unsupported ip type %s", ipType)
	}
	return exactOr, nil
}

// addInexactSearchCondition combine query statements based on inexact ip conditions
func addInexactSearchCondition(orCond []map[string]map[string]interface{},
	ipArr []string, flag string) ([]map[string]map[string]interface{}, error) {
	for _, ip := range ipArr {
		c := make(map[string]interface{})
		c[common.BKDBLIKE] = SpecialCharChange(ip)
		switch flag {
		case INNERONLY:
			orCond = append(orCond, map[string]map[string]interface{}{
				common.BKHostInnerIPField: c,
			})
		case OUTERONLY:
			orCond = append(orCond, map[string]map[string]interface{}{
				common.BKHostOuterIPField: c,
			})
		case IOBOTH:
			orCond = append(orCond, []map[string]map[string]interface{}{{common.BKHostOuterIPField: c},
				{common.BKHostInnerIPField: c}}...)
		default:
			return orCond, fmt.Errorf("unsupported ip.flag %s", flag)
		}
	}
	return orCond, nil
}

// deduplication remove duplicate IP addresses from the IP array
func deduplication(arr []string) []string {
	result := make([]string, 0)
	m := make(map[string]bool) //map的值不重要
	for _, v := range arr {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}
