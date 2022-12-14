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

// convertIpv6ToInt convert an ipv6 address to big int value
func convertIpv6ToInt(ipv6 string) (*big.Int, error) {
	ip := net.ParseIP(ipv6)
	if ip == nil {
		return nil, fmt.Errorf("invalid ipv6 address, data: %s", ipv6)
	}
	intVal := big.NewInt(0).SetBytes(ip)
	return intVal, nil
}

// ConvertIPv6ToFullAddr convert an ipv6 address to a full ipv6 address
func ConvertIPv6ToFullAddr(ipv6 string) (string, error) {
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
		ip, err := ConvertIPv6ToFullAddr(value.(string))
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
				v, err := ConvertIPv6ToFullAddr(item.(string))
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
