/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metadata

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// HostMapStr TODO
// host map with string type ip and operator, can only get host from db with this map
type HostMapStr map[string]interface{}

// UnmarshalBSON TODO
func (h *HostMapStr) UnmarshalBSON(b []byte) error {
	if h == nil {
		return bsonx.ErrNilDocument
	}
	elements, err := bsoncore.Document(b).Elements()
	if err != nil {
		return err
	}

	if *h == nil {
		*h = map[string]interface{}{}
	}
	for _, element := range elements {
		rawValue := element.Value()
		switch element.Key() {
		case common.BKHostInnerIPField:
			innerIP, err := parseBsonStringArrayValueToString(rawValue)
			if err != nil {
				return err
			}
			(*h)[common.BKHostInnerIPField] = string(innerIP)
		case common.BKHostOuterIPField:
			outerIP, err := parseBsonStringArrayValueToString(rawValue)
			if err != nil {
				return err
			}
			(*h)[common.BKHostOuterIPField] = string(outerIP)
		case common.BKOperatorField:
			operator, err := parseBsonStringArrayValueToString(rawValue)
			if err != nil {
				return err
			}
			(*h)[common.BKOperatorField] = string(operator)
		case common.BKBakOperatorField:
			bakOperator, err := parseBsonStringArrayValueToString(rawValue)
			if err != nil {
				return err
			}
			(*h)[common.BKBakOperatorField] = string(bakOperator)
		case common.BKHostInnerIPv6Field:
			innerIPv6, err := parseBsonStringArrayValueToString(rawValue)
			if err != nil {
				return err
			}
			(*h)[common.BKHostInnerIPv6Field] = string(innerIPv6)
		case common.BKHostOuterIPv6Field:
			outerIPv6, err := parseBsonStringArrayValueToString(rawValue)
			if err != nil {
				return err
			}
			(*h)[common.BKHostOuterIPv6Field] = string(outerIPv6)
		default:
			dc := bsoncodec.DecodeContext{Registry: bson.DefaultRegistry}
			vr := bsonrw.NewBSONValueReader(rawValue.Type, rawValue.Data)
			decoder, err := bson.NewDecoderWithContext(dc, vr)
			if err != nil {
				return err
			}
			value := new(interface{})
			err = decoder.Decode(value)
			if err != nil {
				return err
			}
			(*h)[element.Key()] = *value
		}
	}
	return nil
}

func parseBsonStringArrayValueToString(value bsoncore.Value) ([]byte, error) {
	switch value.Type {
	case bsontype.Array:
		rawArray, rem, ok := bsoncore.ReadArray(value.Data)
		if !ok {
			return nil, bsoncore.NewInsufficientBytesError(value.Data, rem)
		}
		array, err := rawArray.Values()
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		arrayLen := len(array)
		for index, arrayValue := range array {
			if arrayValue.Type != bsontype.String {
				return nil, fmt.Errorf("invalid BSON type %v", arrayValue.Type)
			}
			str, rem, ok := bsoncore.ReadString(arrayValue.Data)
			if !ok {
				return nil, bsoncore.NewInsufficientBytesError(arrayValue.Data, rem)
			}
			buf.WriteString(str)
			if index != arrayLen-1 {
				buf.WriteByte(',')
			}
		}
		return buf.Bytes(), nil
	case bsontype.Null:
		return []byte{}, nil
	default:
		return nil, fmt.Errorf("invalid BSON type %v", value.Type)
	}
}

// StringArrayToString TODO
// special field whose string array value is parsed into string value from db
type StringArrayToString string

// UnmarshalBSONValue TODO
func (s *StringArrayToString) UnmarshalBSONValue(typo bsontype.Type, raw []byte) error {
	if s == nil {
		return bsonx.ErrNilDocument
	}
	value := bsoncore.Value{
		Type: typo,
		Data: raw,
	}
	str, err := parseBsonStringArrayValueToString(value)
	if err != nil {
		return err
	}
	*s = StringArrayToString(str)
	return err
}

// HostSpecialFields Special fields in the host attribute, in order to fuzzy query the following fields are stored in
// the database as an array.
var HostSpecialFields = []string{common.BKHostInnerIPField, common.BKHostOuterIPField, common.BKOperatorField,
	common.BKBakOperatorField, common.BKHostInnerIPv6Field, common.BKHostOuterIPv6Field}

// hostIpv6Fields host needs to convert to full format ipv6 field, the field need to in HostSpecialFields
var hostIpv6Fields = map[string]struct{}{
	common.BKHostInnerIPv6Field: {},
	common.BKHostOuterIPv6Field: {},
}

// ConvertHostSpecialStringToArray convert host special string to array
// convert host ip and operator fields value from string to array
// NOTICE: if host special value is empty, convert it to null to trespass the unique check, **do not change this logic**
func ConvertHostSpecialStringToArray(host map[string]interface{}) (map[string]interface{}, error) {
	var err error
	for _, field := range HostSpecialFields {
		value, ok := host[field]
		if !ok {
			continue
		}
		switch v := value.(type) {
		case string:
			v = strings.TrimSpace(v)
			v = strings.Trim(v, ",")
			if len(v) == 0 {
				host[field] = nil
				continue
			}
			items := strings.Split(v, ",")
			if _, ok := hostIpv6Fields[field]; !ok {
				host[field] = items
				continue
			}
			host[field], err = ConvertHostIpv6Val(items)
			if err != nil {
				return nil, err
			}

		case []string:
			if len(v) == 0 {
				host[field] = nil
				continue
			}
			if _, ok := hostIpv6Fields[field]; !ok {
				continue
			}
			host[field], err = ConvertHostIpv6Val(v)
			if err != nil {
				return nil, err
			}

		case []interface{}:
			if len(v) == 0 {
				host[field] = nil
			} else {
				blog.Errorf("host %s type invalid, value %v", field, host[field])
			}
		case nil:
		default:
			blog.Errorf("host %s type invalid, value %v", field, host[field])
		}
	}
	return host, nil
}

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
