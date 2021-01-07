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
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var iteratorJson = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	UseNumber:              true,
}.Froze()

func MarshalToString(v interface{}) (string, error) {
	return iteratorJson.MarshalToString(v)
}

func Marshal(v interface{}) ([]byte, error) {
	return iteratorJson.Marshal(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return iteratorJson.MarshalIndent(v, prefix, indent)
}

func UnmarshalFromString(str string, v interface{}) error {
	return iteratorJson.UnmarshalFromString(str, v)
}

func Unmarshal(data []byte, v interface{}) error {
	return iteratorJson.Unmarshal(data, v)
}

func UnmarshalArray(items []string, result interface{}) error {
	strArrJSON := "[" + strings.Join(items, ",") + "]"
	return iteratorJson.Unmarshal([]byte(strArrJSON), result)
}
