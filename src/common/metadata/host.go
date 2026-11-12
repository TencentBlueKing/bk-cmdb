/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package metadata

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// HostMapStr TODO
// host map with string type ip and operator, can only get host from db with this map
type HostMapStr map[string]interface{}

// UnmarshalBSON TODO
func (h *HostMapStr) UnmarshalBSON(b []byte) error {
	if h == nil {
		return mongo.ErrNilDocument
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
		return mongo.ErrNilDocument
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

// NormalizeHostIPv4List trims and validates the given ipv4 list, the duplicated
// addresses are removed while the order of the remaining ones is kept.
// NOTICE: an empty list means the host has no ipv4 address, it is a valid value.
func NormalizeHostIPv4List(ips []string) ([]string, error) {
	result := make([]string, 0, len(ips))
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if ip == "" {
			return nil, fmt.Errorf("ipv4 address can not be empty")
		}

		parsed := net.ParseIP(ip)
		if parsed == nil || parsed.To4() == nil || strings.Contains(ip, ":") {
			return nil, fmt.Errorf("invalid ipv4 address: %s", ip)
		}
		result = append(result, ip)
	}

	return util.RemoveDuplicatesAndEmptyByMap(result), nil
}

// NormalizeHostIPv6List trims and validates the given ipv6 list, every address is
// converted to the standard format(8 groups, lower case hex) and the duplicated
// addresses are removed while the order of the remaining ones is kept.
// NOTICE: an empty list means the host has no ipv6 address, it is a valid value.
func NormalizeHostIPv6List(ips []string) ([]string, error) {
	result := make([]string, 0, len(ips))
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if ip == "" {
			return nil, fmt.Errorf("ipv6 address can not be empty")
		}
		result = append(result, ip)
	}

	normalized, err := common.ConvertHostIpv6Val(result)
	if err != nil {
		return nil, fmt.Errorf("invalid ipv6 address: %s", err.Error())
	}
	return util.RemoveDuplicatesAndEmptyByMap(normalized), nil
}

// HostIPListEqual compares two ip lists without considering the order, but the
// contents must be the same. The ip list is stored as an array in the database and
// the query matches any element of it, so the order is not a part of the equality.
func HostIPListEqual(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for _, ip := range right {
		if !util.Contains(left, ip) {
			return false
		}
	}
	return true
}

// HostIPListEqualInOrder compares two ip lists in order. The order is meaningful
// because the first inner ip is used as the primary ip, so it is used to decide
// whether the destination ip list is different from the current one.
func HostIPListEqualInOrder(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for idx := range left {
		if left[idx] != right[idx] {
			return false
		}
	}
	return true
}

// ConvertHostIPListToString converts an ip list to the comma joined string which is
// the format accepted by the host special fields, an empty list is converted to an
// empty string which will be stored as null by ConvertHostSpecialStringToArray.
func ConvertHostIPListToString(ips []string) string {
	return strings.Join(ips, ",")
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
			host[field], err = common.ConvertHostIpv6Val(items)
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
			host[field], err = common.ConvertHostIpv6Val(v)
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

// GetHostDisplayName get host display name
func GetHostDisplayName(innerIP string, innerIPv6 string, cloudName string) string {
	if innerIP == "" {
		innerIP = "--"
	}

	if innerIPv6 == "" {
		innerIPv6 = "--"
	}

	return innerIP + "|" + innerIPv6 + "(" + cloudName + ")"
}
