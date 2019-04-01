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
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
)

func getBKTopo(ctx context.Context, db dal.RDB, opt *option) (*Topo, error) {
	result := &Topo{}
	objIds := []string{}
	root, err := getBKAppNode(ctx, db, opt)
	if nil != err {
		return nil, err
	}
	if opt.scope == scopeAll || opt.scope == common.BKInnerObjIDApp {
		assts, err := getAsst(ctx, db, opt)
		if nil != err {
			return nil, err
		}
		topo, err := getMainline(common.BKInnerObjIDApp, assts)
		if nil != err {
			return nil, err
		}
		objIds = append(objIds, topo...)
		pcmap := getPCmap(assts)

		err = getTree(ctx, db, root, pcmap)
		if nil != err {
			return nil, err
		}
		result.Mainline = topo
		result.BizTopo = root
	}

	if opt.scope == scopeAll || opt.scope == common.BKInnerObjIDProc {
		objIds = append(objIds, common.BKInnerObjIDProc)

		bizID, err := root.getInstID()
		if err != nil {
			return nil, err
		}
		procmodules, err := getProcessTopo(ctx, db, opt, bizID)
		if nil != err {
			return nil, err
		}
		proctopo := &ProcessTopo{
			BizName: opt.bizName,
			Procs:   procmodules,
		}
		result.ProcTopos = proctopo
	}

	if opt.mini {

		_, keys, err := getModelAttributes(ctx, db, opt, objIds)
		if nil != err {
			return nil, err
		}

		if result.BizTopo != nil {
			result.BizTopo.walk(func(node *Node) error {
				node.Data = util.CopyMap(node.Data, keys[node.ObjID], []string{common.BKInstParentStr})
				return nil
			})
		}

		if result.ProcTopos != nil {
			for _, proc := range result.ProcTopos.Procs {
				proc.Data = util.CopyMap(proc.Data, append(keys[common.BKInnerObjIDProc], "bind_ip", "port", "protocol", "bk_func_name", "work_path", "bk_start_param_regex"), []string{common.BKInstParentStr, common.BKAppIDField, common.BKOwnerIDField})
			}
		}
	}

	return result, nil
}

func getBKAppNode(ctx context.Context, db dal.RDB, opt *option) (*Node, error) {
	bkapp := newNode(common.BKInnerObjIDApp)
	condition := map[string]interface{}{
		common.BKOwnerIDField: opt.OwnerID,
		common.BKAppNameField: opt.bizName,
	}

	err := db.Table(common.BKTableNameBaseApp).Find(condition).One(ctx, &bkapp.Data)
	if nil != err {
		return nil, fmt.Errorf("getBKAppNode error: %s", err.Error())
	}
	return bkapp, nil
}

func getTree(ctx context.Context, db dal.RDB, root *Node, pcmap map[string]*metadata.Association) error {
	asst := pcmap[root.ObjID]
	if asst == nil {
		return nil
	}

	instID, err := root.getInstID()
	if nil != err {
		return nil
	}

	childCondition := condition.CreateCondition()
	childCondition.Field(common.BKInstParentStr).Eq(instID)

	switch asst.ObjectID {
	case common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule:
		childCondition.Field(common.BKDefaultField).NotGt(0)
	default:
		childCondition.Field(common.BKObjIDField).Eq(asst.ObjectID)
	}

	// blog.InfoJSON("get childs for %s:%d", asst.ObjectID, instID)
	childs := []map[string]interface{}{}
	tablename := common.GetInstTableName(asst.ObjectID)

	err = db.Table(tablename).Find(childCondition.ToMapStr()).All(ctx, &childs)
	if nil != err {
		return fmt.Errorf("get inst for %s error: %s", asst.ObjectID, err.Error())
	}

	for _, child := range childs {
		root.Childs = append(root.Childs, &Node{ObjID: asst.ObjectID, Data: child})
	}

	child := pcmap[asst.ObjectID]
	if child == nil {
		return nil
	}
	for _, child := range root.Childs {
		err = getTree(ctx, db, child, pcmap)
		if nil != err {
			return err
		}
	}
	return nil
}

// get parent -> child mapping
func getPCmap(assts []*metadata.Association) map[string]*metadata.Association {
	m := map[string]*metadata.Association{}
	for _, asst := range assts {
		child := getChileAsst(asst.AsstObjID, assts)
		if child != nil {
			m[asst.AsstObjID] = child
		}
	}
	return m
}

func getChileAsst(objID string, assts []*metadata.Association) *metadata.Association {
	if objID == common.BKInnerObjIDModule {
		return nil
	}
	for index := range assts {
		if assts[index].AsstObjID == objID {
			return assts[index]
		}
	}
	return nil
}

func getMainline(root string, assts []*metadata.Association) ([]string, error) {
	if root == common.BKInnerObjIDModule {
		return []string{common.BKInnerObjIDModule}, nil
	}
	for _, asst := range assts {
		if asst.AsstObjID == root {
			topo, err := getMainline(asst.ObjectID, assts)
			if err != nil {
				return nil, err
			}
			return append([]string{root}, topo...), nil
		}
	}
	return nil, fmt.Errorf("topo association broken: %+v", assts)
}

func getAsst(ctx context.Context, db dal.RDB, opt *option) ([]*metadata.Association, error) {
	assts := []*metadata.Association{}

	cond := condition.CreateCondition()
	cond.Field(common.AssociationKindIDField).Eq(common.AssociationKindMainline)

	err := db.Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).All(ctx, &assts)
	if nil != err {
		return nil, fmt.Errorf("query cc_ObjAsst error: %s", err.Error())
	}
	return assts, nil
}

func getProcessTopo(ctx context.Context, db dal.RDB, opt *option, bizID uint64) ([]*Process, error) {
	// fetch all process module
	procmodules := []ProModule{}
	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	err := db.Table(common.BKTableNameProcModule).Find(cond.ToMapStr()).All(ctx, &procmodules)
	if nil != err {
		return nil, fmt.Errorf("get process faile %s", err.Error())
	}

	procmodMap := map[uint64][]string{} // processID -> modules
	procIDs := []uint64{}
	for _, pm := range procmodules {
		procIDs = append(procIDs, pm.ProcessID)
		procmodMap[pm.ProcessID] = append(procmodMap[pm.ProcessID], pm.ModuleName)
	}

	// fetch all process
	procs := []map[string]interface{}{}
	err = db.Table(common.BKTableNameBaseProcess).Find(cond.ToMapStr()).All(ctx, &procs)
	if nil != err {
		return nil, fmt.Errorf("get process faile %s", err.Error())
	}

	topos := []*Process{}
	for _, proc := range procs {
		topo := Process{
			Data: proc,
		}
		procID, err := getInt64(proc[common.BKProcessIDField])
		if nil != err {
			return nil, err
		}
		topo.Modules = procmodMap[procID]
		topos = append(topos, &topo)
	}

	return topos, nil
}
