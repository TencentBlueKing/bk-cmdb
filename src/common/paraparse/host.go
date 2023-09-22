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
	"encoding/json"
	"fmt"
	"net"
	"strconv"
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
		case common.BKDBEQ, common.BKDBNE:
			output[i.Field], err = common.ConvertIpv6ToFullWord(i.Field, i.Value)
			if err != nil {
				return nil, err
			}
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
func ParseHostIPParams(ipv4Cond metadata.IPInfo, ipv6Cond metadata.IPInfo, output map[string]interface{}) (
	map[string]interface{}, error) {

	var err error
	exactOr := make([]map[string]interface{}, 0)
	embeddedIPv4Addrs := make([]string, 0)

	if len(ipv6Cond.Data) != 0 {
		exactOr, embeddedIPv4Addrs, err = parseIPv6Condition(ipv6Cond, exactOr, &output)
		if err != nil {
			return nil, fmt.Errorf("failed to add ipv6 addresses to condition, err: %v", err)
		}
	}

	ipv4Cond.Data = append(ipv4Cond.Data, embeddedIPv4Addrs...)
	exact := ipv4Cond.Exact
	flag := ipv4Cond.Flag
	if len(ipv4Cond.Data) == 0 && len(exactOr) == 0 {
		return output, nil
	}

	ipv4CloudIDMap, err := splitIPv4Data(ipv4Cond, &output)
	if err != nil {
		return nil, err
	}

	if exact == 1 {
		// exact search
		// filter out illegal IPv4 addresses
		exactOr, err = addIPv4ExactSearchCondition(exactOr, ipv4CloudIDMap, flag)
		if err != nil {
			return nil, err
		}
		output[common.BKDBOR] = exactOr
	} else {
		// not exact search
		orCond := make([]map[string]map[string]interface{}, 0)
		orCond, err = addFuzzyCondition(orCond, ipv4CloudIDMap, flag)
		if err != nil {
			return nil, err
		}
		output[common.BKDBOR] = orCond
	}
	return output, nil
}

// parseIPv6Condition parse IPv6 conditions to full Ipv6 addresses and embedded IPv4 addresses
// only full or abbreviated IPv6 addresses can be used for exact queries, not exact search is not supported
func parseIPv6Condition(ipCond metadata.IPInfo, exactOr []map[string]interface{}, output *map[string]interface{}) (
	[]map[string]interface{}, []string, error) {

	flag := ipCond.Flag
	if ipCond.Exact != 1 {
		return exactOr, nil, nil
	}

	ipv6CloudIDMap, embeddedIPv4Addrs, err := splitIPv6Data(ipCond, output)
	if err != nil {
		return nil, nil, err
	}
	exactOr, err = addIPv6ExactSearchCondition(exactOr, ipv6CloudIDMap, flag)
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

// splitIPv4Data split ipv4 data
func splitIPv4Data(ipCond metadata.IPInfo, output *map[string]interface{}) (map[int64]map[string][]string, error) {
	// 创建一个 map 用于存储结果
	cloudIDMap := make(map[int64]map[string][]string)

	outputCloudIDMap, err := getCloudIDMapByOutput(*output)
	if err != nil {
		return nil, err
	}

	// 遍历切片中的每个字符串
	for _, ipString := range ipCond.Data {
		// 去除空格
		ipString = strings.TrimSpace(ipString)

		// 分割管控区域 ID 和 IP 地址
		colonIndex := strings.Index(ipString, ":")
		ipPart := ipString[colonIndex+1:]
		// 去掉 IP 地址中的方括号，只保留 IP 部分
		ipAddress := strings.Trim(ipPart, "[]")

		if colonIndex == -1 {
			// 如果没有冒号代表未指定管控区域，判断hostcond里面是否有指定
			if outputCloudIDMap != nil {
				for operator, cloudIDArr := range outputCloudIDMap {
					for _, cloudID := range cloudIDArr {
						// 初始化内部的 map
						if _, ok := cloudIDMap[cloudID]; !ok {
							cloudIDMap[cloudID] = make(map[string][]string)
						}
						cloudIDMap[cloudID][operator] = append(cloudIDMap[cloudID][operator], ipString)
						delete(*output, common.BKCloudIDField)
					}
				}
				continue
			}

			// 如果没有冒号，且 output 条件未指定管控区域则直接添加到 "-1" 键下
			// 初始化内部的 map
			if _, ok := cloudIDMap[-1]; !ok {
				cloudIDMap[-1] = make(map[string][]string)
			}
			cloudIDMap[-1][common.BKDBIN] = append(cloudIDMap[-1][common.BKDBIN], ipString)
			continue
		}

		// 转换管控区域 ID 为 int64 类型
		cloudIDStr := ipString[:colonIndex]
		cloudID, err := strconv.ParseInt(cloudIDStr, 10, 64)
		if err != nil {
			continue
		}
		// 初始化内部的 map
		if _, ok := cloudIDMap[cloudID]; !ok {
			cloudIDMap[cloudID] = make(map[string][]string)
		}
		cloudIDMap[cloudID][common.BKDBIN] = append(cloudIDMap[cloudID][common.BKDBIN], ipAddress)
	}
	return cloudIDMap, nil
}

// splitIPv6Data split ipv6 data
func splitIPv6Data(ipCond metadata.IPInfo, output *map[string]interface{}) (map[int64]map[string][]string, []string, error) {
	// 创建一个 map 用于存储结果
	cloudIDMap := make(map[int64]map[string][]string)
	embeddedIPv4Addrs := make([]string, 0)

	outputCloudIDMap, err := getCloudIDMapByOutput(*output)
	if err != nil {
		return nil, nil, err
	}

	// 遍历切片中的每个字符串
	for _, ipString := range ipCond.Data {
		// 去除空格
		ipString = strings.TrimSpace(ipString)

		// 分割管控区域 ID 和 IP 地址
		colonIndex := strings.Index(ipString, ":[")
		ipPart := ipString[colonIndex+1:]
		// 去掉 IP 地址中的方括号，只保留 IP 部分
		ipAddress := strings.Trim(ipPart, "[]")

		ipAddr, err := common.GetIPv4IfEmbeddedInIPv6(ipAddress)
		if err != nil {
			continue
		}
		// 对于兼容IPv4的嵌入式IPv6地址，::127.0.0.1和::ffff:127.0.0.1这两种格式的地址，存放于ipv4字段中，所以使用ipv4的字段查询
		if !strings.Contains(ipAddr, ":") {
			embeddedIPv4Addrs = append(embeddedIPv4Addrs, ipString)
			continue
		}

		fullIpv6Addr, err := common.ConvertIPv6ToStandardFormat(ipAddr)
		if err != nil {
			continue
		}

		if colonIndex == -1 {

			// 如果没有冒号代表未指定管控区域，判断hostcond里面是否有指定
			if outputCloudIDMap != nil {
				for operator, cloudIDArr := range outputCloudIDMap {
					for _, cloudID := range cloudIDArr {
						// 初始化内部的 map
						if _, ok := cloudIDMap[cloudID]; !ok {
							cloudIDMap[cloudID] = make(map[string][]string)
						}
						cloudIDMap[cloudID][operator] = append(cloudIDMap[cloudID][operator], ipString)
						delete(*output, common.BKCloudIDField)
					}
				}
				continue
			}

			// 如果没有冒号，则直接添加到 "-1" 键下
			// 初始化内部的 map
			if _, ok := cloudIDMap[-1]; !ok {
				cloudIDMap[-1] = make(map[string][]string)
			}
			cloudIDMap[-1][common.BKDBIN] = append(cloudIDMap[-1][common.BKDBIN], ipString)
			continue
		}

		// 转换管控区域 ID 为 int64 类型
		cloudIDStr := ipString[:colonIndex]
		cloudID, err := strconv.ParseInt(cloudIDStr, 10, 64)
		if err != nil {
			continue
		}

		// 初始化内部的 map
		if _, ok := cloudIDMap[cloudID]; !ok {
			cloudIDMap[cloudID] = make(map[string][]string)
		}
		cloudIDMap[cloudID][common.BKDBIN] = append(cloudIDMap[cloudID][common.BKDBIN], fullIpv6Addr)
	}
	return cloudIDMap, embeddedIPv4Addrs, nil
}

func getCloudIDMapByOutput(output map[string]interface{}) (map[string][]int64, error) {
	outputCloudIDMap := make(map[string][]int64)

	cloudIDCond, ok := output[common.BKCloudIDField].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	hasOperator := false
	if cloudIDArray, ok := cloudIDCond[common.BKDBIN].([]interface{}); ok {
		for _, cloudID := range cloudIDArray {
			cloudIDStr, ok := cloudID.(json.Number)
			if !ok {
				return nil, fmt.Errorf("conversion to type 'json.Number' failed")
			}

			cloudIDInt64, err := cloudIDStr.Int64()
			if err != nil {
				return nil, fmt.Errorf("conversion to type int64 failed")
			}

			outputCloudIDMap[common.BKDBIN] = append(outputCloudIDMap[common.BKDBIN], cloudIDInt64)
			hasOperator = true
		}
	}

	if cloudIDArray, ok := cloudIDCond[common.BKDBNIN].([]interface{}); ok {
		for _, cloudID := range cloudIDArray {
			cloudIDStr, ok := cloudID.(json.Number)
			if !ok {
				return nil, fmt.Errorf("conversion to type 'json.Number' failed")
			}

			cloudIDInt64, err := cloudIDStr.Int64()
			if err != nil {
				return nil, fmt.Errorf("conversion to type int64 failed")
			}

			outputCloudIDMap[common.BKDBIN] = append(outputCloudIDMap[common.BKDBIN], cloudIDInt64)
			hasOperator = true
		}
	}

	if hasOperator == false {
		return nil, fmt.Errorf("cloudID query condition is not equal to '$in' or '$nin' failed")
	}
	return outputCloudIDMap, nil
}

// addExactSearchCondition combine query statements based on exact ipv4 conditions
func addIPv4ExactSearchCondition(exactOr []map[string]interface{}, ipv4CloudIDMap map[int64]map[string][]string,
	flag string) ([]map[string]interface{}, error) {

	for cloudID, ipv4OperatorMap := range ipv4CloudIDMap {
		for operator, ipv4Arr := range ipv4OperatorMap {
			// 校验、去重
			ipv4Arr = filterHostIP(ipv4Arr)
			exactIP := map[string]interface{}{operator: deduplication(ipv4Arr)}
			switch flag {
			case INNERONLY:
				if cloudID == -1 {
					exactOr = append(exactOr, mapstr.MapStr{
						common.BKHostInnerIPField: exactIP,
					})
					continue
				}
				exactOr = append(exactOr, mapstr.MapStr{
					common.BKCloudIDField:     cloudID,
					common.BKHostInnerIPField: exactIP,
				})

			case OUTERONLY:
				if cloudID == -1 {
					exactOr = append(exactOr, mapstr.MapStr{
						common.BKHostOuterIPField: exactIP,
					})
					continue
				}
				exactOr = append(exactOr, mapstr.MapStr{
					common.BKCloudIDField:     cloudID,
					common.BKHostOuterIPField: exactIP,
				})

			case IOBOTH:
				if cloudID == -1 {
					exactOr = append(exactOr,
						mapstr.MapStr{
							common.BKHostInnerIPField: exactIP,
						},
						mapstr.MapStr{
							common.BKHostOuterIPField: exactIP,
						})
					continue
				}
				exactOr = append(exactOr,
					mapstr.MapStr{
						common.BKCloudIDField:     cloudID,
						common.BKHostInnerIPField: exactIP,
					},
					mapstr.MapStr{
						common.BKCloudIDField:     cloudID,
						common.BKHostOuterIPField: exactIP,
					})

			default:
				return exactOr, fmt.Errorf("unsupported ip.flag %s", flag)
			}
		}
	}

	return exactOr, nil
}

// addIPv6ExactSearchCondition combine query statements based on exact ipv6 conditions
func addIPv6ExactSearchCondition(exactOr []map[string]interface{}, ipv6CloudIDMap map[int64]map[string][]string,
	flag string) ([]map[string]interface{}, error) {

	for cloudID, ipv6OperatorMap := range ipv6CloudIDMap {
		for operator, ipv6Arr := range ipv6OperatorMap {
			// 校验、去重
			ipv6Arr = filterHostIP(ipv6Arr)
			exactIP := map[string]interface{}{operator: deduplication(ipv6Arr)}
			switch flag {
			case INNERONLY:
				if cloudID == -1 {
					exactOr = append(exactOr, mapstr.MapStr{
						common.BKHostInnerIPv6Field: exactIP,
					})
					continue
				}
				exactOr = append(exactOr, mapstr.MapStr{
					common.BKCloudIDField:       cloudID,
					common.BKHostInnerIPv6Field: exactIP,
				})

			case OUTERONLY:
				if cloudID == -1 {
					exactOr = append(exactOr, mapstr.MapStr{
						common.BKHostOuterIPv6Field: exactIP,
					})
					continue
				}
				exactOr = append(exactOr, mapstr.MapStr{
					common.BKCloudIDField:       cloudID,
					common.BKHostOuterIPv6Field: exactIP,
				})

			case IOBOTH:
				if cloudID == -1 {
					exactOr = append(exactOr,
						mapstr.MapStr{
							common.BKHostInnerIPv6Field: exactIP,
						},
						mapstr.MapStr{
							common.BKHostOuterIPv6Field: exactIP,
						})
					continue
				}
				exactOr = append(exactOr,
					mapstr.MapStr{
						common.BKCloudIDField:       cloudID,
						common.BKHostInnerIPv6Field: exactIP,
					},
					mapstr.MapStr{
						common.BKCloudIDField:       cloudID,
						common.BKHostOuterIPv6Field: exactIP,
					})

			default:
				return exactOr, fmt.Errorf("unsupported ip.flag %s", flag)
			}
		}
	}

	return exactOr, nil
}

// addFuzzyCondition combine query statements based on inexact ip conditions
func addFuzzyCondition(orCond []map[string]map[string]interface{}, ipCloudIDMap map[int64]map[string][]string,
	flag string) ([]map[string]map[string]interface{}, error) {

	for cloudID, ipOperatorMap := range ipCloudIDMap {
		cloudIDOperator := make(map[string]interface{})
		cloudIDOperator[common.BKDBEQ] = cloudID
		for _, ipArr := range ipOperatorMap {
			for _, ip := range ipArr {
				c := make(map[string]interface{})
				c[common.BKDBLIKE] = SpecialCharChange(ip)
				switch flag {
				case INNERONLY:
					if cloudID == -1 {
						orCond = append(orCond, map[string]map[string]interface{}{
							common.BKHostInnerIPField: c,
						})
						continue
					}
					orCond = append(orCond, map[string]map[string]interface{}{
						common.BKCloudIDField:     cloudIDOperator,
						common.BKHostInnerIPField: c,
					})

				case OUTERONLY:
					if cloudID == -1 {
						orCond = append(orCond, map[string]map[string]interface{}{
							common.BKHostOuterIPField: c,
						})
						continue
					}
					orCond = append(orCond, map[string]map[string]interface{}{
						common.BKCloudIDField:     cloudIDOperator,
						common.BKHostOuterIPField: c,
					})

				case IOBOTH:
					if cloudID == -1 {
						orCond = append(orCond, []map[string]map[string]interface{}{
							{
								common.BKHostOuterIPField: c,
							},
							{
								common.BKHostInnerIPField: c,
							}}...)
						continue
					}
					orCond = append(orCond, []map[string]map[string]interface{}{
						{
							common.BKCloudIDField:     cloudIDOperator,
							common.BKHostOuterIPField: c,
						},
						{
							common.BKCloudIDField:     cloudIDOperator,
							common.BKHostInnerIPField: c,
						}}...)

				default:
					return orCond, fmt.Errorf("unsupported ip.flag %s", flag)
				}
			}
		}
	}
	return orCond, nil
}
