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

package findopt

import (
	"configcenter/src/common"
)

// SortItem used to define the sort condition
type SortItem struct {
	Name       string
	Descending bool
}

// FieldItem used to define the field codnition
type FieldItem struct {
	Name string
	Hide bool
}

// Opts options used to find and modify
type Opts struct {
	Fields []FieldItem
	Sort   []SortItem
	Limit  int64
	Skip   int64
}

// One find one options
type One struct {
	Opts
}

// Many find many options
type Many struct {
	Opts
}

// FindAndModify find and modify options
type FindAndModify struct {
	Opts
	Remove bool
	Upsert bool
	New    bool
}

// DefaultOpts default opts
var DefaultOpts = Opts{Limit: common.BKNoLimit}
