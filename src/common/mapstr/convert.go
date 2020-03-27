/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mapstr

import (
	"bytes"
	"encoding/json"
)

// DecodeFromMapStr convert input into json, then decode json into data
// 接口背景：mapstr 直接解析结构体实现的不完整，有很多坑点，已知问题：结构体中指针类型会导致 mapstr 解析结构体异常。
// 新的问题：mapstr 转json时数据会丢失
func DecodeFromMapStr(data interface{}, input MapStr) error {
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}
	d := json.NewDecoder(bytes.NewReader(inputBytes))
	d.UseNumber()
	err = d.Decode(data)
	return err
}
