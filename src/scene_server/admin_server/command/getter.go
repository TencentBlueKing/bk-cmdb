package command

import (
	"configcenter/src/common"
	"configcenter/src/source_controller/api/metadata"
	"configcenter/src/source_controller/common/commondata"
	"configcenter/src/storage"
	"fmt"
)

func getBKTopo(db storage.DI, opt *option) (*Topo, error) {
	assts, err := getAsst(db, opt)
	if nil != err {
		return nil, err
	}
	topo, err := getMainline(common.BKInnerObjIDApp, assts)
	if nil != err {
		return nil, err
	}
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

	return &Topo{Mainline: topo, Topo: root}, nil
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

func getTree(db storage.DI, root *Node, pcmap map[string]*metadata.ObjectAsst) error {
	asst := pcmap[root.ObjID]
	if asst == nil {
		return nil
	}

	instID, err := root.getInstID()
	if nil != err {
		return nil
	}
	condition := map[string]interface{}{
		"bk_parent_id": instID,
	}

	// blog.InfoJSON("get childs for %s:%d", asst.ObjectID, instID)
	childs := []map[string]interface{}{}
	tablename := commondata.GetInstTableName(asst.ObjectID)
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

func getPCmap(assts []*metadata.ObjectAsst) map[string]*metadata.ObjectAsst {
	m := map[string]*metadata.ObjectAsst{}
	for _, asst := range assts {
		child := getChileAsst(asst.AsstObjID, assts)
		if child != nil {
			m[asst.AsstObjID] = child
		}
	}
	return m
}

func getChileAsst(objID string, assts []*metadata.ObjectAsst) *metadata.ObjectAsst {
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

func getMainline(root string, assts []*metadata.ObjectAsst) ([]string, error) {
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

func getAsst(db storage.DI, opt *option) ([]*metadata.ObjectAsst, error) {
	assts := []*metadata.ObjectAsst{}
	condition := map[string]interface{}{
		common.BKOwnerIDField: opt.OwnerID,
		"bk_object_att_id":    "bk_childid",
	}
	err := db.GetMutilByCondition("cc_ObjAsst", nil, condition, &assts, "", 0, 0)
	if nil != err {
		return nil, fmt.Errorf("query cc_ObjAsst error: %s", err.Error())
	}
	return assts, nil
}
