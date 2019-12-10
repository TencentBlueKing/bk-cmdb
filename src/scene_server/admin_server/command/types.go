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
	"strings"

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
	ObjID    string                 `json:"bk_obj_id,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
	Children []*Node                `json:"childs,omitempty"`
	nodeKey  string
	mark     string
}

func (n *Node) getChildObjID() string {
	for _, child := range n.Children {
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

func (n *Node) getChildInstNames() (instNames []string) {
	instNameField := n.getChildInstNameField()
	for _, child := range n.Children {
		name, ok := child.Data[instNameField].(string)
		if !ok {
			blog.Errorf("child has no inst name field %#v", child.Data[instNameField])
			continue
		}
		instNames = append(instNames, name)
	}
	return
}
func (n *Node) getChildInstNameField() string {
	for _, child := range n.Children {
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
	return &Node{ObjID: objID, Data: map[string]interface{}{}, Children: []*Node{}}
}

// result: map[parentkey]map[childkey]node
func (n *Node) walk(walkFunc func(node *Node) error) error {
	for _, child := range n.Children {
		if err := child.walk(walkFunc); nil != err {
			return err
		}
	}
	if err := walkFunc(n); nil != err {
		return err
	}
	return nil
}

// getNodeKey outputs ==> `{parentKey}-{objectID}[{key1}:{value1},{key2}:{value2}]`
func (n *Node) getNodeKey(parentKey string, keys []string) (nodeKey string) {
	if "" != parentKey {
		nodeKey = parentKey + "-"
	}
	nodeKey += n.ObjID + "["
	kv := make([]string, 0)
	for _, key := range keys {
		item := fmt.Sprintf("%s:%s", key, n.Data[key])
		kv = append(kv, item)
	}
	nodeKey += strings.Join(kv, ",") + "]"
	return nodeKey
}

// BKTopo 蓝鲸安装拓扑的的结构
type BKTopo struct {
	// map[process name]map[key]value
	Proc               []map[string]interface{} `json:"proc"`
	ServiceTemplateArr []BKServiceTemplate      `json:"service_template"`
	Topo               BKBizTopo                `json:"topo"`
}

// BKServiceTemplate 服务模版与进程关系
type BKServiceTemplate struct {
	Name string `json:"name"`
	// 只能有前面两个有效, 没有内容的时候，使用Default
	ServiceCategoryName []string `json:"service_category_name"`
	// ServiceCategoryName,服务分类转换成分类ID
	ServiceCategoryID int64    `json:"-"`
	BindProcess       []string `json:"bind_proc"`
}

// BKBizTopo business set,module
type BKBizTopo struct {
	SetArr    []map[string]interface{} `json:"set"`
	ModuleArr []BKBizModule            `json:"module"`
}

// BKBizModule business module info
type BKBizModule struct {
	SetName         string                 `json:"bk_set_name"`
	ServiceTemplate string                 `json:"service_template"`
	Info            map[string]interface{} `json:"info"`
}

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
	BizName   string     `json:"bk_biz_name"`
	Processes []*Process `json:"procs"`
}

const actionCreate = "create"
const actionUpdate = "update"
