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

package logics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	sceneUtil "configcenter/src/scene_server/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
)

type asstObjectInst struct {
	*backbone.Engine
	pheader        http.Header
	ownerID        string
	fields         map[string]*metadata.ObjAttDes
	asstPrimaryKey map[string][]metadata.ObjAttDes // import object assocate object primary fields [object id][]object fiedls
	asstInstConds  map[string][]interface{}        // import object assocate fields search instance data condition
	// import object assocate object primary fields and instance data relation map, map[object id]map[primary key val string]instance id
	instPrimayIDMap map[string]map[string]int64
}

func NewAsstObjectInst(pheader http.Header, engine *backbone.Engine, ownerID string, fields map[string]*metadata.ObjAttDes) *asstObjectInst {
	return &asstObjectInst{
		Engine:  engine,
		ownerID: ownerID,
		fields:  fields,
		pheader: pheader,
	}
}

// GetObjAsstObjectPrimaryKey  get instance assocate object primary property fields
func (a *asstObjectInst) GetObjAsstObjectPrimaryKey() error {
	ret := make(map[string][]metadata.ObjAttDes)

	for _, f := range a.fields {
		if util.IsAssocateProperty(f.PropertyType) {
			fields, err := a.getObjectFields(f.AssociationID, common.BKPropertyIDField)
			if nil != err {
				blog.Errorf("get object  assocate property %s error:%s", f.PropertyID, err.Error())
				return err
			}
			var primaryFields []metadata.ObjAttDes
			for _, f := range fields {
				if true == f.IsOnly {
					primaryFields = append(primaryFields, f)
				}
			}
			ret[f.AssociationID] = primaryFields
		}
	}
	a.asstPrimaryKey = ret
	return nil

}

// SetObjAsstPropertyVal set instance assocate object value to property fields
func (a *asstObjectInst) SetObjAsstPropertyVal(inst map[string]interface{}) error {

	for key, val := range inst {
		f, ok := a.fields[key]
		if false == ok {
			continue
		}

		if util.IsAssocateProperty(f.PropertyType) {
			strInsts, _ := val.(string)
			if common.ExcelDelAsstObjectRelation == strings.TrimSpace(strInsts) {
				inst[key] = ""
				continue
			}
			// assocate object relation no change
			if "" == strInsts {
				continue
			}

			insts := strings.Split(strInsts, common.ExcelAsstPrimaryKeyRowChar)
			var strIds []string
			for _, inst := range insts {
				inst = strings.TrimSpace(inst)
				if "" == inst {
					continue
				}
				id, err := a.GetIDsByExcelStr(f.AssociationID, inst)
				if nil != err {
					return err
				}
				strIds = append(strIds, fmt.Sprintf("%d", id))
			}

			inst[key] = strings.Join(strIds, common.InstAsstIDSplit)

		}
	}
	return nil
}

// InitInstFromData get assocate object instance data, return map[row]error,  error
func (a *asstObjectInst) InitInstFromData(infos map[int64]map[string]interface{}) (map[int64]error, error) {
	rowErr := a.getAsstObjectConds(infos)
	if 0 != len(rowErr) {
		return rowErr, nil
	}

	err := a.getAsstInstByAsstObjectConds()

	if nil != err {
		return nil, err
	}

	return nil, nil
}

// GetIDsByExcelStr Get a string of data based on multiple primary key values (multiple primary key values are sorted by field name, separated by ##)
func (a *asstObjectInst) GetIDsByExcelStr(objID, key string) (int64, error) {
	if nil == a.instPrimayIDMap {
		return 0, fmt.Errorf("%s primary %s not found", objID, key)
	}
	objIDInst, ok := a.instPrimayIDMap[objID]
	if false == ok {
		return 0, fmt.Errorf("%s primary %s not found", objID, key)
	}
	id, ok := objIDInst[key]
	if false == ok {
		return 0, fmt.Errorf("%s primary %s not found", objID, key)
	}
	return id, nil

}

// getAsstInstByAsstObjectConds  parse import instance assocate property field value to search paramaters
func (a *asstObjectInst) getAsstInstByAsstObjectConds() error {
	ret := make(map[string]map[string]int64)
	for objID, conds := range a.asstInstConds {
		isExist := false
		for _, f := range a.fields {
			if f.AssociationID == objID {
				isExist = true
			}
		}

		if false == isExist {
			continue
		}
		searchObjID := objID

		if !util.IsInnerObject(objID) {
			searchObjID = common.BKINnerObjIDObject
		}

		_, data, err := a.getInstData(searchObjID, conds)
		if err != nil {
			return err
		}
		if nil == data {
			ret[objID] = make(map[string]int64, 0)
			continue
		}
		primaryKey, _ := a.asstPrimaryKey[objID]
		for _, item := range data {
			keys := make([]string, 0)
			for _, f := range primaryKey {
				key, ok := item[f.PropertyID]
				if false == ok {
					errMsg := a.Language.Languagef("import_str_asst_str_query_data_format_error", objID, f.PropertyID)
					blog.Error(errMsg)
					return errors.New(errMsg)
				}
				keys = append(keys, fmt.Sprintf("%v", key))
			}
			_, ok := ret[objID]
			if false == ok {
				ret[objID] = make(map[string]int64)
			}
			id, err := a.getInstIDFromMapByObjectID(objID, item)
			if nil != err {
				err := fmt.Errorf("get object %s  inst id error, inst info:%v, err:%s ", objID, item, err.Error())
				blog.Errorf("%s ", err.Error())

				return err
			}

			ret[objID][strings.Join(keys, common.ExcelAsstPrimaryKeySplitChar)] = id

		}

	}
	a.instPrimayIDMap = ret
	return nil

}

func (a *asstObjectInst) getInstIDFromMapByObjectID(objType string, mapInst mapstr.MapStr) (int64, error) {
	idField := common.GetInstIDField(objType)
	idInterface, ok := mapInst[idField]
	if false == ok {
		return 0, fmt.Errorf("%s %s not found", objType, idField)
	}
	return util.GetInt64ByInterface(idInterface)
}

// getInstData  get instance data by condition
func (a *asstObjectInst) getInstData(objectID string, conds []interface{}) (int, []mapstr.MapStr, error) {
	query := &metadata.QueryInput{
		Condition: common.KvMap{common.BKDBOR: conds},
	}
	result, err := a.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), objectID, a.pheader, query)
	if err != nil || (err == nil && !result.Result) {
		return 0, nil, fmt.Errorf("%v, %v", err, result.ErrMsg)
	}

	return result.Data.Count, result.Data.Info, nil
}

// getAsstObjectConds get all  assocatie field object condition
func (a *asstObjectInst) getAsstObjectConds(infos map[int64]map[string]interface{}) map[int64]error {

	errs := make(map[int64]error, 0)
	asstMap := make(map[string][]interface{}) //map[AssociationID][]condition

	for rowIndex, info := range infos {
		for key, val := range info {
			f, ok := a.fields[key]
			if false == ok {
				continue
			}
			if util.IsAssocateProperty(f.PropertyType) {

				asstFields, ok := a.asstPrimaryKey[f.AssociationID]
				if false == ok {
					errs[rowIndex] = errors.New(a.Language.Languagef("import_asst_property_str_not_found", key))
					continue
				}

				strVal, ok := val.(string)
				if false == ok {
					errs[rowIndex] = errors.New(a.Language.Languagef("import_property_str_format_error", key))
					continue
				}

				if common.ExcelDelAsstObjectRelation == strings.TrimSpace(strVal) {
					continue
				}
				rows := strings.Split(strVal, common.ExcelAsstPrimaryKeyRowChar)

				asstConds := make([]interface{}, 0)
				for _, row := range rows {
					if "" == row {
						continue
					}
					primaryKeys := strings.Split(row, common.ExcelAsstPrimaryKeySplitChar)
					if len(primaryKeys) != len(asstFields) {
						errs[rowIndex] = errors.New(a.Language.Languagef("import_asst_property_str_primary_count_len", key))
						continue
					}
					conds := common.KvMap{}
					if false == util.IsInnerObject(f.AssociationID) {
						conds[common.BKObjIDField] = f.AssociationID
					}
					for i, val := range primaryKeys {

						asstf := asstFields[i]
						var err error
						conds[asstf.PropertyID], err = sceneUtil.ConvByPropertytype(&asstf, val)
						if nil != err {
							errs[rowIndex] = errors.New(a.Language.Languagef("import_asst_property_str_primary_count_len", key))
							continue
						}
					}
					asstConds = append(asstConds, conds)

				}

				_, ok = asstMap[f.AssociationID]
				if ok {
					asstMap[f.AssociationID] = append(asstMap[f.AssociationID], asstConds...)
				} else {
					asstMap[f.AssociationID] = asstConds
				}

			}

		}

	}
	a.asstInstConds = asstMap

	return errs
}

//GetObjectFields get object fields
func (a *asstObjectInst) getObjectFields(objID, sort string) ([]metadata.ObjAttDes, error) {
	page := metadata.BasePage{Start: 0, Limit: common.BKNoLimit, Sort: common.BKPropertyIDField}
	query := hutil.NewOperation().WithObjID(objID).WithOwnerID(a.ownerID).WithPage(page).Data()
	result, err := a.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), a.pheader, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("search host attributes failed, err: %v, result err: %s", err, result.ErrMsg)
	}
	attributesDesc := make([]metadata.ObjAttDes, 0)
	for _, att := range result.Data {
		attributesDesc = append(attributesDesc, metadata.ObjAttDes{Attribute: att})
	}
	for idx, attr := range attributesDesc {
		if !util.IsAssocateProperty(attr.PropertyType) {
			continue
		}

		cond := hutil.NewOperation().WithPropertyID(attr.PropertyID).WithOwnerID(a.ownerID).WithObjID(attr.ObjectID).Data()
		assResult, err := a.CoreAPI.ObjectController().Meta().SelectObjectAssociations(context.Background(), a.pheader, cond)
		if err != nil || (err == nil && !result.Result) {
			return nil, fmt.Errorf("search host obj associations failed, err: %v, result err: %s", err, result.ErrMsg)
		}

		if 0 < len(assResult.Data) {
			attributesDesc[idx].AssociationID = assResult.Data[0].AsstObjID // by the rules, only one id
			attributesDesc[idx].AsstForward = assResult.Data[0].AsstForward // by the rules, only one id
		}
	}
	return attributesDesc, nil
}

// SetMapFields set import object property fields
func (a *asstObjectInst) SetMapFields(objID string) error {
	fields, err := a.getObjectFields(objID, "")
	if nil != err {
		return err
	}
	ret := make(map[string]*metadata.ObjAttDes)
	for index, f := range fields {
		ret[f.PropertyID] = &fields[index]
	}
	a.fields = ret
	return nil
}
