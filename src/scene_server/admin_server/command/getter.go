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
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage"
)

func getBKTopo(db storage.DI, opt *option) (*Topo, error) {
	result := &Topo{}
	objIds := []string{}
	if opt.scope == "all" || opt.scope == common.BKInnerObjIDApp {
		assts, err := getAsst(db, opt)
		if nil != err {
			return nil, err
		}
		topo, err := getMainline(common.BKInnerObjIDApp, assts)
		if nil != err {
			return nil, err
		}
		objIds = append(objIds, topo...)
		root, err := getBKAppNode(db, opt)
		if nil != err {
			return nil, err
		}
		pcmap := getPCmap(assts)
		// blog.InfoJSON("%s", pcmap)
		err = getTree(db, root, pcmap)
		if nil != err {
			return nil, err
		}
		result.Mainline = topo
		result.BizTopo = root
	}

	if opt.scope == "all" || opt.scope == common.BKInnerObjIDProc {
		objIds = append(objIds, common.BKInnerObjIDProc)
		procmodules, err := getProcessTopo(db, opt)
		if nil != err {
			return nil, err
		}
		proctopo := &ProcessTopo{
			BizName: common.BKAppName,
			Procs:   procmodules,
		}
		result.ProcTopos = proctopo
	}

	if opt.mini {

		_, keys, err := getModelAttributes(db, opt, objIds)
		if nil != err {
			return nil, err
		}

		if result.BizTopo != nil {
			result.BizTopo.walk(func(node *Node) error {
				node.Data = util.CopyMap(node.Data, keys[node.ObjID], []string{"bk_parent_id"})
				return nil
			})
		}

		if result.ProcTopos != nil {
			for _, proc := range result.ProcTopos.Procs {
				proc.Data = util.CopyMap(proc.Data, append(keys[common.BKInnerObjIDProc], "bind_ip", "port", "protocol", "bk_func_name", "work_path"), []string{"bk_parent_id", "bk_biz_id", "bk_supplier_account"})
			}
		}
	}

	return result, nil
}

func getBKAppNode(db storage.DI, opt *option) (*Node, error) {
	bkapp := newNode(common.BKInnerObjIDApp)
	condition := map[string]interface{}{
		common.BKOwnerIDField: opt.OwnerID,
		common.BKAppNameField: common.BKAppName,
	}
	err := db.GetOneByCondition(common.BKTableNameBaseApp, nil, condition, &bkapp.Data)
	if nil != err {
		return nil, fmt.Errorf("getBKAppNode error: %s", err.Error())
	}
	return bkapp, nil
}

func getTree(db storage.DI, root *Node, pcmap map[string]*metadata.Association) error {
	asst := pcmap[root.ObjID]
	if asst == nil {
		return nil
	}

	instID, err := root.getInstID()
	if nil != err {
		return nil
	}
	condition := map[string]interface{}{
		common.BKInstParentStr: instID,
	}

	switch asst.ObjectID {
	case common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule:
		condition["default"] = map[string]interface{}{common.BKDBEQ: 0}
	default:
		condition[common.BKObjIDField] = asst.ObjectID
	}

	// blog.InfoJSON("get childs for %s:%d", asst.ObjectID, instID)
	childs := []map[string]interface{}{}
	tablename := common.GetInstTableName(asst.ObjectID)
	err = db.GetMutilByCondition(tablename, nil, condition, &childs, "", 0, 0)
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
		err = getTree(db, child, pcmap)
		if nil != err {
			return err
		}
	}
	return nil
}

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

func getAsst(db storage.DI, opt *option) ([]*metadata.Association, error) {
	assts := []*metadata.Association{}
	condition := map[string]interface{}{
		common.BKOwnerIDField:  opt.OwnerID,
		common.BKObjAttIDField: common.BKChildStr,
	}
	err := db.GetMutilByCondition("cc_ObjAsst", nil, condition, &assts, "", 0, 0)
	if nil != err {
		return nil, fmt.Errorf("query cc_ObjAsst error: %s", err.Error())
	}
	return assts, nil
}

func getProcessTopo(db storage.DI, opt *option) ([]*Process, error) {
	// fetch all process
	procs := []map[string]interface{}{}
	err := db.GetMutilByCondition(common.BKTableNameBaseProcess, nil, map[string]interface{}{}, &procs, "", 0, 0)
	if nil != err {
		return nil, fmt.Errorf("get process faile %s", err.Error())
	}

	// fetch all process module
	procmodules := []ProModule{}
	err = db.GetMutilByCondition(common.BKTableNameProcModule, nil, map[string]interface{}{}, &procmodules, "", 0, 0)
	if nil != err {
		return nil, fmt.Errorf("get process faile %s", err.Error())
	}

	procmodMap := map[int64][]string{} // processID -> modules
	for _, pm := range procmodules {
		procmodMap[pm.ProcessID] = append(procmodMap[pm.ProcessID], pm.ModuleName)
	}

	topos := []*Process{}
	for _, proc := range procs {
		topo := Process{
			Data: proc,
		}
		procID, err := getInt64(proc["bk_process_id"])
		if nil != err {
			return nil, err
		}
		topo.Modules = procmodMap[procID]
		topos = append(topos, &topo)
	}

	return topos, nil
}
