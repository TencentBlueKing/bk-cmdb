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
	"configcenter/src/common/blog"
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
	// NOCLOUDID 未指定云区域时map的键
	NOCLOUDID int64 = -1
	// NOTEXISTSIGN 分割符不存在
	NOTEXISTSIGN int = -1
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
		exactOr, embeddedIPv4Addrs, err = parseIPv6Condition(ipv6Cond, exactOr, output)
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

	ipv4CloudIDMap, err := splitIPv4Data(ipv4Cond)
	if err != nil {
		return nil, err
	}

	if exact == 1 {
		// exact search
		// filter out illegal IPv4 addresses
		exactOr, err = addIPv4ExactSearchCondition(exactOr, ipv4CloudIDMap, output, flag)
		if err != nil {
			return nil, err
		}
		output[common.BKDBOR] = exactOr
	} else {
		// not exact search
		orCond := make([]map[string]interface{}, 0)
		orCond, err = addFuzzyCondition(orCond, ipv4CloudIDMap, output, flag)
		if err != nil {
			return nil, err
		}
		output[common.BKDBOR] = orCond
	}
	return output, nil
}

// parseIPv6Condition parse IPv6 conditions to full Ipv6 addresses and embedded IPv4 addresses
// only full or abbreviated IPv6 addresses can be used for exact queries, not exact search is not supported
func parseIPv6Condition(ipCond metadata.IPInfo, exactOr []map[string]interface{}, output map[string]interface{}) (
	[]map[string]interface{}, []string, error) {

	flag := ipCond.Flag
	if ipCond.Exact != 1 {
		return exactOr, nil, nil
	}

	ipv6CloudIDMap, embeddedIPv4Addrs, err := splitIPv6Data(ipCond)
	if err != nil {
		return nil, nil, err
	}

	exactOr, err = addIPv6ExactSearchCondition(exactOr, ipv6CloudIDMap, output, flag)
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
/*
该方法用于分割出IPv4条件中IP数组的管控区域和IPv4，返回存储管控区域与IP的map用于构建查询条件，如下IP条件数组：
"data": [
	"1:172.0.0.1",
	"1:172.0.0.2",
	"2:172.0.0.3",
	"172.0.0.4"
]

处理后返回的map为：
{
	1: ["172.0.0.1", "172.0.0.2"],
	2: ["172.0.0.3"],
	-1: ["172.0.0.4"]
}
*/
func splitIPv4Data(ipCond metadata.IPInfo) (map[int64][]string, error) {
	// 创建一个 map 用于存储分割结果
	cloudIDMap := make(map[int64][]string)

	// 遍历 IPInfo 中的每个字符串
	for _, ipString := range ipCond.Data {
		// 去除空格
		ipString = strings.TrimSpace(ipString)

		// 分割管控区域ID 和 IP地址
		colonIndex := strings.Index(ipString, ":")
		// 如果没有 ":" 代表未直接指定管控区域，则直接添加到 NOCLOUDID 键下
		if colonIndex == NOTEXISTSIGN {
			cloudIDMap[NOCLOUDID] = append(cloudIDMap[NOCLOUDID], ipString)
			continue
		}

		// 去掉 IP地址中的方括号，只保留 IP
		ipAddress := strings.Trim(ipString[colonIndex+1:], "[]")
		// 转换管控区域ID 为 int64 类型
		cloudIDStr := ipString[:colonIndex]
		cloudIDInt64, err := strconv.ParseInt(cloudIDStr, 10, 64)
		if err != nil {
			blog.Errorf("cloudID is invalid failed, err: %v, cloudID: %s", err, cloudIDStr)
			continue
		}

		cloudIDMap[cloudIDInt64] = append(cloudIDMap[cloudIDInt64], ipAddress)
	}
	return cloudIDMap, nil
}

// splitIPv6Data 该方法用于分割出IPv6条件中IP数组的管控区域和IPv6，返回存储管控区域与IP的map用于构建查询条件，具体解释参考splitIPv4Data方法解释
func splitIPv6Data(ipCond metadata.IPInfo) (map[int64][]string, []string, error) {
	// 创建一个 map 用于存储分割结果
	cloudIDMap := make(map[int64][]string)
	embeddedIPv4Addrs := make([]string, 0)

	// 遍历切片中的每个字符串
	for _, ipString := range ipCond.Data {
		// 去除空格
		ipString = strings.TrimSpace(ipString)

		// 分割管控区域 ID 和 IP 地址
		colonIndex := strings.Index(ipString, ":[")
		ipPart := ipString[colonIndex+1:]
		if colonIndex == NOTEXISTSIGN {
			ipPart = ipString
		}
		// 去掉 IP 地址中的方括号，只保留 IP
		ipAddress := strings.Trim(ipPart, "[]")

		// 对于兼容IPv4的嵌入式IPv6地址，::127.0.0.1和::ffff:127.0.0.1这两种格式的地址，存放于ipv4字段中，所以使用ipv4的字段查询
		ipAddr, err := common.GetIPv4IfEmbeddedInIPv6(ipAddress)
		if err != nil {
			blog.Errorf("get ipv4 if embedded in ipv6 failed, err: %v, ip: %s", err, ipAddress)
			continue
		}
		if !strings.Contains(ipAddr, ":") {
			embeddedIPv4Addrs = append(embeddedIPv4Addrs, ipString)
			continue
		}

		fullIpv6Addr, err := common.ConvertIPv6ToStandardFormat(ipAddr)
		if err != nil {
			blog.Errorf("convert ipv6 to standard format failed, err: %v, ip: %s", err, ipAddr)
			continue
		}

		// 如果没有 ":[" 代表未直接指定管控区域，则直接添加到 NOCLOUDID 键下
		if colonIndex == NOTEXISTSIGN {
			cloudIDMap[NOCLOUDID] = append(cloudIDMap[NOCLOUDID], ipString)
			continue
		}

		// 转换管控区域 ID 为 int64 类型
		cloudIDStr := ipString[:colonIndex]
		cloudIDInt64, err := strconv.ParseInt(cloudIDStr, 10, 64)
		if err != nil {
			blog.Errorf("cloudID is invalid failed, err: %v, cloudID: %s", err, cloudIDStr)
			continue
		}

		cloudIDMap[cloudIDInt64] = append(cloudIDMap[cloudIDInt64], fullIpv6Addr)
	}
	return cloudIDMap, embeddedIPv4Addrs, nil
}

// getCloudIDMapByOutput 从output中获取管控区域id添加到outputCloudIDMap
func getCloudIDMapByOutput(output map[string]interface{}) (map[string][]int64, error) {

	outputCloudIDMap := make(map[string][]int64)
	var err error

	cloudIDCond, ok := output[common.BKCloudIDField].(map[string]interface{})
	if !ok {
		return nil, nil
	}
	hasOperator := false

	if cloudIDArray, ok := cloudIDCond[common.BKDBIN].([]interface{}); ok {
		hasOperator = true
		outputCloudIDMap, err = addCloudIDInMap(cloudIDArray, outputCloudIDMap)
		if err != nil {
			return nil, err
		}
	}
	if cloudIDArray, ok := cloudIDCond[common.BKDBNIN].([]interface{}); ok {
		hasOperator = true
		outputCloudIDMap, err = addCloudIDInMap(cloudIDArray, outputCloudIDMap)
		if err != nil {
			return nil, err
		}
	}

	if !hasOperator {
		return nil, fmt.Errorf("cloudID query condition is not equal to '$in' or '$nin' failed")
	}
	return outputCloudIDMap, nil
}

// addCloudIDInMap 添加管控区域id到outputCloudIDMap
func addCloudIDInMap(cloudIDArray []interface{}, outputCloudIDMap map[string][]int64) (map[string][]int64, error) {
	for _, cloudID := range cloudIDArray {
		cloudIDStr, ok := cloudID.(json.Number)
		if !ok {
			return nil, fmt.Errorf("conversion to type 'json.Number' failed")
		}

		cloudIDInt64, err := cloudIDStr.Int64()
		if err != nil {
			return nil, err
		}

		outputCloudIDMap[common.BKDBNIN] = append(outputCloudIDMap[common.BKDBNIN], cloudIDInt64)
	}
	return outputCloudIDMap, nil
}

// addExactSearchCondition combine query statements based on exact ipv4 conditions
func addIPv4ExactSearchCondition(exactOr []map[string]interface{}, ipv4CloudIDMap map[int64][]string,
	output map[string]interface{}, flag string) ([]map[string]interface{}, error) {

	outputCloudIDMap, err := getCloudIDMapByOutput(output)
	if err != nil {
		return nil, err
	}

	for cloudID, ipv4Arr := range ipv4CloudIDMap {
		// 校验、去重
		ipv4Arr = filterHostIP(ipv4Arr)
		exactIP := map[string]interface{}{common.BKDBIN: deduplication(ipv4Arr)}

		ipv4MapCond, err := getIPv4MapCond(flag, cloudID, outputCloudIDMap, exactIP)
		if err != nil {
			return nil, err
		}
		exactOr = append(exactOr, ipv4MapCond...)
	}

	return exactOr, nil
}

// addIPv6ExactSearchCondition combine query statements based on exact ipv6 conditions
func addIPv6ExactSearchCondition(exactOr []map[string]interface{}, ipv6CloudIDMap map[int64][]string,
	output map[string]interface{}, flag string) ([]map[string]interface{}, error) {

	outputCloudIDMap, err := getCloudIDMapByOutput(output)
	if err != nil {
		return nil, err
	}

	for cloudID, ipv6Arr := range ipv6CloudIDMap {
		// 校验、去重
		ipv6Arr = filterHostIP(ipv6Arr)
		exactIP := map[string]interface{}{common.BKDBIN: deduplication(ipv6Arr)}

		ipv6MapCond, err := getIPv6MapCond(flag, cloudID, outputCloudIDMap, exactIP)
		if err != nil {
			return nil, err
		}
		exactOr = append(exactOr, ipv6MapCond...)
	}

	return exactOr, nil
}

// getIPv4MapCond 获取ipv4相关的查询条件
/*
该方法用于构建ipv4的精确查询条件并返回
1、当参数中的IP前直接指定了管控区域如："4:127.0.0.1"，则构建的条件：
{
	"bk_cloud_id": 4,
	"bk_host_innerip": {"$in": ["127.0.0.1", "......其它管控区域相同的IP"]}
}

2、当参数中的IP未直接指定了管控区域如："127.0.0.1"，而在主机查询条件中指定了管控区域[4,5,6]，则构建的条件：
{
	"bk_cloud_id": {"$in": [4,5,6]},
	"bk_host_innerip": {"$in": ["127.0.0.1", "......其它未直接指定管控区域的IP"]}
}

3、当参数中的IP未直接指定了管控区域如："127.0.0.1"，而在主机查询条件中也未指定管控区域，则构建的条件：
{
	"bk_host_innerip": {"$in": ["127.0.0.1", "......其它未直接指定管控区域的IP"]}
}
*/
func getIPv4MapCond(flag string, cloudID int64, outputCloudIDMap map[string][]int64,
	exactIP map[string]interface{}) ([]map[string]interface{}, error) {

	cloudIDArr, ok := outputCloudIDMap[common.BKDBIN]
	if ok {
		cloudIDArr = util.IntArrayUnique(cloudIDArr)
	}
	switch flag {
	case INNERONLY:
		// 直接指定管控区域的IP
		if cloudID != NOCLOUDID {
			return []map[string]interface{}{{
				common.BKCloudIDField:     cloudID,
				common.BKHostInnerIPField: exactIP,
			}}, nil
		}

		// 在hostCond中有指定
		if ok {
			return []map[string]interface{}{{
				common.BKCloudIDField:     map[string]interface{}{common.BKDBIN: cloudIDArr},
				common.BKHostInnerIPField: exactIP,
			}}, nil
		}
		// 在hostCond中未指定
		return []map[string]interface{}{{
			common.BKHostInnerIPField: exactIP,
		}}, nil

	case OUTERONLY:
		if cloudID != NOCLOUDID {
			return []map[string]interface{}{{
				common.BKCloudIDField:     cloudID,
				common.BKHostOuterIPField: exactIP,
			}}, nil
		}

		if ok {
			return []map[string]interface{}{{
				common.BKCloudIDField:     map[string]interface{}{common.BKDBIN: cloudIDArr},
				common.BKHostOuterIPField: exactIP,
			}}, nil
		}
		return []map[string]interface{}{{
			common.BKHostOuterIPField: exactIP,
		}}, nil

	case IOBOTH:
		if cloudID != NOCLOUDID {
			return []map[string]interface{}{{
				common.BKCloudIDField:     cloudID,
				common.BKHostInnerIPField: exactIP,
			}, {
				common.BKCloudIDField:     cloudID,
				common.BKHostOuterIPField: exactIP,
			}}, nil

		}

		if ok {
			return []map[string]interface{}{{
				common.BKCloudIDField:     map[string]interface{}{common.BKDBIN: cloudIDArr},
				common.BKHostInnerIPField: exactIP,
			}, {
				common.BKCloudIDField:     map[string]interface{}{common.BKDBIN: cloudIDArr},
				common.BKHostOuterIPField: exactIP,
			}}, nil
		}
		return []map[string]interface{}{{
			common.BKHostInnerIPField: exactIP,
		}, {
			common.BKHostOuterIPField: exactIP,
		}}, nil

	default:
		return nil, fmt.Errorf("unsupported ip.flag %s", flag)
	}
}

// getIPv6MapStrCond 获取ipv6相关的查询条件，该方法用于构建ipv6的精确查询条件并返回，条件样式参考getIPv4MapCond方法解释
func getIPv6MapCond(flag string, cloudID int64, outputCloudIDMap map[string][]int64,
	exactIP map[string]interface{}) ([]map[string]interface{}, error) {

	cloudIDArr, ok := outputCloudIDMap[common.BKDBIN]
	if ok {
		cloudIDArr = util.IntArrayUnique(cloudIDArr)
	}
	switch flag {
	case INNERONLY:
		// 直接指定管控区域的IP
		if cloudID != NOCLOUDID {
			return []map[string]interface{}{{
				common.BKCloudIDField:       cloudID,
				common.BKHostInnerIPv6Field: exactIP,
			}}, nil
		}

		// 在hostCond中有指定
		if ok {
			return []map[string]interface{}{{
				common.BKCloudIDField:       map[string]interface{}{common.BKDBIN: cloudIDArr},
				common.BKHostInnerIPv6Field: exactIP,
			}}, nil
		}
		// 在hostCond中未指定
		return []map[string]interface{}{{
			common.BKHostInnerIPv6Field: exactIP,
		}}, nil

	case OUTERONLY:
		if cloudID != NOCLOUDID {
			return []map[string]interface{}{{
				common.BKCloudIDField:       cloudID,
				common.BKHostOuterIPv6Field: exactIP,
			}}, nil
		}

		if ok {
			return []map[string]interface{}{{
				common.BKCloudIDField:       map[string]interface{}{common.BKDBIN: cloudIDArr},
				common.BKHostOuterIPv6Field: exactIP,
			}}, nil
		}
		return []map[string]interface{}{{
			common.BKHostOuterIPv6Field: exactIP,
		}}, nil

	case IOBOTH:
		if cloudID != NOCLOUDID {
			return []map[string]interface{}{{
				common.BKCloudIDField:       cloudID,
				common.BKHostInnerIPv6Field: exactIP,
			}, {
				common.BKCloudIDField:       cloudID,
				common.BKHostOuterIPv6Field: exactIP,
			}}, nil
		}

		if ok {
			return []map[string]interface{}{{
				common.BKCloudIDField:       map[string]interface{}{common.BKDBIN: cloudIDArr},
				common.BKHostInnerIPv6Field: exactIP,
			}, {
				common.BKCloudIDField:       map[string]interface{}{common.BKDBIN: cloudIDArr},
				common.BKHostOuterIPv6Field: exactIP,
			}}, nil
		}
		return []map[string]interface{}{{
			common.BKHostInnerIPv6Field: exactIP,
		}, {
			common.BKHostOuterIPv6Field: exactIP,
		}}, nil

	default:
		return nil, fmt.Errorf("unsupported ip.flag %s", flag)
	}
}

// addFuzzyCondition combine query statements based on inexact ip conditions
func addFuzzyCondition(orCond []map[string]interface{}, ipCloudIDMap map[int64][]string,
	output map[string]interface{}, flag string) ([]map[string]interface{}, error) {

	outputCloudIDMap, err := getCloudIDMapByOutput(output)
	if err != nil {
		return nil, err
	}

	for cloudID, ipArr := range ipCloudIDMap {
		ipArr = deduplication(ipArr)

		for _, ip := range ipArr {
			ipRegex := make(map[string]interface{})
			ipRegex[common.BKDBLIKE] = SpecialCharChange(ip)

			fuzzySearchCond, err := getFuzzyCond(flag, cloudID, outputCloudIDMap, ipRegex)
			if err != nil {
				return nil, err
			}
			orCond = append(orCond, fuzzySearchCond...)
		}
	}
	return orCond, nil
}

// getFuzzyMapStrCond 获取模糊查询条件
/*
该方法用于构建ipv4的模糊查询条件并返回，ipv6不支持模糊查询
1、当参数中的IP前直接指定了管控区域如："4:127.0."，则构建的条件：
{
	"bk_cloud_id": 4,
	"bk_host_innerip": {"$regex": "127.0."}
}

2、当参数中的IP未直接指定了管控区域如："127.0."，而在主机查询条件中指定了管控区域[4,5,6]，则构建的条件：
{
	"bk_cloud_id": {"$in": [4,5,6]},
	"bk_host_innerip": {"$regex": "127.0."}
}

3、当参数中的IP未直接指定了管控区域如："127.0."，而在主机查询条件中也未指定管控区域，则构建的条件：
{
	"bk_host_innerip": {"$regex": "127.0."}
}
*/
func getFuzzyCond(flag string, cloudID int64, outputCloudIDMap map[string][]int64,
	ipRegex map[string]interface{}) ([]map[string]interface{}, error) {

	cloudIDArr, ok := outputCloudIDMap[common.BKDBIN]
	if ok {
		cloudIDArr = util.IntArrayUnique(cloudIDArr)
	}
	switch flag {
	case INNERONLY:
		if cloudID != NOCLOUDID {
			return []map[string]interface{}{{
				common.BKCloudIDField:     cloudID,
				common.BKHostInnerIPField: ipRegex,
			}}, nil
		}
		if ok {
			return []map[string]interface{}{{
				common.BKCloudIDField:     map[string]interface{}{common.BKDBIN: cloudIDArr},
				common.BKHostInnerIPField: ipRegex,
			}}, nil
		}
		return []map[string]interface{}{{
			common.BKHostInnerIPField: ipRegex,
		}}, nil

	case OUTERONLY:
		if cloudID != NOCLOUDID {
			return []map[string]interface{}{{
				common.BKCloudIDField:     cloudID,
				common.BKHostOuterIPField: ipRegex,
			}}, nil

		}
		if ok {
			return []map[string]interface{}{{
				common.BKCloudIDField:     map[string]interface{}{common.BKDBIN: cloudIDArr},
				common.BKHostOuterIPField: ipRegex,
			}}, nil
		}
		return []map[string]interface{}{{
			common.BKHostOuterIPField: ipRegex,
		}}, nil

	case IOBOTH:
		if cloudID != NOCLOUDID {
			return []map[string]interface{}{{
				common.BKCloudIDField:     cloudID,
				common.BKHostInnerIPField: ipRegex,
			}, {
				common.BKCloudIDField:     cloudID,
				common.BKHostOuterIPField: ipRegex,
			}}, nil

		}

		if ok {
			return []map[string]interface{}{{
				common.BKCloudIDField:     map[string]interface{}{common.BKDBIN: cloudIDArr},
				common.BKHostInnerIPField: ipRegex,
			}, {
				common.BKCloudIDField:     map[string]interface{}{common.BKDBIN: cloudIDArr},
				common.BKHostOuterIPField: ipRegex,
			}}, nil
		}
		return []map[string]interface{}{{
			common.BKHostInnerIPField: ipRegex,
		}, {
			common.BKHostOuterIPField: ipRegex,
		}}, nil

	default:
		return nil, fmt.Errorf("unsupported ip.flag %s", flag)
	}
}
