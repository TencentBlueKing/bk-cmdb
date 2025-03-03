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

package json

import (
	"bytes"
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// CutJsonDataWithFields cut jsonData and only return the "fields" be targeted.
// jsonData can not be nil, and must be a json string
func CutJsonDataWithFields(jsonData *string, fields []string) *string {
	if jsonData == nil {
		empty := ""
		return &empty
	}
	if len(fields) == 0 || *jsonData == "" {
		return jsonData
	}
	elements := gjson.GetMany(*jsonData, fields...)
	last := len(fields) - 1
	jsonBuffer := bytes.Buffer{}
	jsonBuffer.Write([]byte{'{'})
	for idx, field := range fields {
		jsonBuffer.Write([]byte{'"'})
		jsonBuffer.Write([]byte(field))
		jsonBuffer.Write([]byte{'"'})
		jsonBuffer.Write([]byte{':'})
		if elements[idx].Raw == "" {
			jsonBuffer.Write([]byte("null"))
		} else {
			jsonBuffer.Write([]byte(elements[idx].Raw))
		}
		if idx != last {
			jsonBuffer.Write([]byte{','})
		}
	}
	jsonBuffer.Write([]byte{'}'})
	cutOff := jsonBuffer.String()
	return &cutOff
}

// ReplaceJsonKey replace the oldKey with newKey in jsonData
func ReplaceJsonKey(jsonData []byte, keyMap map[string]string) ([]byte, error) {
	var err error
	for oldKey, newKey := range keyMap {
		jsonData, err = sjson.SetRawBytes(jsonData, newKey, []byte(gjson.GetBytes(jsonData, oldKey).Raw))
		if err != nil {
			return nil, fmt.Errorf("set %s key using %s value failed, err: %v", newKey, oldKey, err)
		}
		jsonData, err = sjson.DeleteBytes(jsonData, oldKey)
		if err != nil {
			return nil, fmt.Errorf("remove %s key failed, err: %v", oldKey, err)
		}
	}
	return jsonData, nil
}
