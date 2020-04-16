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

package topo_tree

import (
	"errors"
)

// if hit count larger than overHead, then the request will 
// return immediately with an error.
const overHead = 20

var OverHeadError = errors.New("hit data is overhead")

// all the search option is all matched with regexp and ignore case.
type SearchOption struct {
	BusinessID   int64       `json:"bk_biz_id"`
	// if BusinessID is > 0, then BusinessName will be ignored
	BusinessName string      `json:"bk_biz_name"`
	SetName      string      `json:"bk_set_name"`
	ModuleName   string      `json:"bk_module_name"`
	Level        CustomLevel `json:"bk_level"`
}

func (s SearchOption) Validate() error  {	
	if s.BusinessID == 0 && len(s.BusinessName) == 0 {
		return errors.New("bk_biz_id and bk_biz_name can not be empty at the same time")
	}
		
	return nil
}

// business topology custom level describe.
type CustomLevel struct {
	Object   string `json:"bk_obj_id"`
	InstName string `json:"bk_inst_name"`
}

type Topology struct {
	BusinessID   int64  `json:"bk_biz_id"`
	BusinessName string `json:"bk_biz_name"`
	Tree         Tree   `json:"bk_topo_tree"`
}

type Tree struct {
	Object   string `json:"bk_obj_id"`
	InstName string `json:"bk_inst_name"`
	InstID   int64  `json:"bk_inst_id"`
	Children []Tree `json:"children"`
}
