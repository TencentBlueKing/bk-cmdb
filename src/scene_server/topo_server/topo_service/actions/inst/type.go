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

package inst

type condition struct {
	InstID []int `json:"inst_ids"`
}

type deleteCondition struct {
	condition `json:",inline"`
}

type updateCondition struct {
	InstID   int                    `json:"inst_id"`
	InstInfo map[string]interface{} `json:"datas"`
}

type operation struct {
	Delete deleteCondition `json:"delete"`
	Update []updateCondition `json:"update"`
}

// instNameAsst  association inst name
type instNameAsst struct {
	ID         string                 `json:"id"`
	ObjID      string                 `json:"bk_obj_id"`
	ObjIcon    string                 `json:"bk_obj_icon"`
	InstID     int                    `json:"bk_inst_id"`
	ObjectName string                 `json:"bk_obj_name"`
	InstName   string                 `json:"bk_inst_name"`
	InstInfo   map[string]interface{} `json:"inst_info,omitempty"`
}

// commonInstTopo common inst topo
type commonInstTopo struct {
	instNameAsst
	Count    int            `json:"count"`
	Children []instNameAsst `json:"children"`
}

// commonInstTopoV2 common inst topo
type commonInstTopoV2 struct {
	Prev interface{} `json:"prev"`
	Next interface{} `json:"next"`
	Curr interface{} `json:"curr"`
}

type instTopoSort []commonInstTopo

func (cli instTopoSort) Len() int {
	return len(cli)
}

func (cli instTopoSort) Less(i, j int) bool {
	return cli[i].ID < cli[j].ID
}

func (cli instTopoSort) Swap(i, j int) {
	cli[i], cli[j] = cli[j], cli[i]
}

type instAsstSort []instNameAsst

func (cli instAsstSort) Len() int {
	return len(cli)
}

func (cli instAsstSort) Less(i, j int) bool {
	return cli[i].ID < cli[j].ID
}

func (cli instAsstSort) Swap(i, j int) {
	cli[i], cli[j] = cli[j], cli[i]
}
