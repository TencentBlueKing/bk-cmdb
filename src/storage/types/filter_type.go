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
 
package types

import (
	"github.com/mongodb/mongo-go-driver/bson"

	"configcenter/src/common"
)

type Filter interface{}

type FilterBuilder map[string]interface{}

func NewFilterBuilder() FilterBuilder {
	return FilterBuilder{}
}

func (f FilterBuilder) Build() Filter {
	return f
}

func (f FilterBuilder) ParseBytes(data []byte) error {
	return bson.Unmarshal(data, f)
}
func (f FilterBuilder) ParseStuct(data interface{}) error {
	out, err := bson.Marshal(data)
	if nil != err {
		return err
	}
	return bson.Unmarshal(out, f)
}
func (f FilterBuilder) Eq(field string, val interface{}) FilterBuilder {
	f[field] = val
	return f
}
func (f FilterBuilder) NotEq(field string, val interface{}) FilterBuilder {
	f[field] = map[string]interface{}{
		common.BKDBNE: val,
	}
	return f
}
func (f FilterBuilder) Like(field string, val interface{}) FilterBuilder {
	f[field] = map[string]interface{}{
		common.BKDBLIKE: val,
	}
	return f
}
func (f FilterBuilder) In(field string, val interface{}) FilterBuilder {
	f[field] = map[string]interface{}{
		common.BKDBIN: val,
	}
	return f
}
func (f FilterBuilder) NotIn(field string, val interface{}) FilterBuilder {
	f[field] = map[string]interface{}{
		common.BKDBNIN: val,
	}
	return f
}
func (f FilterBuilder) LowerThan(field string, val interface{}) FilterBuilder {
	f[field] = map[string]interface{}{
		common.BKDBLT: val,
	}
	return f
}
func (f FilterBuilder) LowerOrEq(field string, val interface{}) FilterBuilder {
	f[field] = map[string]interface{}{
		common.BKDBLTE: val,
	}
	return f
}
func (f FilterBuilder) GreaterThan(field string, val interface{}) FilterBuilder {
	f[field] = map[string]interface{}{
		common.BKDBGT: val,
	}
	return f
}
func (f FilterBuilder) GreaterOrEq(field string, val interface{}) FilterBuilder {
	f[field] = map[string]interface{}{
		common.BKDBGTE: val,
	}
	return f
}
func (f FilterBuilder) Exists(field string) FilterBuilder {
	f[field] = map[string]interface{}{
		common.BKDBExists: true,
	}
	return f
}

func (f FilterBuilder) Or(exprs ...FilterBuilder) FilterBuilder {
	f[common.BKDBOR] = exprs
	return f
}

func (f FilterBuilder) And(exprs ...FilterBuilder) FilterBuilder {
	f[common.BKDBAND] = exprs
	return f
}
