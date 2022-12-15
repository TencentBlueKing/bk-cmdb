/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package common

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"reflect"
	"strings"
)

// ConvertHostIpv6Val convert host ipv6 value
func ConvertHostIpv6Val(items []string) ([]string, error) {
	var err error
	for idx, val := range items {
		items[idx], err = ConvertIPv6ToStandardFormat(val)
		if err != nil {
			return nil, err
		}
	}

	return items, nil
}

// convertIpv6ToInt convert an ipv6 address to big int value
func convertIpv6ToInt(ipv6 string) (*big.Int, error) {
	ip := net.ParseIP(ipv6)
	if ip == nil {
		return nil, fmt.Errorf("invalid ipv6 address, data: %s", ipv6)
	}
	intVal := big.NewInt(0).SetBytes(ip)
	return intVal, nil
}

// convertIPv6ToFullAddr convert an ipv6 address to a full ipv6 address
func convertIPv6ToFullAddr(ipv6 string) (string, error) {
	if !strings.Contains(ipv6, ":") {
		return "", fmt.Errorf("address %s is not ipv6 address", ipv6)
	}

	intVal, err := convertIpv6ToInt(ipv6)
	if err != nil {
		return "", err
	}

	b255 := new(big.Int).SetBytes([]byte{255})
	buf := make([]byte, 2)
	part := make([]string, 8)
	pos := 0
	tmpInt := new(big.Int)
	var i uint
	for i = 0; i < 16; i += 2 {
		tmpInt.Rsh(intVal, 120-i*8).And(tmpInt, b255)
		bytes := tmpInt.Bytes()
		if len(bytes) > 0 {
			buf[0] = bytes[0]
		} else {
			buf[0] = 0
		}
		tmpInt.Rsh(intVal, 120-(i+1)*8).And(tmpInt, b255)
		bytes = tmpInt.Bytes()
		if len(bytes) > 0 {
			buf[1] = bytes[0]
		} else {
			buf[1] = 0
		}
		part[pos] = hex.EncodeToString(buf)
		pos++
	}

	return strings.Join(part, ":"), nil
}

// ConvertIPv6ToStandardFormat convert ipv6 address to standard format
// :: => 0000:0000:0000:0000:0000:0000:0000:0000
// ::127.0.0.1 => 0000:0000:0000:0000:0000:0000:127.0.0.1
func ConvertIPv6ToStandardFormat(address string) (string, error) {
	if ip := net.ParseIP(address); ip == nil {
		return "", fmt.Errorf("address %s is invalid", address)
	}

	if !strings.Contains(address, ":") {
		return "", fmt.Errorf("address %s is not ipv6 address", address)
	}

	ipv6FullAddr, err := convertIPv6ToFullAddr(address)
	if err != nil {
		return "", err
	}

	addrs := strings.Split(address, ":")
	if !strings.Contains(addrs[len(addrs)-1], ".") {
		return ipv6FullAddr, nil
	}

	if ip := net.ParseIP(addrs[len(addrs)-1]); ip == nil {
		return "", fmt.Errorf("address %s is invalid", address)
	}

	ipv6FullAddrs := strings.Split(ipv6FullAddr, ":")
	var result string
	for i := 0; i <= len(ipv6FullAddrs)-3; i++ {
		result += ipv6FullAddrs[i] + ":"
	}
	return result + addrs[len(addrs)-1], nil
}

// GetIPv4IfEmbeddedInIPv6 get ipv4 address if it is embedded in ipv6 address
// ::ffff:127.0.0.1 => 127.0.0.1, ::127.0.0.1 => 127.0.0.1
func GetIPv4IfEmbeddedInIPv6(address string) (string, error) {
	if ip := net.ParseIP(address); ip == nil {
		return "", fmt.Errorf("address %s is invalid", address)
	}

	if !strings.Contains(address, ":") {
		return "", fmt.Errorf("address %s is not ipv6 address", address)
	}

	ipv6Addr, err := convertIPv6ToFullAddr(address)
	if err != nil {
		return "", err
	}
	ipv6Addrs := strings.Split(ipv6Addr, ":")
	for i := 0; i <= len(ipv6Addrs)-3; i++ {
		if i != len(ipv6Addrs)-3 && ipv6Addrs[i] != "0000" {
			return address, nil
		}

		if i == len(ipv6Addrs)-3 && ipv6Addrs[i] != "0000" && ipv6Addrs[i] != "ffff" {
			return address, nil
		}
	}

	addrs := strings.Split(address, ":")
	if !strings.Contains(addrs[len(addrs)-1], ".") {
		return address, nil
	}

	if ip := net.ParseIP(addrs[len(addrs)-1]); ip == nil {
		return "", fmt.Errorf("address %s is invalid", address)
	}

	return addrs[len(addrs)-1], nil
}

// ConvertIpv6ToFullWord convert the ipv6 field into a complete format.
// for the converted scene at this time, there are only two types of value,
// string and slice, because the operators involved in mongo can only be
// one of the four cases "$eq", "$ne", "$in" and "$nin".
func ConvertIpv6ToFullWord(field string, value interface{}) (interface{}, error) {

	// only supports the transformation of the ipv6 field of the internal and external network.
	if field != BKHostInnerIPv6Field && field != BKHostOuterIPv6Field {
		return value, nil
	}

	var data interface{}
	switch reflect.ValueOf(value).Kind() {
	case reflect.String:
		ip, err := ConvertIPv6ToStandardFormat(value.(string))
		if err != nil {
			return nil, err
		}
		data = ip
	case reflect.Array, reflect.Slice:
		v := reflect.ValueOf(value)
		length := v.Len()
		if length == 0 {
			return value, nil
		}

		result := make([]interface{}, 0)
		// each element in the array or slice should be of the same basic type.
		for i := 0; i < length; i++ {
			item := v.Index(i).Interface()
			switch item.(type) {
			case string:
				v, err := ConvertIPv6ToStandardFormat(item.(string))
				if err != nil {
					return nil, err
				}
				result = append(result, v)
			default:
				return value, nil
			}
		}
		data = result
	default:
		return value, nil
	}

	return data, nil
}
