package command

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/source_controller/api/metadata"
	"configcenter/src/source_controller/common/commondata"
	"configcenter/src/source_controller/common/instdata"
	"configcenter/src/storage"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

func importer(db storage.DI, opt *option) error {
	file, err := os.OpenFile(opt.position, os.O_RDONLY, os.ModePerm)
	if nil != err {
		return err
	}
	defer file.Close()

	topo := new(Topo)
	json.NewDecoder(file).Decode(topo)
	if nil != err {
		return err
	}
	err = importBKBiz(db, opt, topo)
	if nil != err {
		return fmt.Errorf("import topo error: %s", err.Error())
	}

	return nil
}

func importBKBiz(db storage.DI, opt *option, tar *Topo) error {
	cur, err := getBKTopo(db, opt)
	if err != nil {
		return fmt.Errorf("get src topo faile %s", err.Error())
	}

	//topo check
	if !compareSlice(tar.Mainline, cur.Mainline) {
		return fmt.Errorf("different topo mainline found, your expecting import topo is [%s], but the existing topo is [%s]",
			strings.Join(tar.Mainline, "->"), strings.Join(cur.Mainline, "->"))
	}

	modelAttributes, modelKeys, err := getModelAttributes(db, opt, cur)
	if err != nil {
		return err
	}

	screate, supdate, sdelete := compareTopo(cur, tar, modelAttributes, modelKeys)
	sort.SliceStable(screate, func(i, j int) bool {
		return screate[i].nodekey < screate[j].nodekey
	})
	sort.SliceStable(supdate, func(i, j int) bool {
		return supdate[i].nodekey < supdate[j].nodekey
	})
	sort.SliceStable(sdelete, func(i, j int) bool {
		return sdelete[i].nodekey > sdelete[j].nodekey
	})
	for _, node := range screate {
		tablename := commondata.GetInstTableName(node.ObjID)
		id, err := db.GetIncID(tablename)
		if nil != err {
			return fmt.Errorf("get increase id for  %s error: %s", tablename, err.Error())
		}
		idfield := instdata.GetIDNameByType(node.ObjID)
		node.Data[idfield] = id
		for _, child := range node.Childs {
			child.Data["bk_parent_id"] = id
		}
		_, err = db.Insert(tablename, node.Data)
		if nil != err {
			return fmt.Errorf("insert date to %s error: %s", tablename, err.Error())
		}
		blog.Infof("inserted %s to data: %+v", node.ObjID, node.Data)
	}
	for _, node := range supdate {
		tablename := commondata.GetInstTableName(node.ObjID)
		condition := getModifyCondition(node.Data, modelKeys[node.ObjID])
		updatedate := getUpdateData(node)
		err = db.UpdateByCondition(tablename, updatedate, condition)
		if nil != err {
			return fmt.Errorf("insert date to %s error: %s", tablename, err.Error())
		}
		blog.Infof("updated %s by %+v to data: %+v", node.ObjID, condition, updatedate)
	}
	for _, node := range sdelete {
		tablename := commondata.GetInstTableName(node.ObjID)
		condition := getModifyCondition(node.Data, modelKeys[node.ObjID])
		err = db.DelByCondition(tablename, condition)
		if nil != err {
			return fmt.Errorf("insert date to %s error: %s", tablename, err.Error())
		}
		blog.Infof("deleted %s to data: %+v", node.ObjID, condition)
	}

	return nil
}

func getModifyCondition(n map[string]interface{}, keys []string) map[string]interface{} {
	condition := map[string]interface{}{}
	for _, key := range keys {
		condition[key] = n[key]
	}
	return condition
}

var ignoreKeys = map[string]bool{
	"_id":           true,
	"create_time":   true,
	"bk_parent_id":  true,
	"default":       true,
	"bk_biz_id":     true,
	"bk_module_id":  true,
	"bk_process_id": true,
	"bk_inst_id":    true,
}

func getUpdateData(n *Node) map[string]interface{} {
	data := map[string]interface{}{}
	for key, value := range n.Data {
		if ignoreKeys[key] {
			continue
		}
		data[key] = value
	}
	return data
}

// compareSlice returns whether slice a,b exactly equal
func compareSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for index := range a {
		if a[index] != b[index] {
			return false
		}
	}
	return true
}

// compareTopo compare a & b topo and returns there difference
func compareTopo(cur, tar *Topo, modelAttributes map[string][]metadata.ObjectAttDes, modelKeys map[string][]string) (screate, supdate, sdelete []*Node) {
	tarExp := map[string]map[string]*Node{}
	curExp := map[string]map[string]*Node{}
	tar.Topo.expand("", modelKeys, tarExp)
	cur.Topo.expand("", modelKeys, curExp)

	for tarparent, tarchilds := range tarExp {
		curchilds := curExp[tarparent]
		for tarchildKey, tarchildNode := range tarchilds {
			if curchilds == nil {
				screate = append(screate, tarchildNode)
				continue
			}
			curchildNode := curchilds[tarchildKey]
			if curchildNode == nil {
				screate = append(screate, tarchildNode)
			} else {
				supdate = append(supdate, tarchildNode)
			}
		}
	}

	for curparent, curchilds := range curExp {
		tarchilds := tarExp[curparent]
		for curchildKey, curchildNode := range curchilds {
			if tarchilds == nil {
				sdelete = append(sdelete, curchildNode)
				continue
			}
			tarchildNode := tarchilds[curchildKey]
			if tarchildNode == nil {
				sdelete = append(sdelete, curchildNode)
			}

		}
	}

	return
}

func restoreIDs(tar *Node, modelAttributes map[string][]metadata.ObjectAttDes, modelKeys map[string][]string) {

}

func restoreID(db storage.DI, opt *option, node *Node, help *idhelper) error {
	switch node.ObjID {
	case "biz":
		condition := getModifyCondition(node.Data, []string{common.BKAppNameField})
		app := map[string]interface{}{}
		err := db.GetOneByCondition(commondata.GetInstTableName(node.ObjID), nil, condition, &app)
		if nil != err {
			return fmt.Errorf("get biz of %+v error: %s", condition, err.Error())
		}
		bizID, err := getInt64(app[common.BKAppIDField])
		if nil != err {
			return fmt.Errorf("get bizID faile, data: %+v, error: %s", app, err.Error())
		}
		help.bizID = bizID
		node.Data[common.BKAppIDField] = bizID
	case "set":
		condition := getModifyCondition(node.Data, []string{common.BKSetNameField})
		set := map[string]interface{}{}
		err := db.GetOneByCondition(commondata.GetInstTableName(node.ObjID), nil, condition, &set)
		if nil != err {
			return fmt.Errorf("get biz of %+v error: %s", condition, err.Error())
		}
		setID, err := getInt64(set[common.BKSetIDField])
		if nil != err {
			return fmt.Errorf("get setID faile, data: %+v, error: %s", set, err.Error())
		}
		node.Data[common.BKAppIDField] = setID

	case "module":

	}
}

type idhelper struct {
	bizID    int64
	setID    int64
	module   int64
	parentID int64
}

// getModelAttributes returns the model attributes
func getModelAttributes(db storage.DI, opt *option, cur *Topo) (modelAttributes map[string][]metadata.ObjectAttDes, modelKeys map[string][]string, err error) {
	condition := map[string]interface{}{
		common.BKObjIDField: map[string]interface{}{
			"$in": cur.Mainline,
		},
	}

	attributes := []metadata.ObjectAttDes{}
	err = db.GetMutilByCondition("cc_ObjAttDes", nil, condition, &attributes, "", 0, 0)
	if nil != err {
		return nil, nil, fmt.Errorf("faile to getModelAttributes for %v, error: %s", cur.Mainline, err.Error())
	}

	modelAttributes = map[string][]metadata.ObjectAttDes{}
	modelKeys = map[string][]string{}
	for _, att := range attributes {
		if att.IsOnly {
			modelKeys[att.ObjectID] = append(modelKeys[att.ObjectID], att.PropertyID)
		}
		modelAttributes[att.ObjectID] = append(modelAttributes[att.ObjectID], att)
	}

	for index := range modelKeys {
		sort.Strings(modelKeys[index])
	}

	return modelAttributes, modelKeys, nil
}
