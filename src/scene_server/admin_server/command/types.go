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

package command

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
)

type option struct {
	OwnerID  string
	position string
	dryrun   bool
	mini     bool
	scope    string
	bizName  string
}

// Node topo node define
type Node struct {
	ObjID   string                 `json:"bk_obj_id,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Childs  []*Node                `json:"childs,omitempty"`
	nodekey string
	mark    string
}

func (n *Node) getChildObjID() string {
	for _, child := range n.Childs {
		return child.ObjID
	}
	return ""
}

func (n *Node) getInstID() (uint64, error) {
	id, err := getInt64(n.Data[common.GetInstIDField(n.ObjID)])
	if nil != err {
		return 0, fmt.Errorf("node has no instID: %+v", *n)
	}
	return id, nil
}

func (n *Node) getChilDInstNames() (instnames []string) {
	instnamefield := n.getChilDInstNameField()
	for _, child := range n.Childs {
		name, ok := child.Data[instnamefield].(string)
		if !ok {
			blog.Errorf("child has no instname field %#v", child.Data[instnamefield])
			continue
		}
		instnames = append(instnames, name)
	}
	return
}
func (n *Node) getChilDInstNameField() string {
	for _, child := range n.Childs {
		return common.GetInstNameField(child.ObjID)
	}
	return ""
}
func (n *Node) getInstNameField() string {
	return common.GetInstNameField(n.ObjID)
}

func getInt64(v interface{}) (uint64, error) {
	switch tv := v.(type) {
	case int8:
		return uint64(tv), nil
	case int16:
		return uint64(tv), nil
	case int32:
		return uint64(tv), nil
	case int64:
		return uint64(tv), nil
	case int:
		return uint64(tv), nil
	case uint8:
		return uint64(tv), nil
	case uint16:
		return uint64(tv), nil
	case uint32:
		return uint64(tv), nil
	case uint64:
		return uint64(tv), nil
	case uint:
		return uint64(tv), nil
	case float32:
		return uint64(tv), nil
	case float64:
		return uint64(tv), nil
	default:
		return 0, fmt.Errorf("v is not number : %+v", v)
	}
}

func newNode(objID string) *Node {
	return &Node{ObjID: objID, Data: map[string]interface{}{}, Childs: []*Node{}}
}

// result: map[parentkey]map[childkey]node
func (n *Node) walk(walkfunc func(node *Node) error) error {
	for _, child := range n.Childs {
		if err := child.walk(walkfunc); nil != err {
			return err
		}
	}
	if err := walkfunc(n); nil != err {
		return err
	}
	return nil
}

// // nodekey: model2[key1:value,key2:value]-model2[key1:value,key2:value]
// func (n *Node) getNodekey(parentKey string, keys []string) (nodekey string) {
// 	nodekey = n.ObjID + "["
// 	if "" != parentKey {
// 		nodekey = parentKey + "-" + nodekey
// 	}
// 	kv := []string{}
// 	for _, key := range keys {
// 		kv = append(kv, key+":"+fmt.Sprint(n.Data[key]))
// 	}
// 	nodekey = strings.Join(kv, ",") + "]"
// 	return nodekey
// }

// Topo define
type Topo struct {
	Mainline  []string     `json:"mainline,omitempty"`
	BizTopo   *Node        `json:"biz_topo,omitempty"`
	ProcTopos *ProcessTopo `json:"proc_topo,omitempty"`
}

type ProModule struct {
	ProcessID  uint64 `json:"bk_process_id" bson:"bk_process_id,omitempty"`
	ModuleName string `json:"bk_module_name" bson:"bk_module_name,omitempty"`
	BizID      uint64 `json:"bk_biz_id" bson:"bk_biz_id,omitempty"`
	OwnerID    string `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

type Process struct {
	Data    map[string]interface{} `json:"data"`
	Modules []string               `json:"modules"`
}
type ProcessTopo struct {
	BizName string     `json:"bk_biz_name"`
	Procs   []*Process `json:"procs"`
}

const actionCreate = "create"
const actionUpdate = "update"
