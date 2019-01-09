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

package mapstr

import (
	"reflect"
	"strings"
)

// GetTags parse a object and get the all tags
func GetTags(target interface{}, tagName string) []string {

	targetType := reflect.TypeOf(target)
	switch targetType.Kind() {
	default:
		break
	case reflect.Ptr:
		targetType = targetType.Elem()

	}

	numField := targetType.NumField()
	tags := make([]string, 0)
	for i := 0; i < numField; i++ {
		structField := targetType.Field(i)
		if tag, ok := structField.Tag.Lookup("field"); ok {
			tags = append(tags, tag)
		}
	}
	return tags

}

// SetValueToMapStrByTags  convert a struct to MapStr by tags default tag name is field
func SetValueToMapStrByTags(source interface{}) MapStr {
	return SetValueToMapStrByTagsWithTagName(source, "field")
}

// SetValueToMapStrByTagsWithTagName convert a struct to MapStr by tags
func SetValueToMapStrByTagsWithTagName(source interface{}, tagName string) MapStr {

	values := MapStr{}
	if nil == source {
		return values
	}

	targetType := reflect.TypeOf(source)
	targetValue := reflect.ValueOf(source)
	switch targetType.Kind() {
	case reflect.Ptr:
		targetType = targetType.Elem()
		targetValue = targetValue.Elem()

		if targetType.Kind() == reflect.Ptr {
			return SetValueToMapStrByTagsWithTagName(targetValue.Interface(), tagName)
		}

	}

	numField := targetType.NumField()
	for i := 0; i < numField; i++ {
		structField := targetType.Field(i)
		tag, ok := structField.Tag.Lookup(tagName)
		if !ok && !structField.Anonymous {
			continue
		}

		if (0 == len(tag) || strings.Contains(tag, "ignoretomap")) && !structField.Anonymous {
			continue
		}
		tags := strings.Split(tag, ",")
		if 0 == len(tag) {
			tags = []string{structField.Name}
		}

		fieldValue := targetValue.FieldByName(structField.Name)
		if !fieldValue.CanInterface() {
			continue
		}

		switch structField.Type.Kind() {
		case reflect.String, reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Map:
			values.Set(tags[0], fieldValue.Interface())
		case reflect.Struct:
			innerMapStr := SetValueToMapStrByTagsWithTagName(fieldValue.Interface(), tagName)
			values.Set(tags[0], innerMapStr)

		case reflect.Ptr:

			innerValue := dealPointer(fieldValue, tags[0], tagName)
			values.Set(tags[0], innerValue)

		}

	}

	return values
}

// SetValueToStructByTags set the struct object field value by tags, default tag name is field
func SetValueToStructByTags(target interface{}, values MapStr) error {
	return SetValueToStructByTagsWithTagName(target, values, "field")
}

// SetValueToStructByTagsWithTagName set the struct object field value by tags
func SetValueToStructByTagsWithTagName(target interface{}, values MapStr, tagName string) error {

	targetType := reflect.TypeOf(target)
	targetValue := reflect.ValueOf(target)

	return parseStruct(targetType, targetValue, values, tagName)
}
