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
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
)

func getBKTopo(ctx context.Context, db dal.RDB, opt *option) (*Topo, error) {
	result := &Topo{}
	objIds := make([]string, 0)
	root, err := getBKAppNode(ctx, db, opt)
	if nil != err {
		return nil, err
	}
	if opt.scope == "all" || opt.scope == common.BKInnerObjIDApp {
		assts, err := getMainlineAssociation(ctx, db, opt)
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
		procModules, err := getProcessTopo(ctx, db, opt, bizID)
		if nil != err {
			return nil, err
		}
		procTopo := &ProcessTopo{
			BizName:   opt.bizName,
			Processes: procModules,
		}
		result.ProcTopos = procTopo
	}

	if opt.mini {

		_, keys, err := getModelAttributes(ctx, db, opt, objIds)
		if nil != err {
			return nil, err
		}

		if result.BizTopo != nil {
			err := result.BizTopo.walk(func(node *Node) error {
				node.Data = util.CopyMap(node.Data, keys[node.ObjID], []string{common.BKInstParentStr})
				return nil
			})
			if err != nil {
				blog.Errorf("walk through biz topo failed, err: %+v", err)
			}
		}

		if result.ProcTopos != nil {
			for _, proc := range result.ProcTopos.Processes {
				proc.Data = util.CopyMap(proc.Data, append(keys[common.BKInnerObjIDProc], "bind_ip", "port", "protocol", "bk_func_name", "work_path", "bk_start_param_regex"), []string{common.BKInstParentStr, common.BKAppIDField, common.BKOwnerIDField})
			}
		}
	}

	return result, nil
}

func getBKAppNode(ctx context.Context, db dal.RDB, opt *option) (*Node, error) {
	bkApp := newNode(common.BKInnerObjIDApp)
	cond := map[string]interface{}{
		common.BKOwnerIDField: opt.OwnerID,
		common.BKAppNameField: opt.bizName,
	}

	err := db.Table(common.BKTableNameBaseApp).Find(cond).One(ctx, &bkApp.Data)
	if nil != err {
		return nil, fmt.Errorf("getBKAppNode error: %s", err.Error())
	}
	return bkApp, nil
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

	childFilter := map[string]interface{}{
		common.BKInstParentStr: instID,
	}

	switch asst.ObjectID {
	case common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule:
		childFilter[common.BKDefaultField] = map[string]interface{}{
			common.BKDBEQ: common.DefaultFlagDefaultValue,
		}
	default:
		childFilter[common.BKObjIDField] = asst.ObjectID
	}

	// blog.InfoJSON("get children for %s:%d", asst.ObjectID, instID)
	children := make([]map[string]interface{}, 0)
	tableName := common.GetInstTableName(asst.ObjectID)

	err = db.Table(tableName).Find(childFilter).All(ctx, &children)
	if nil != err {
		return fmt.Errorf("get inst for %s error: %s", asst.ObjectID, err.Error())
	}

	for _, child := range children {
		root.Children = append(root.Children, &Node{ObjID: asst.ObjectID, Data: child})
	}

	child := pcmap[asst.ObjectID]
	if child == nil {
		return nil
	}
	for _, child := range root.Children {
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

func getMainlineAssociation(ctx context.Context, db dal.RDB, opt *option) ([]*metadata.Association, error) {
	assts := make([]*metadata.Association, 0)

	filter := map[string]interface{}{
		common.AssociationKindIDField: common.AssociationKindMainline,
	}

	err := db.Table(common.BKTableNameObjAsst).Find(filter).All(ctx, &assts)
	if nil != err {
		return nil, fmt.Errorf("query cc_ObjAsst error: %s", err.Error())
	}
	return assts, nil
}

func getProcessTopo(ctx context.Context, db dal.RDB, opt *option, bizID uint64) ([]*Process, error) {
	// fetch all process module
	procModules := make([]ProModule, 0)
	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
	}
	err := db.Table(common.BKTableNameProcModule).Find(filter).All(ctx, &procModules)
	if nil != err {
		return nil, fmt.Errorf("get process faile %s", err.Error())
	}

	procModMap := map[uint64][]string{} // processID -> modules
	procIDs := make([]uint64, 0)
	for _, pm := range procModules {
		procIDs = append(procIDs, pm.ProcessID)
		procModMap[pm.ProcessID] = append(procModMap[pm.ProcessID], pm.ModuleName)
	}

	// fetch all process
	procs := make([]map[string]interface{}, 0)
	err = db.Table(common.BKTableNameBaseProcess).Find(filter).All(ctx, &procs)
	if nil != err {
		return nil, fmt.Errorf("get process faile %s", err.Error())
	}

	topos := make([]*Process, 0)
	for _, proc := range procs {
		topo := Process{
			Data: proc,
		}
		procID, err := getInt64(proc[common.BKProcessIDField])
		if nil != err {
			return nil, err
		}
		topo.Modules = procModMap[procID]
		topos = append(topos, &topo)
	}

	return topos, nil
}
