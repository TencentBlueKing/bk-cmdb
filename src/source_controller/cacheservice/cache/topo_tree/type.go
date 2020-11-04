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
	BusinessID int64 `json:"bk_biz_id"`
	// if BusinessID is > 0, then BusinessName will be ignored
	BusinessName string      `json:"bk_biz_name"`
	SetName      string      `json:"bk_set_name"`
	ModuleName   string      `json:"bk_module_name"`
	Level        CustomLevel `json:"bk_level"`
}

func (s SearchOption) Validate() error {
	if s.BusinessID == 0 && len(s.BusinessName) == 0 {
		return errors.New("bk_biz_id and bk_biz_name can not be empty at the same time")
	}

	// check custom level
	if len(s.Level.Object) == 0 && len(s.Level.InstName) != 0 {
		return errors.New("bk_obj_id can not be empty")
	}

	if len(s.Level.Object) != 0 && len(s.Level.InstName) == 0 {
		return errors.New("bk_inst_name can not be empty")
	}

	hitCustom := false
	if len(s.Level.Object) != 0 && len(s.Level.InstName) != 0 {
		hitCustom = true
	}

	hitSet := false
	if len(s.SetName) != 0 {
		hitSet = true
	}

	if hitCustom && hitSet {
		return errors.New("bk_level, bk_set_name, bk_module can only be one of them")
	}

	hitModule := false
	if len(s.ModuleName) != 0 {
		hitModule = true
	}

	if (hitCustom && hitModule) || (hitSet && hitModule) {
		return errors.New("bk_level, bk_set_name, bk_module can only be one of them")
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
	Trees        []Tree `json:"bk_topo_tree"`
}

type Tree struct {
	Object   string `json:"bk_obj_id"`
	InstName string `json:"bk_inst_name"`
	InstID   int64  `json:"bk_inst_id"`
	Children []Tree `json:"children"`
}

type instance struct {
	name     string
	id       int64
	parentID int64
}

// search one of the node's path in the the business topology.
type SearchNodePathOption struct {
	Business int64           `json:"bk_biz_id"`
	Nodes    []*MainlineNode `json:"bk_nodes"`
}

type MainlineNode struct {
	// mainline topology object
	Object string `json:"bk_obj_id"`
	// object's instance id, eg, bk_set_id, bk_module_id...
	InstanceID int64 `json:"bk_inst_id"`
}

// NodePaths describes a topology instance's path from itself to the top biz.
// cause a host may exists in multiple modules, so it's may have several paths.
type NodePaths struct {
	*MainlineNode
	InstanceName string `json:"bk_inst_name"`
	// host may have multiple paths
	Paths [][]Node `json:"bk_paths"`
}

type Node struct {
	// mainline topology object
	Object string `json:"bk_obj_id"`
	// object's instance id, eg, bk_set_id, bk_module_id...
	InstanceID int64 `json:"bk_inst_id"`
	// instance's name, eg: bk_set_name, bk_module_name...
	InstanceName string `json:"bk_inst_name"`
	// node's parent id, only used for internal, do not return to user.
	ParentID int64 `json:"-"`
}

type module struct {
	ID       int64  `json:"bk_module_id"`
	Name     string `json:"bk_module_name"`
	Biz      int64  `json:"bk_biz_id"`
	ParentID int64  `json:"bk_parent_id"`
}

type set struct {
	ID       int64  `json:"bk_set_id"`
	Name     string `json:"bk_set_name"`
	Biz      int64  `json:"bk_biz_id"`
	ParentID int64  `json:"bk_parent_id"`
}

type biz struct {
	ID   int64  `json:"bk_biz_id"`
	Name string `json:"bk_biz_name"`
}

type custom struct {
	ID       int64  `json:"bk_inst_id"`
	ParentID int64  `json:"bk_parent_id"`
	Name     string `json:"bk_inst_name"`
}
