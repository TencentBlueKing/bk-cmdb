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
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
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
	// TYPEIPV4 IPv4类型
	TYPEIPV4 string = "ipv4"
	// TYPEIPV6 IPv6类型
	TYPEIPV6 string = "ipv6"
)

// FieldAndCondition 存放构建IP条件时的字段名称、管控区域条件和IP条件
type FieldAndCondition struct {
	CloudIDField     string
	HostInnerIPField string
	HostOuterIPField string
	CloudIDCond      map[string]interface{}
	IPCond           map[string]interface{}
}

// SplitIPResult 存放IP与管控区域分割后结果
type SplitIPResult struct {
	CloudIdIPMap      map[int64][]string
	EmbeddedIPv4Addrs []string
	NoCloudIdIP       []string
}

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
func ParseHostIPParams(ipv4Cond metadata.IPInfo, ipv6Cond metadata.IPInfo, output map[string]interface{},
	rid string) (map[string]interface{}, error) {

	exact := ipv4Cond.Exact
	if exact != 1 && len(ipv4Cond.Data) > 10 {
		return nil, fmt.Errorf("the number of IP condition in fuzzy query cannot more than 10")
	}
	var err error
	exactOr := make([]map[string]interface{}, 0)
	embeddedIPv4Addrs := make([]string, 0)

	if len(ipv6Cond.Data) != 0 {
		exactOr, embeddedIPv4Addrs, err = parseIPv6Condition(ipv6Cond, exactOr, output, rid)
		if err != nil {
			return nil, fmt.Errorf("failed to add ipv6 addresses to condition, err: %v", err)
		}
	}

	ipv4Cond.Data = append(ipv4Cond.Data, embeddedIPv4Addrs...)
	flag := ipv4Cond.Flag
	if len(ipv4Cond.Data) == 0 && len(exactOr) == 0 {
		return output, nil
	}

	splitIPResult, err := splitIPData(ipv4Cond, TYPEIPV4, rid)
	if err != nil {
		return nil, err
	}
	if len(splitIPResult.CloudIdIPMap) == 0 && len(splitIPResult.NoCloudIdIP) == 0 && len(exactOr) != 0 {
		output[common.BKDBOR] = exactOr
		return output, nil
	}

	fieldAndCond := FieldAndCondition{
		CloudIDField:     common.BKCloudIDField,
		HostInnerIPField: common.BKHostInnerIPField,
		HostOuterIPField: common.BKHostOuterIPField,
	}
	if exact == 1 {
		// exact search
		// filter out illegal IPv4 addresses
		ipv4ExactCond, err := addExactSearchCondition(splitIPResult, fieldAndCond, output, flag)
		if err != nil {
			return nil, err
		}
		exactOr = append(exactOr, ipv4ExactCond...)
		// 此处判断当设置了ip条件如：5h:127.0.0.x, j:127.0.0.1 但因为设置的所有ip或管控区域不合法，导致最后构建的条件为空，从而会查询出所有主机，不符合用户预期
		// 所以如果设置了ip条件，但所有ip或管控区域都不合法，就设置一个空查询从而不返回任何主机数据
		if len(exactOr) == 0 {
			exactOr = append(exactOr, []map[string]interface{}{{
				common.BKHostInnerIPField: map[string]interface{}{common.BKDBIN: []string{}},
			}}...)
		}
		output[common.BKDBOR] = exactOr
	} else {
		// not exact search
		orCond := make([]map[string]interface{}, 0)
		orCond, err = addFuzzyCondition(splitIPResult, fieldAndCond, output, flag)
		if err != nil {
			return nil, err
		}
		output[common.BKDBOR] = orCond
	}
	return output, nil
}

// parseIPv6Condition parse IPv6 conditions to full Ipv6 addresses and embedded IPv4 addresses
// only full or abbreviated IPv6 addresses can be used for exact queries, not exact search is not supported
func parseIPv6Condition(ipCond metadata.IPInfo, exactOr []map[string]interface{}, output map[string]interface{},
	rid string) ([]map[string]interface{}, []string, error) {

	flag := ipCond.Flag
	if ipCond.Exact != 1 {
		return exactOr, nil, nil
	}

	splitIPResult, err := splitIPData(ipCond, TYPEIPV6, rid)
	if err != nil {
		return nil, nil, err
	}
	if len(splitIPResult.CloudIdIPMap) == 0 && len(splitIPResult.NoCloudIdIP) == 0 {
		return exactOr, splitIPResult.EmbeddedIPv4Addrs, nil
	}

	fieldAndCond := FieldAndCondition{
		CloudIDField:     common.BKCloudIDField,
		HostInnerIPField: common.BKHostInnerIPv6Field,
		HostOuterIPField: common.BKHostOuterIPv6Field,
	}
	exactOr, err = addExactSearchCondition(splitIPResult, fieldAndCond, output, flag)
	if err != nil {
		return nil, nil, err
	}
	return exactOr, splitIPResult.EmbeddedIPv4Addrs, nil
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

// splitIPData 该方法用于分割出IP条件中IP数组的管控区域和IP，
// 返回存储管控区域与IP的map以及兼容IPv4的嵌入式IPv6地址的切片
/*
如下IP条件数组：
"data": [
	"1:127.0.0.1",
	"1:127.0.0.2",
	"2:127.0.0.3",
	"127.0.0.4"
]

处理后返回的map为：
{
	1: ["127.0.0.1", "127.0.0.2"],
	2: ["127.0.0.3"],
	-1: ["127.0.0.4"]
}
*/
func splitIPData(ipCond metadata.IPInfo, ipType string, rid string) (SplitIPResult, error) {
	// 创建一个 map 用于存储分割结果
	cloudIdIpMap := make(map[int64][]string)
	embeddedIPv4Addrs := make([]string, 0)
	noCloudIdIP := make([]string, 0)
	var colonIndex int
	var fullIPAddr string

	// 遍历切片中的每个字符串
	for _, ipString := range ipCond.Data {
		// 去除空格
		ipString = strings.TrimSpace(ipString)
		// 分割管控区域ID 和 IP地址
		colonIndex = strings.Index(ipString, ":")
		// 去掉 IP地址中的方括号，只保留 IP
		fullIPAddr = strings.Trim(ipString[colonIndex+1:], "[]")

		if ipType == TYPEIPV6 {
			// 分割管控区域 ID 和 IP 地址
			colonIndex = strings.Index(ipString, ":[")
			ipPart := ipString[colonIndex+1:]
			if colonIndex == -1 {
				ipPart = ipString
			}
			// 去掉 IP 地址中的方括号，只保留 IP
			ipAddress := strings.Trim(ipPart, "[]")

			// 对于兼容IPv4的嵌入式IPv6地址，::127.0.0.1和::ffff:127.0.0.1这两种格式的地址，存放于ipv4字段中，所以使用ipv4的字段查询
			ipAddr, err := common.GetIPv4IfEmbeddedInIPv6(ipAddress)
			if err != nil {
				blog.Errorf("get ipv4 if embedded in ipv6 failed, err: %v, ip: %s, rid: %s", err, ipAddress, rid)
				continue
			}
			if !strings.Contains(ipAddr, ":") {
				embeddedIPv4Addrs = append(embeddedIPv4Addrs, ipString)
				continue
			}

			fullIPAddr, err = common.ConvertIPv6ToStandardFormat(ipAddr)
			if err != nil {
				blog.Errorf("convert ipv6 to standard format failed, err: %v, ip: %s, rid: %s", err, ipAddr, rid)
				continue
			}
		}

		// 如果没有 ":[" 代表未直接指定管控区域，则直接添加到 noCloudIdIP下
		if colonIndex == -1 {
			noCloudIdIP = append(noCloudIdIP, ipString)
			continue
		}

		// 转换管控区域 ID 为 int64 类型
		cloudIDStr := ipString[:colonIndex]
		cloudIDInt64, err := strconv.ParseInt(cloudIDStr, 10, 64)
		if err != nil {
			blog.Errorf("cloudID is invalid failed, err: %v, cloudID: %s, rid: %s", err, cloudIDStr, rid)
			continue
		}
		cloudIdIpMap[cloudIDInt64] = append(cloudIdIpMap[cloudIDInt64], fullIPAddr)
	}

	splitIPResult := SplitIPResult{
		CloudIdIPMap:      cloudIdIpMap,
		EmbeddedIPv4Addrs: embeddedIPv4Addrs,
		NoCloudIdIP:       noCloudIdIP,
	}
	return splitIPResult, nil
}

// addExactSearchCondition combine query statements based on exact ip conditions
/*
该方法用于构建ipv4和ipv6的精确查询条件并返回
获取ipv4相关的查询条件：
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

获取ipv6相关的查询条件：
1、当参数中的IP前直接指定了管控区域如："4:[0000:0000:0000:0000:0000:0000:0000:1234]"，则构建的条件：
{
	"bk_cloud_id": 4,
	"bk_host_innerip_v6": {"$in": ["0000:0000:0000:0000:0000:0000:0000:1234", "......其它管控区域相同的IP"]}
}

2、当参数中的IP未直接指定了管控区域如："127.0.0.1"，而在主机查询条件中指定了管控区域[4,5,6]，则构建的条件：
{
	"bk_cloud_id": {"$in": [4,5,6]},
	"bk_host_innerip_v6": {"$in": ["0000:0000:0000:0000:0000:0000:0000:1234", "......其它未直接指定管控区域的IP"]}
}

3、当参数中的IP未直接指定了管控区域如："0000:0000:0000:0000:0000:0000:0000:1234"，而在主机查询条件中也未指定管控区域，则构建的条件：
{
	"bk_host_innerip_v6": {"$in": ["0000:0000:0000:0000:0000:0000:0000:1234", "......其它未直接指定管控区域的IP"]}
}
*/
func addExactSearchCondition(splitIPResult SplitIPResult, fieldAndCond FieldAndCondition, output map[string]interface{}, flag string) ([]map[string]interface{}, error) {

	exactOr := make([]map[string]interface{}, 0)

	if len(splitIPResult.NoCloudIdIP) != 0 {
		ips := filterHostIP(splitIPResult.NoCloudIdIP)
		exactIP := map[string]interface{}{common.BKDBIN: deduplication(ips)}
		fieldAndCond.IPCond = exactIP
		ipCond, err := getIPCond(flag, output, fieldAndCond)
		if err != nil {
			return nil, err
		}
		exactOr = append(exactOr, ipCond...)
	}

	for cloudID, ipArr := range splitIPResult.CloudIdIPMap {
		fieldAndCond.CloudIDCond = map[string]interface{}{common.BKDBEQ: cloudID}
		ipArr = filterHostIP(ipArr)
		fieldAndCond.IPCond = map[string]interface{}{common.BKDBIN: deduplication(ipArr)}

		ipCond, err := getIPCond(flag, output, fieldAndCond)
		if err != nil {
			return nil, err
		}
		exactOr = append(exactOr, ipCond...)
	}

	return exactOr, nil
}

// addFuzzyCondition combine query statements based on inexact ip conditions
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
func addFuzzyCondition(splitIPResult SplitIPResult, fieldAndCond FieldAndCondition, output map[string]interface{},
	flag string) ([]map[string]interface{}, error) {

	orCond := make([]map[string]interface{}, 0)

	// 处理未设置管控区域的IP
	if len(splitIPResult.NoCloudIdIP) != 0 {
		splitIPResult.NoCloudIdIP = deduplication(splitIPResult.NoCloudIdIP)
		regexString := SpecialCharChange(splitIPResult.NoCloudIdIP[0])
		for i := 1; i < len(splitIPResult.NoCloudIdIP); i++ {
			regexString = regexString + "|" + SpecialCharChange(splitIPResult.NoCloudIdIP[i])
		}

		fieldAndCond.IPCond = map[string]interface{}{common.BKDBLIKE: regexString}
		fuzzySearchCond, err := getIPCond(flag, output, fieldAndCond)
		if err != nil {
			return nil, err
		}
		orCond = append(orCond, fuzzySearchCond...)
	}

	// 处理设置了管控区域的IP
	for cloudID, ipArr := range splitIPResult.CloudIdIPMap {
		ipArr = deduplication(ipArr)
		regexString := SpecialCharChange(ipArr[0])
		for i := 1; i < len(ipArr); i++ {
			regexString = regexString + "|" + SpecialCharChange(ipArr[i])
		}

		fieldAndCond.CloudIDCond = map[string]interface{}{common.BKDBEQ: cloudID}
		fieldAndCond.IPCond = map[string]interface{}{common.BKDBLIKE: regexString}

		fuzzySearchCond, err := getIPCond(flag, output, fieldAndCond)
		if err != nil {
			return nil, err
		}
		orCond = append(orCond, fuzzySearchCond...)

	}
	return orCond, nil
}

// getCloudIDCond 获取output中的管控区域条件
func getCloudIDCond(output map[string]interface{}) (map[string]interface{}, bool, error) {
	cloudIDCond, ok := output[common.BKCloudIDField].(map[string]interface{})
	if ok {
		_, inExist := cloudIDCond[common.BKDBIN]
		_, ninExist := cloudIDCond[common.BKDBNIN]
		if !inExist && !ninExist {
			return nil, false, fmt.Errorf("cloudID query condition is not equal to '$in' or '$nin' failed")
		}
		return cloudIDCond, true, nil
	}
	return cloudIDCond, false, nil
}

// getIPCond 获取ip查询条件
func getIPCond(flag string, output map[string]interface{}, fieldAndCond FieldAndCondition) ([]map[string]interface{},
	error) {

	cloudIDCond, ok, err := getCloudIDCond(output)
	if err != nil {
		return nil, err
	}

	switch flag {
	case INNERONLY:
		if fieldAndCond.CloudIDCond != nil {
			return []map[string]interface{}{{
				fieldAndCond.CloudIDField:     fieldAndCond.CloudIDCond,
				fieldAndCond.HostInnerIPField: fieldAndCond.IPCond,
			}}, nil
		}
		if ok {
			return []map[string]interface{}{{
				fieldAndCond.CloudIDField:     cloudIDCond,
				fieldAndCond.HostInnerIPField: fieldAndCond.IPCond,
			}}, nil
		}
		return []map[string]interface{}{{
			fieldAndCond.HostInnerIPField: fieldAndCond.IPCond,
		}}, nil

	case OUTERONLY:
		if fieldAndCond.CloudIDCond != nil {
			return []map[string]interface{}{{
				fieldAndCond.CloudIDField:     fieldAndCond.CloudIDCond,
				fieldAndCond.HostOuterIPField: fieldAndCond.IPCond,
			}}, nil

		}
		if ok {
			return []map[string]interface{}{{
				fieldAndCond.CloudIDField:     cloudIDCond,
				fieldAndCond.HostOuterIPField: fieldAndCond.IPCond,
			}}, nil
		}
		return []map[string]interface{}{{
			fieldAndCond.HostOuterIPField: fieldAndCond.IPCond,
		}}, nil

	case IOBOTH:
		if fieldAndCond.CloudIDCond != nil {
			return []map[string]interface{}{{
				fieldAndCond.CloudIDField:     fieldAndCond.CloudIDCond,
				fieldAndCond.HostInnerIPField: fieldAndCond.IPCond,
			}, {
				fieldAndCond.CloudIDField:     fieldAndCond.CloudIDCond,
				fieldAndCond.HostOuterIPField: fieldAndCond.IPCond,
			}}, nil
		}

		if ok {
			return []map[string]interface{}{{
				fieldAndCond.CloudIDField:     cloudIDCond,
				fieldAndCond.HostInnerIPField: fieldAndCond.IPCond,
			}, {
				fieldAndCond.CloudIDField:     cloudIDCond,
				fieldAndCond.HostOuterIPField: fieldAndCond.IPCond,
			}}, nil
		}
		return []map[string]interface{}{{
			fieldAndCond.HostInnerIPField: fieldAndCond.IPCond,
		}, {
			fieldAndCond.HostOuterIPField: fieldAndCond.IPCond,
		}}, nil

	default:
		return nil, fmt.Errorf("unsupported ip.flag %s", flag)
	}
}
