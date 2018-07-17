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

package filter

type Filter struct {
	// must match these key and value pairs
	MustKV map[string]interface{}

	// used for operations like eq, lt, lte, nin etc.
	// compare sequences is irrelevant with the filter result.
	Compares []Compare
}

func (f *Filter) ToDoc() Doc {
	d := make(Doc)
	for k, v := range f.MustKV {
		d[k] = v
	}

	for idx := range f.Compares {
		for k, v := range f.Compares[idx].ToDocument() {
			d[k] = v
		}
	}

	return d
}

type Compare struct {
	Key      string
	Value    interface{}
	Operator OperatorType
}

func (c *Compare) ToDocument() Doc {
	return newDoc(c.Key, c.Value, c.Operator)
}

type Doc map[string]interface{}

func newDoc(key string, value interface{}, oper OperatorType) Doc {
	return map[string]interface{}{key: map[string]interface{}{string(oper): value}}
}
